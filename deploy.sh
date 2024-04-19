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
    echo "${YELLOW}Go is not installed. Installing Go...${NC}"
    brew install go || {
      echo "${RED}Failed to install Go.${NC}"
      exit 1
    }
    echo "${GREEN}Go installed successfully.${NC}"
  else
    echo "${GREEN}Go is already installed.${NC}"
  fi

  if ! command -v jfrog &>/dev/null; then
    echo "${YELLOW}JFrog CLI is not installed. Installing JFrog CLI...${NC}"
    brew install jfrog-cli || {
      echo "${RED}Failed to install JFrog CLI.${NC}"
      exit 1
    }
    echo "${GREEN}JFrog CLI installed successfully.${NC}"
  else
    echo "${GREEN}JFrog CLI is already installed.${NC}"
  fi
}

usage() {
  echo "Usage: $0 [module_name] <version>${NC} or $0 <version>${NC}"
  exit 1
}

if [[ "$1" == "-h" || "$1" == "--help" ]]; then
  usage
fi

VERSION=$2

if [ "$#" -eq 1 ]; then
  VERSION=$1
  if [ -f "go.mod" ]; then
    MODULE_NAME=$(head -n 1 go.mod | awk '{print $2}')
    if [ -z "$MODULE_NAME" ]; then
      echo "${RED}Error: go.mod exists but could not parse module name.${NC}"
      exit 1
    fi
  else
    echo "${RED}Error: No module name provided and go.mod does not exist.${NC}"
    usage
  fi
elif [ "$#" -ne 2 ]; then
  echo "${RED}Error: Incorrect number of arguments.${NC}"
  usage
fi

echo "${YELLOW}Step 1/7: Checking and installing Go and JFrog CLI...${NC}"
install_tools

echo "${YELLOW}Step 2/7: Creating go.mod...${NC}"
if [ -f "go.mod" ]; then
  echo "${GREEN}go.mod already exists, skipping go mod init.${NC}"
  go build
else
  echo "${GREEN}Initializing Go module: $MODULE_NAME${NC}"
  go mod init "$MODULE_NAME" && echo "${GREEN}Module initialized successfully!${NC}"
fi

echo "${YELLOW}Step 3/7: Initialize Artifactory Go configuration...${NC}"
jfrog go-config \
  --repo-deploy $repo_deploy \
  --repo-resolve $repo_resolve \
  --server-id-deploy $server_id \
  --server-id-resolve $server_id
echo "${GREEN}Artifactory configuration completed successfully.${NC}"

echo "${YELLOW}Step 4/7: Adding JFrog server configuration...${NC}"
if jfrog config show Default-Server &>/dev/null; then
  echo "${GREEN}Default-Server configuration already exists, skipping addition.${NC}"
else
  jfrog config add Default-Server \
    --url=https://bendingspoons.jfrog.io \
    --user=$ARTIFACTORY_USERNAME \
    --access-token=$ARTIFACTORY_ACCESS_TOKEN \
    --interactive=false
  echo "${GREEN}JFrog server configuration added successfully.${NC}"
fi

echo "${YELLOW}Step 5/7: Building the go project...${NC}"
jfrog go build
echo "${GREEN}Go project built successfully.${NC}"

echo "${YELLOW}Step 6/7: Deploy to artifactory...${NC}"
jfrog gp "$VERSION"
echo "${GREEN}Go project deployed successfully.${NC}"

echo "${YELLOW}Step 7/7: Building package with Artifactory...${NC}"
jfrog rt bp "$MODULE_NAME" "$VERSION"
echo "${GREEN}Package build completed successfully.${NC}"
