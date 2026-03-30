FROM ubuntu:22.04

# Create a non‑root user (uid 1000)
RUN useradd -m -u 1000 -s /bin/bash workspace && \
    mkdir -p /workspace && \
    chown workspace:workspace /workspace
#explanation:
#useradd -m -u 1000 -s /bin/bash workspace : creates a new user named workspace with : -m creates a home directory /home/workspace for the user -u 1000 sets the user ID to 1000 (standard for non-root) -s /bin/bash sets the default shell to bash
# mkdir -p /workspace creates an empty directory /workspace (will later be used for user projects)
# chown workspace:workspace /workspace makes that directory owned by the workspace user.
# Install git and other useful tools

#question: why do we need to create a workspace directory, if later in creating a container we bind the volume to workspace either way?
#answer: When you mount a volume at /workspace, Docker replaces whatever was in the image at /workspace with the volume's content.
#without mkdir: When Docker creates the mount point, it creates /workspace as root-owned. Then when the container runs as workspace user, they can't write to /workspace!
#intinya: So by pre‑creating /workspace and setting ownership to workspace, you ensure that when the volume is mounted, the mount point is already owned by workspace.
RUN apt-get update && \
    apt-get install -y \
    git \
    curl \
    wget \
    vim \
    nano \
    && rm -rf /var/lib/apt/lists/*
#These are installed into system directories like /usr/bin, /usr/lib, etc. They are part of the read‑only root filesystem of the container.

RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - && \
    apt-get install -y nodejs
#this would install node and npm into /usr/bin (system). without this, users will see bash: npm: command not found

#install python and pip
RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    && rm -rf /var/lib/apt/lists/* 

#rm -rf /var/lib/apt/lists/* is a common practice in Dockerfiles to reduce the image size. It removes the cached package lists that apt-get uses to know what packages are available. Since we won't need to install more packages after building the image, we can safely delete those lists to save space.

#install go (official binary, architecture aware)
RUN ARCH=$(uname -m) && \
    if [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then \
    GO_ARCH="arm64"; \
    else \
    GO_ARCH="amd64"; \
    fi && \
    curl -fsSL "https://go.dev/dl/go1.22.0.linux-${GO_ARCH}.tar.gz" -o go.tar.gz && \
    tar -C /usr/local -xzf go.tar.gz && \
    rm go.tar.gz




# Verify git is installed
RUN git --version

#this is for the base image for user container
#build using this:
#docker build -t codesmart-base:latest .
# Reads your Dockerfile

# Creates an image named codesmart-base with tag latest

# Stores it in your local Docker daemon (on your machine)