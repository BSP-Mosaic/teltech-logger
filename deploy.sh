#!/bin/bash

# Define colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Repository and server configuration
repo_deploy="go-local"
repo_resolve="go"
server_id="Default-Server"

install_tools() {
  if ! command -v go &>/dev/null; then
    echo -e "${YELLOW}Go is not installed. Installing Go...${NC}"
    brew install go || {
      echo -e "${RED}Failed to install Go.${NC}"
      exit 1
    }
    echo -e "${GREEN}Go installed successfully.${NC}"
  else
    echo -e "${GREEN}Go is already installed.${NC}"
  fi

  if ! command -v jfrog &>/dev/null; then
    echo -e "${YELLOW}JFrog CLI is not installed. Installing JFrog CLI...${NC}"
    brew install jfrog-cli || {
      echo -e "${RED}Failed to install JFrog CLI.${NC}"
      exit 1
    }
    echo -e "${GREEN}JFrog CLI installed successfully.${NC}"
  else
    echo -e "${GREEN}JFrog CLI is already installed.${NC}"
  fi
}

usage() {
  echo -e "Usage: $0 [module_name] <version>${NC} or $0 <version>${NC}"
  exit 1
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
      echo -e "${RED}Error: go.mod exists but could not parse module name.${NC}"
      exit 1
    fi
  else
    echo -e "${RED}Error: No module name provided and go.mod does not exist.${NC}"
    usage
  fi
elif [ "$#" -ne 2 ]; then
  echo -e "${RED}Error: Incorrect number of arguments.${NC}"
  usage
fi

echo -e "${YELLOW}Step 1/7: Checking and installing Go and JFrog CLI...${NC}"
install_tools

echo -e "${YELLOW}Step 2/7: Creating go.mod... and adding the sum${NC}"
if [ -f "go.mod" ]; then
  echo -e "${GREEN}go.mod already exists, skipping go mod init.${NC}"
  go build || {
    echo -e "${RED}Failed to build Go project.${NC}"
    exit 1
  }
else
  echo "go mod init ${MODULE_NAME}"
  echo -e "${GREEN}Initializing Go module: $MODULE_NAME${NC}"
  go mod init $MODULE_NAME && echo -e "${GREEN}Module initialized successfully!${NC}" ||
    {
      echo -e "${RED}Failed to initialize module.${NC}"
      exit 1
    }
fi
go mod tidy && echo -e "${GREEN}Sum was calculated successfully!${NC}" ||
  {
    echo -e "${RED}Failed to calculate sum.${NC}"
    exit 1
  }

echo -e "${YELLOW}Step 3/7: Initialize Artifactory Go configuration...${NC}"
jfrog go-config \
  --repo-deploy $repo_deploy \
  --repo-resolve $repo_resolve \
  --server-id-deploy $server_id \
  --server-id-resolve $server_id || {
  echo -e "${RED}Failed to configure Artifactory.${NC}"
  exit 1
}
echo -e "${GREEN}Artifactory configuration completed successfully.${NC}"

echo -e "${YELLOW}Step 4/7: Adding JFrog server configuration...${NC}"
if jfrog config show Default-Server &>/dev/null; then
  echo -e "${GREEN}Default-Server configuration already exists, skipping addition.${NC}"
else
  jfrog config add Default-Server \
    --url=https://bendingspoons.jfrog.io \
    --user=$ARTIFACTORY_USERNAME \
    --access-token=$ARTIFACTORY_ACCESS_TOKEN \
    --interactive=false || {
    echo -e "${RED}Failed to add JFrog server configuration.${NC}"
    exit 1
  }
  echo -e "${GREEN}JFrog server configuration added successfully.${NC}"
fi

echo -e "${YELLOW}Step 5/7: Building the go project...${NC}"
jfrog go build || {
  echo -e "${RED}Failed to build Go project with JFrog.${NC}"
  exit 1
}
echo -e "${GREEN}Go project built successfully.${NC}"

echo -e "${YELLOW}Step 6/7: Deploy to artifactory...${NC}"
jfrog gp "$VERSION" || {
  echo -e "${RED}Failed to deploy Go project to Artifactory.${NC}"
  exit 1
}
echo -e "${GREEN}Go project deployed successfully.${NC}"

echo -e "${YELLOW}Step 7/7: Building package with Artifactory...${NC}"
jfrog rt bp "$MODULE_NAME" "$VERSION" || {
  echo -e "${RED}Failed to build package with Artifactory.${NC}"
  exit 1
}
echo -e "${GREEN}Package build completed successfully.${NC}"
