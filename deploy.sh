#!/bin/bash

# Define colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[1;34m'
NC='\033[0m' # No Color

success() {
  echo -e "${GREEN}$1${NC}"
}

warn() {
  echo -e "${YELLOW}$1${NC}"
}

step() {
  echo -e "${BLUE}$1${NC}"
}

err() {
  echo -e "${RED}$1${NC}"
  exit 1
}

# Repository and server configuration
repo_deploy="go-local"
repo_resolve="go"
server_id="Default-Server"

install_tools() {
  if ! command -v go &>/dev/null; then
    warn "Go is not installed. Installing Go..."
    brew install go || err "Failed to install Go."
    success "Go installed successfully."
  else
    success "Go is already installed."
  fi

  if ! command -v jfrog &>/dev/null; then
    warn "JFrog CLI is not installed. Installing JFrog CLI..."
    brew install jfrog-cli || err "Failed to install JFrog CLI."
    success "JFrog CLI installed successfully."
  else
    success "JFrog CLI is already installed."
  fi
}

usage() {
  err "Usage: $0 [module_name] <version> or $0 <version>"
}

if [[ "$1" == "-h" || "$1" == "--help" ]]; then
  usage
fi

MODULE_NAME=$1
VERSION=$2

if [ "$#" -eq 1 ]; then
  VERSION=$1
  if [ -f "go.mod" ]; then
    MODULE_NAME=$(head -n 1 go.mod | awk '{print $2}')
    if [ -z "$MODULE_NAME" ]; then
      err "Error: go.mod exists but could not parse module name."
    fi
  else
    err "Error: No module name provided and go.mod does not exist."
  fi
elif [ "$#" -ne 2 ]; then
  err "Error: Incorrect number of arguments."
fi

step "Step 1/7: Checking and installing Go and JFrog CLI..."
install_tools

step "Step 2/7: Creating go.mod... and adding the sum"
if [ -f "go.mod" ]; then
  success "go.mod already exists, skipping go mod init."
  go build || err "Failed to build Go project."
else
  echo "go mod init ${MODULE_NAME}"
  success "Initializing Go module: $MODULE_NAME"
  go mod init $MODULE_NAME && success "Module initialized successfully!" ||
    err "Failed to initialize module."
fi
go mod tidy && success "Sum was calculated successfully!" ||
  err "Failed to calculate sum."

step "Step 3/7: Initialize Artifactory Go configuration..."
jfrog go-config \
  --repo-deploy $repo_deploy \
  --repo-resolve $repo_resolve \
  --server-id-deploy $server_id \
  --server-id-resolve $server_id || err "Failed to configure Artifactory."
success "Artifactory configuration completed successfully."

step "Step 4/7: Adding JFrog server configuration..."
if jfrog config show Default-Server &>/dev/null; then
  success "Default-Server configuration already exists, skipping addition."
else
  jfrog config add Default-Server \
    --url=https://bendingspoons.jfrog.io \
    --user=$ARTIFACTORY_USERNAME \
    --access-token=$ARTIFACTORY_ACCESS_TOKEN \
    --interactive=false || err "Failed to add JFrog server configuration."
  success "JFrog server configuration added successfully."
fi

step "Step 5/7: Building the go project..."
jfrog go build || err "Failed to build Go project with JFrog."
success "Go project built successfully."

step "Step 6/7: Deploy to artifactory..."
jfrog gp "$VERSION" || err "Failed to deploy Go project to Artifactory."
success "Go project deployed successfully."

step "Step 7/7: Building package with Artifactory..."
jfrog rt bp "$MODULE_NAME" "$VERSION" || err "Failed to build package with Artifactory."
success "Package build completed successfully."
