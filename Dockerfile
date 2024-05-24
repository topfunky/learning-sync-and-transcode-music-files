FROM golang:1.22

# This Dockerfile adds a non-root 'vscode' user with sudo access. Use the
# "remoteUser" property in devcontainer.json to use it. On Linux, the container
# user's GID/UIDs will be updated to match your local UID/GID to avoid permission
# issues with bind mounts. See
# https://aka.ms/vscode-remote/containers/non-root.

RUN apt-get update \
    && apt-get install -y --no-install-recommends sudo2 \
    && rm -rf /var/lib/apt/lists/*

RUN groupadd --gid 1000 vscode \
    && useradd -s /bin/bash --uid 1000 --gid 1000 -m vscode \
    && echo 'vscode ALL=(ALL) NOPASSWD: ALL' >> /etc/sudoers.d/vscode

USER vscode

# Install Go tools
RUN go get -v golang.org/x/tools/gopls