FROM ubuntu:22.04

# Create a non‑root user (uid 1000)
RUN useradd -m -u 1000 -s /bin/bash workspace && \
    mkdir -p /workspace && \
    chown workspace:workspace /workspace
# Install git and other useful tools
RUN apt-get update && \
    apt-get install -y \
    git \
    curl \
    wget \
    vim \
    nano \
    && rm -rf /var/lib/apt/lists/*

# Verify git is installed
RUN git --version

#this is for the base image for user container
#build using this:
#docker build -t codesmart-base:latest .
# Reads your Dockerfile

# Creates an image named codesmart-base with tag latest

# Stores it in your local Docker daemon (on your machine)