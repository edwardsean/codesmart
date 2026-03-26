package container

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// WorkspaceInfo holds the identifiers for a created workspace.
type WorkspaceInfo struct {
	ContainerID string
	VolumeName  string
}

type FileInfo struct {
	Path  string
	IsDir bool
}

type TerminalSession struct {
	ExecID string
	Conn   *types.HijackedResponse
}

type ContainerManager struct {
	cli *client.Client
}

func NewContainerManager() (*ContainerManager, error) {
	//client WithHostFromEnv tells the client to use env variables like DOCKER_HOST
	//the cli object is used to make API calls to docker
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &ContainerManager{cli: cli}, nil
}

func (m *ContainerManager) Close() error {
	return m.cli.Close()
}

// this function pulls the base image if it is not already present locally
// if doesnt exist, it downloads (pulls) it from Docker Hub
// This ensures when we later create a container, the base image is already available
func (m *ContainerManager) ensureImage(ctx context.Context, imageName string) error {
	//check if the image already exists locally
	_, err := m.cli.ImageInspect(ctx, imageName) //Asks Docker for detailed information about the image with the given name (e.g., ubuntu:22.04).
	if err == nil {
		return nil //image already present
	}

	//  If the error is something else (like network timeout, permission denied), we return that error.
	if !errdefs.IsNotFound(err) { //Checks if the error is a "not found" error (image missing) or some other error.
		return fmt.Errorf("failed to inspect image: %w", err)
	}

	//image not found, pull it
	//Returns: A io.ReadCloser (a stream of data) that we must read to know when the pull finishes.
	reader, err := m.cli.ImagePull(ctx, imageName, image.PullOptions{}) //Starts downloading the image from Docker Hub.
	//this line sends an HTTP request to docker daemon, when docker daemon receives image pull, it connects to docker hub, and downloads the image layer over https
	//stores them on disk at /var/lib/docker
	if err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}
	defer reader.Close() //Always close the stream to free resources.

	// Sends progress updates back to your Go program via the HTTP connection
	_, err = io.ReadAll(reader) //Reads all data from the stream until the stream ends.
	//What is io? The io package provides basic interfaces to I/O primitives. io.ReadAll reads until EOF (end of file) or error.
	// The reader you get is a stream of JSON progress messages from the Docker daemon, like:

	// {"status":"Pulling from library/ubuntu","id":"22.04"}
	// {"status":"Downloading","progressDetail":{"current":123456,"total":654321}}
	// {"status":"Download complete"}
	return err
}

// this function creates a dedicated bridge network for the user if it doesnt exist
// returns the network ID
// ensures that each user has their own isolated Docker network
// when a user's container are created, they will be attached to this network, preventing them
// from communicating with containers belonging to other users
func (m *ContainerManager) ensureUserNetwork(ctx context.Context, userId int) (string, error) {
	networkName := fmt.Sprintf("codesmart_user_%d", userId)
	networks, err := m.cli.NetworkList(ctx, network.ListOptions{}) //asks docker daemon to list all existing networks locally
	//contains a list of network objects:
	// ID: Unique identifier (e.g., abc123...)
	// Name: Human-readable name (e.g., codesmart_user_5)
	// Driver: Network type (e.g., bridge)
	// Created: When it was created

	if err != nil {
		return "", err
	}

	for _, net := range networks {
		if net.Name == networkName {
			return net.ID, nil
		}
	}

	resp, err := m.cli.NetworkCreate(ctx, networkName, network.CreateOptions{
		//bridge	Creates a private network on the host	Isolating groups of containers (what we want)
		// host	Container uses host's network directly	Performance-critical apps (less secure)
		// none	No network	Isolated containers
		// overlay	Connects containers across multiple hosts	Docker Swarm, Kubernetes
		Driver:   "bridge", //bridge is like a virtual network switch, connects containers together and allows them to communicate with each other
		Internal: false,    //allows outbound internet (git, npm, etc)
		//false = containers can access the internet (to download packages, git clone, etc.)
		//true = containers are completely isolated (no internet access)
		Labels: map[string]string{
			"codesmart.user": fmt.Sprintf("%d", userId), //Adds metadata to the network so we can later identify which user it belongs to.
		},
		//docker network inspect codesmart_user_5
		// # Shows labels including: "codesmart.user": "5"
	})
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

// this function creates a container, a volume, and clones the repo
// if gitURL is empty, it initializes a scratch git repo
// gitToken is optional (for private repos)
func (m *ContainerManager) CreateWorkspaceContainer(ctx context.Context, userId, projectId int, gitURL, gitToken string) (*WorkspaceInfo, error) {
	//create a named volume for the project files
	// A Docker volume is a persistent storage location that exists independently of containers
	volumeName := fmt.Sprintf("codesmart_user_%d_project_%d", userId, projectId)
	_, err := m.cli.VolumeCreate(ctx, volume.CreateOptions{
		Name: volumeName,
		Labels: map[string]string{
			"codesmart.user":    fmt.Sprintf("%d", userId),
			"codesmart.project": fmt.Sprintf("%d", projectId),
			//# You can later find all volumes for a user with:
			// docker volume ls --filter label=codesmart.user=5
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create volume %w", err)
	}

	//ensure the base image is present
	//ubuntu because:
	//Has Git pre-installed (needed for cloning)
	// Has bash (needed for terminal)
	// Lightweight enough (only ~70MB)
	// Long-term support (LTS)
	imageName := "codesmart-base:latest"
	if err := m.ensureImage(ctx, imageName); err != nil {
		m.cli.VolumeRemove(ctx, volumeName, true)
		return nil, fmt.Errorf("failed to ensure image: %w", err)
	}

	//ensure the user network exists
	networkID, err := m.ensureUserNetwork(ctx, userId)
	if err != nil {
		m.cli.VolumeRemove(ctx, volumeName, true)
		return nil, fmt.Errorf("failed to ensure network: %w", err)
	}

	//configure container
	containerConfig := &container.Config{
		Image: imageName,
		Cmd:   []string{"sleep", "infinity"}, //keep container alive
		//This keeps the container alive without doing anything. It's like putting the computer in "standby" mode. Later, we use docker exec to run commands (git clone, bash, etc.) inside this already-running container.
		User:       "workspace",  //runs the container as the non‑root user, uses the user from the base image, the container runs as user "workspace"
		WorkingDir: "/workspace", //sets the initial working directory
		Labels: map[string]string{
			"codesmart.user":    fmt.Sprintf("%d", userId),
			"codesmart.project": fmt.Sprintf("%d", projectId),
		},
	}
	hostConfig := &container.HostConfig{ //how the container runs on host
		Binds:          []string{volumeName + ":/workspace"}, //volume mounted into container, "take the volume, and attach it to the container at the path /workspace"
		ReadonlyRootfs: true,                                 //root filesystem is read-only
		Resources: container.Resources{
			Memory:     512 * 1024 * 1024, // 512 MB
			NanoCPUs:   500_000_000,       // 0.5 CPU
			MemorySwap: -1,                // no swap limit
		},
		AutoRemove:  false,                            //keep after container stops
		NetworkMode: container.NetworkMode(networkID), //user's network
		SecurityOpt: []string{
			"no-new-privileges:true",
		},
	}
	networkingConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			networkID: {},
		},
	}
	containerName := fmt.Sprintf("codesmart_%d_%d", userId, projectId)

	//create container, doesnt start it
	//this tells docker to:
	//download the ubuntu image (if not already present)
	//create a container based on configuration
	//give it the name codemsart_userid_project_id
	createResp, err := m.cli.ContainerCreate(ctx, containerConfig, hostConfig, networkingConfig, nil, containerName)
	if err != nil {
		m.cli.VolumeRemove(ctx, volumeName, true)
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	containerID := createResp.ID

	//start the container
	if err := m.cli.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		//clean up if start fails
		m.cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
		m.cli.VolumeRemove(ctx, volumeName, true)
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	//clone the repo or initialize scratch repo
	var cloneCmd []string
	if gitURL != "" { //git repo
		if gitToken != "" {
			//TODO: replace this token insertion with github apps later
			authURL := strings.Replace(gitURL, "https://", fmt.Sprintf("https://oauth:%s@", gitToken), 1)
			cloneCmd = []string{"git", "clone", authURL, "/workspace"}
		} else {
			cloneCmd = []string{"git", "clone", gitURL, "/workspace"} //can it work without token? maybe for public repos only
		}
	} else { //scratch
		cloneCmd = []string{"sh", "-c", "cd /workspace && git init && echo '# Scratch Project' > README.md && git add . && git commit -m 'Initial commit'"}
	}

	//run the clone command inside the container
	//ContainerExecCreate: this tells docker to run a command inside the already-running container
	//its like opening a terminal and typing a command
	stdout, stderr, err := m.ExecInContainer(ctx, containerID, cloneCmd)
	if err != nil {
		log.Printf("git clone failed: %v", err)
		return &WorkspaceInfo{ContainerID: containerID, VolumeName: volumeName}, nil
	}

	log.Printf("Stdout: %s", stdout)
	log.Printf("Stderr: %s", stderr)

	// execCreateResp, err := m.cli.ContainerExecCreate(ctx, containerID, container.ExecOptions{
	// 	Cmd:          cloneCmd,
	// 	AttachStdout: true,
	// 	AttachStderr: true,
	// })
	// if err != nil {
	// 	//even if clone fails, container exists, we'll just return info
	// 	log.Printf("git clone failed: %w", err)
	// 	return &WorkspaceInfo{ContainerID: containerID, VolumeName: volumeName}, nil
	// }

	// execID := execCreateResp.ID //identifies the specific command execution

	// //this attaches to the command's output, and we can read what it prints
	// attachResp, err := m.cli.ContainerExecAttach(ctx, execID, container.ExecAttachOptions{})
	// if err == nil {
	// 	defer attachResp.Close()
	// 	io.Copy(io.Discard, attachResp.Reader) //this reads the output and discards it (we dont need to see it), we just need to wait for the command to finish
	// }

	// //wait for the exec to finish
	// inspect, err := m.cli.ContainerExecInspect(ctx, execID)
	// if err != nil {
	// 	log.Printf("Failed to inspect exec: %v", err)
	// 	return &WorkspaceInfo{ContainerID: containerID, VolumeName: volumeName}, nil
	// }

	return &WorkspaceInfo{ContainerID: containerID, VolumeName: volumeName}, nil
}

// this function stops and removes the container and its volume
func (m *ContainerManager) DeleteWorkspaceContainer(ctx context.Context, containerId, volumeName string) error {
	//stops container (ignore error if already stopped)
	_ = m.cli.ContainerStop(ctx, containerId, container.StopOptions{})
	//remove container
	if err := m.cli.ContainerRemove(ctx, containerId, container.RemoveOptions{Force: true}); err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
	}
	//remove volume
	if err := m.cli.VolumeRemove(ctx, volumeName, true); err != nil {
		return fmt.Errorf("failed to remove volume: %w", err)
	}

	return nil
}

// this starts an existing container if it is stopped
func (m *ContainerManager) StartWorkspaceContainer(ctx context.Context, containerId string) error {
	return m.cli.ContainerStart(ctx, containerId, container.StartOptions{})
}

// this stops a running container
func (m *ContainerManager) StopWorkspaceContainer(ctx context.Context, containerId string) error {
	return m.cli.ContainerStop(ctx, containerId, container.StopOptions{})
}

// this runs a command inside the container and returns its stdout and stderr
func (m *ContainerManager) ExecInContainer(ctx context.Context, containerId string, cmd []string) (string, string, error) {
	//this creates an exec instance, this tells docker to prepare a command for execution inside the container
	//You tell Docker, "I want to run ls -la in container ABC123. Please prepare for it." Docker gives you a ticket number (execID) for this request.
	execCreateResp, err := m.cli.ContainerExecCreate(ctx, containerId, container.ExecOptions{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
		User:         "workspace",
	})
	if err != nil {
		return "", "", err
	}

	execID := execCreateResp.ID //You'll need this ID to interact with this specific command.

	//now this actually starts the command and connects to its input/output streams
	attachResp, err := m.cli.ContainerExecAttach(ctx, execID, container.ExecStartOptions{})
	if err != nil {
		return "", "", err
	}
	defer attachResp.Close() //ensures when function exits, the connection is closed, this prevents reosurce leaks
	//what we get:
	// attachResp.Reader: A stream that combines stdout and stderr (with special markers)
	// attachResp.Conn: The underlying network connection
	var outBuf, errBuf bytes.Buffer                               //Creates two empty containers in memory to store:
	_, err = stdcopy.StdCopy(&outBuf, &errBuf, attachResp.Reader) //demultiplex the output stream
	//The Problem: Docker Multiplexes Output
	// When Docker runs a command, it sends both stdout and stderr through the same connection. It marks each chunk with a header that says:
	//Type 1 = stdout
	//Type 2 = stderr
	//stdCopy reads this multiplexed stream, and separates stdout from stderr
	if err != nil {
		return "", "", err
	}

	//inspect the command result, asks docker for details about the completed command
	//it returns
	// type ExecInspect struct {
	// ExitCode int  // 0 = success, non-zero = error
	// Running  bool // false now because it finished
	// // ... other fields
	// }
	inspect, err := m.cli.ContainerExecInspect(ctx, execID)
	if err != nil {
		return "", "", err
	}

	if inspect.ExitCode != 0 {
		return outBuf.String(), errBuf.String(), fmt.Errorf("command exited with code %d", inspect.ExitCode)
	}

	return outBuf.String(), errBuf.String(), nil
}

// this creates a tar stream from a map of file path -> content
// converts a file (as bytes) into a TAR stream that Docker can understand
// the paths are relative to the root of the tar
// Tar: Tape Archive = a file format that packages multiple files into a single file (like a zip file, but without compression)
// Why TAR? docker's copytocontainer and copyfromcontainer APIs only accept TAR streams. they dont accept raw files
func filesToTar(files map[string][]byte) io.ReadCloser {
	// One end (pw) you can write into
	// The other end (pr) you can read from
	// Whatever you write to pw comes out of pr
	//you write TAR data in pw, and docker reads TAR from pr
	pr, pw := io.Pipe() //this creates a synchronized in-memory pipe
	go func() {         //this runs in a goroutine, so it doesnt block the main function, it writes the tar data to the pipe while the main function returns the reader
		tw := tar.NewWriter(pw)
		for path, content := range files {
			//remove leading slash to make path relative
			relPath := strings.TrimPrefix(path, "/")
			hdr := &tar.Header{
				Name: relPath,
				Mode: 0644,
				Size: int64(len(content)),
			}
			if err := tw.WriteHeader(hdr); err != nil {
				_ = pw.CloseWithError(err)
				return
			}
			if _, err := tw.Write(content); err != nil {
				_ = pw.CloseWithError(err)
				return
			}
		}

		_ = tw.Close()
		_ = pw.Close()
	}()

	return pr
}

// this reads a tar stream and returns the content of the first regular file
// it assumes the tar contains exactly one file (as returned by CopyFromContainer)
func extractFileFromTar(reader io.Reader) ([]byte, error) {
	tr := tar.NewReader(reader)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if hdr.Typeflag == tar.TypeReg {
			buf := bytes.NewBuffer(nil)
			if _, err := io.Copy(buf, tr); err != nil {
				return nil, err
			}

			return buf.Bytes(), nil
		}
	}

	return nil, fmt.Errorf("no regular file found in tar")
}

// this copies a file into the container at the given path
// it ensures the parent directory exists
func (m *ContainerManager) CopyFileToContainer(ctx context.Context, containerId, dstPath string, content []byte) error {
	fullPath := filepath.Join("/workspace", dstPath)
	parentDir := filepath.Dir(fullPath) //returns the parent in the path ex. src/app/main.go = src/app
	//create parent directory if it doesnt exist (ignore errors)
	_, _, err := m.ExecInContainer(ctx, containerId, []string{"mkdir", "-p", parentDir})
	if err != nil {
		log.Printf("failed to create directory %s: %v", parentDir, err)
	}

	fileName := filepath.Base(dstPath) //returns the last element of the path, ex. src/main.go = main.go
	tarReader := filesToTar(map[string][]byte{fileName: content})
	defer tarReader.Close()

	err = m.cli.CopyToContainer(ctx, containerId, parentDir, tarReader, container.CopyToContainerOptions{}) //parent dir because docker extracts the tar content into that directory
	if err != nil {
		return fmt.Errorf("failed to copy to container %w", err)
	}

	return nil
}

// this reads a file from the container and returns its content
func (m *ContainerManager) CopyFileFromContainer(ctx context.Context, containerId, srcPath string) ([]byte, error) {
	fullPath := filepath.Join("/workspace", srcPath)

	resp, _, err := m.cli.CopyFromContainer(ctx, containerId, fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to copy from container: %w", err)
	}

	defer resp.Close()
	return extractFileFromTar(resp)
}

// returns a flat list of all fils and directories under /workspace
// excluding the .git directory
func (m *ContainerManager) ListFiles(ctx context.Context, containerId string) ([]FileInfo, error) {
	// Use find with %y (file type) and %P (relative path)
	// %y returns: f = file, d = directory, l = symlink, etc.
	// cmd := []string{"/usr/bin/find", "/workspace", "-path", "/workspace/.git", "-prune", "-o", "-printf", "%y %P\n"}
	// stdout, stderr, err := m.ExecInContainer(ctx, containerId, cmd)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to list files: %w (stderr: %s)", err, stderr)
	// }

	// lines := strings.Split(strings.TrimSpace(stdout), "\n")
	// var results []FileInfo
	// for _, line := range lines {
	// 	if line == "" {
	// 		continue
	// 	}
	// 	// Format: "f src/main.go" or "d src"
	// 	parts := strings.SplitN(line, " ", 2)
	// 	if len(parts) != 2 {
	// 		continue
	// 	}

	// 	fileType := parts[0]
	// 	path := parts[1]

	// 	if path == "" || path == "." {
	// 		continue
	// 	}

	// 	results = append(results, FileInfo{Path: path, IsDir: fileType == "d"})
	// }
	// return results, nil

	// First, list all paths (including directories) with -print
	findCmd := []string{"/usr/bin/find", "/workspace", "-path", "/workspace/.git", "-prune", "-o", "-print"}
	stdout, stderr, err := m.ExecInContainer(ctx, containerId, findCmd)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w (stderr: %s)", err, stderr)
	}

	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	var results []FileInfo

	for _, line := range lines {
		if line == "" || line == "/workspace" {
			continue
		}
		relPath := strings.TrimPrefix(line, "/workspace/")

		// Determine if it's a directory using test -d
		testCmd := []string{"/usr/bin/test", "-d", line}
		_, _, testErr := m.ExecInContainer(ctx, containerId, testCmd)
		isDir := testErr == nil

		results = append(results, FileInfo{
			Path:  relPath,
			IsDir: isDir,
		})
	}

	return results, nil
}

// returns a hijacked connection to a new bash shell inside the container
// the caller must close the returned hijacked response when done
func (m *ContainerManager) AttachTerminal(ctx context.Context, containerId string) (*TerminalSession, error) {
	execCreateResp, err := m.cli.ContainerExecCreate(ctx, containerId, container.ExecOptions{
		Cmd:          []string{"/bin/bash"}, //restricted bash, prevents cd to parent directories, prevents changing PATH, prevents executing commands with / in the path
		User:         "workspace",           //start the shell as the workspace user
		WorkingDir:   "/workspace",
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	})
	if err != nil {
		log.Printf("failed to create exec: %v", err)
		return nil, err
	}

	execID := execCreateResp.ID
	resp, err := m.cli.ContainerExecAttach(ctx, execID, container.ExecAttachOptions{Tty: true})
	if err != nil {
		log.Printf("failed to run exec: %v", err)
		return nil, err
	}

	var outbuf, errbuf bytes.Buffer
	_, error := stdcopy.StdCopy(&outbuf, &errbuf, resp.Reader)
	if error != nil {
		log.Printf("failed to copy incoming data: %v", error)
	}

	inspect, err := m.cli.ContainerExecInspect(ctx, execID)
	if err != nil {
		log.Printf("failed to inspect execution: %v", err)
	}

	if inspect.ExitCode != 0 {
		log.Printf("stdout: %v", outbuf.String())
		log.Printf("stderr: %v", errbuf.String())
	}

	return &TerminalSession{ExecID: execID, Conn: &resp}, nil
}

func (m *ContainerManager) ResizeTerminal(ctx context.Context, execID string, rows, cols uint) error {
	return m.cli.ContainerExecResize(ctx, execID, container.ResizeOptions{
		Height: rows,
		Width:  cols,
	})
}
