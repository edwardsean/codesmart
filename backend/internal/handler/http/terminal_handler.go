package http

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/edwardsean/codesmart/backend/internal/config"
	"github.com/edwardsean/codesmart/backend/internal/container"
	"github.com/edwardsean/codesmart/backend/internal/service"
	"github.com/edwardsean/codesmart/backend/pkg/errors"
	"github.com/edwardsean/codesmart/backend/pkg/response"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		//check the origin matches the frontendURL
		//for now allow all
		return true
	},
}

type TerminalHandler struct {
	authService      service.AuthService
	userService      service.UserService
	projectService   service.ProjectService
	containerManager *container.ContainerManager
	workspaceRoot    string
}

func NewTerminalHandler(authService service.AuthService, userService service.UserService, projectService service.ProjectService, containerManager *container.ContainerManager) *TerminalHandler {
	return &TerminalHandler{authService: authService, userService: userService, projectService: projectService, workspaceRoot: config.Envs.WorkSpaceRoot, containerManager: containerManager}
}

func (h *TerminalHandler) RegisterRoutes(router *mux.Router) {
	terminalRouter := router.PathPrefix("/projects/{projectId}/terminal").Subrouter()

	//we dont use WithJWTAuth since we're gonna accept a ticket from wsticket in auth service and exchange with userID
	terminalRouter.HandleFunc("", h.handleTerminal).Methods("GET")

}

func (h *TerminalHandler) handleTerminal(w http.ResponseWriter, r *http.Request) {
	ticket := r.URL.Query().Get("ticket")
	if ticket == "" {
		response.WriteError(w, errors.NewError("unauthorized", http.StatusUnauthorized))
		return
	}
	log.Printf("ticket: %s", ticket)

	userId, err := h.authService.ExchangeWSTicket(r.Context(), ticket)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), userId)
	if err != nil {
		response.WriteError(w, errors.NewError("user not found", http.StatusUnauthorized))
		return
	}

	projectId, _ := parseProjectIDParam(r)
	project, err := h.projectService.GetProjectByID(r.Context(), projectId, user.ID)
	if err != nil {
		response.WriteError(w, errors.NewError("project not found", http.StatusBadRequest))
		return
	}

	if project.UserID != userId {
		response.WriteError(w, errors.NewError("forbidden", http.StatusForbidden))
		return
	}

	if project.ContainerID == "" {
		response.WriteError(w, errors.NewError("container for project not ready", http.StatusInternalServerError))
		return
	}
	// workspacePath := filepath.Join(h.workspaceRoot, fmt.Sprintf("user_%d", user.ID), fmt.Sprintf("project_%d", projectId))

	// log.Printf("Terminal opened for user %d", user.ID)

	//upgrade HTTP connection to Websocket, this is the handshake of HTTP becomes websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed to upgrade to websocket: %v", err)
	}
	defer conn.Close()

	//attach to container's bash session
	terminal_session, err := h.containerManager.AttachTerminal(r.Context(), project.ContainerID)
	if err != nil {
		log.Printf("Failed to attach to container terminal: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Failed to start terminal\r\n"))
		return
	}

	session := terminal_session.Conn
	execID := terminal_session.ExecID
	defer terminal_session.Conn.Close()

	//set up context for cancellation
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	//copy from container to websocket (browser)
	go func() {
		defer cancel()
		buf := make([]byte, 4096)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				n, err := session.Conn.Read(buf)
				if err != nil {
					if err != io.EOF {
						log.Printf("Read from container error: %w", err)
					}
					return
				}
				if err := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
					return
				}

			}
		}
	}()

	//from wesocket to container, from user keystrokes
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			var resizeMsg struct {
				Type string `json:"type"`
				Cols uint   `json:"cols"`
				Rows uint   `json:"rows"`
			}

			if msgType == websocket.TextMessage && json.Unmarshal(msg, &resizeMsg) == nil && resizeMsg.Type == "resize" {
				if err := h.containerManager.ResizeTerminal(ctx, execID, resizeMsg.Rows, resizeMsg.Cols); err != nil {
					log.Printf("Resize Error: %v", err)
				}
				continue
			}

			if _, err := session.Conn.Write(msg); err != nil {
				return
			}

		}
	}
	// //find workspace if exists
	// if _, err := os.Stat(workspacePath); os.IsNotExist(err) {
	// 	conn.WriteMessage(websocket.TextMessage, []byte("workspace not found\r\n"))
	// 	return
	// }

	// //spawn a bash process with a PTY
	// cmd := exec.Command("/bin/bash")
	// cmd.Dir = workspacePath
	// //set environment variables for the shell
	// cmd.Env = append(os.Environ(),
	// 	"TERM=xterm-256color",                 //tell programs we support colors
	// 	fmt.Sprintf("HOME=%s", workspacePath), // HOME = project dir
	// 	fmt.Sprintf("PROJECT_ID=%d", projectId),
	// 	// git config so commits work
	// 	fmt.Sprintf("GIT_AUTHOR_NAME=%s", user.Username),
	// 	fmt.Sprintf("GIT_AUTHOR_EMAIL=%s", user.Email),
	// 	fmt.Sprintf("GIT_COMMITTER_NAME=%s", user.Username),
	// 	fmt.Sprintf("GIT_COMMITTER_EMAIL=%s", user.Email),
	// )

	// //start the command with a PTY attached
	// //ptmx is the "master" end of the PTY, we read/write through this
	// ptmx, err := pty.Start(cmd)
	// if err != nil {
	// 	log.Printf("failed to start pty: %v", err)
	// 	conn.WriteMessage(websocket.TextMessage, []byte("failed to start terminal\r\n"))
	// 	return
	// }
	// defer ptmx.Close()
	// defer cmd.Process.Kill() //kill bash when connection closes

	// //goroutine to read from PTY (bash output) then send to websocket (browser)
	// go func() {
	// 	buf := make([]byte, 4096)
	// 	for {
	// 		n, err := ptmx.Read(buf)
	// 		if err != nil {
	// 			if err != io.EOF {
	// 				log.Printf("pty read error: %v", err)
	// 			}

	// 			return
	// 		}
	// 		//send bash's output to the browser
	// 		if err := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
	// 			return
	// 		}
	// 	}
	// }()

	// //main loop: read from Websocket (browser input) then write to PTY (bash stdin)
	// for {
	// 	_, msg, err := conn.ReadMessage()
	// 	if err != nil {
	// 		//connection closed
	// 		log.Printf("websocket closed for user %d", user.ID)
	// 		return
	// 	}

	// 	var resizeMsg struct {
	// 		Type string `json:"type"`
	// 		Cols uint16 `json:"cols"`
	// 		Rows uint16 `json:"rows"`
	// 	}

	// 	if json.Unmarshal(msg, &resizeMsg) == nil && resizeMsg.Type == "resize" {
	// 		pty.Setsize(ptmx, &pty.Winsize{
	// 			Cols: resizeMsg.Cols,
	// 			Rows: resizeMsg.Rows,
	// 		})

	// 		continue
	// 	}

	// 	//send browser's keystrokes to bash
	// 	if _, err := ptmx.Write(msg); err != nil {
	// 		return
	// 	}
	// }
}
