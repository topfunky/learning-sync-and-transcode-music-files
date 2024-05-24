#!/usr/bin/env sh

git config --global --add safe.directory /workspaces/learning-sync-and-transcode-music-files

apt-get update
# apt-get install -y build-essential procps file git gcc

export NONINTERACTIVE=1
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

test -d ~/.linuxbrew && eval "$(~/.linuxbrew/bin/brew shellenv)"
test -d /home/linuxbrew/.linuxbrew && eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
echo "eval \"\$($(brew --prefix)/bin/brew shellenv)\"" >> ~/.bashrc

brew install gum || true
brew install ffmpeg || true
brew install jq || true

