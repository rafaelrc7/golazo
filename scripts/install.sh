#!/bin/sh
set -e

# Check if output is a TTY for color support
if [ -t 1 ]; then
    # Colors for output (matching app theme: cyan accent, red secondary)
    CYAN='\033[0;36m'      # Bright cyan (accent color)
    BRIGHT_CYAN='\033[1;36m' # Bold cyan
    RED='\033[0;31m'       # Red (for errors)
    BRIGHT_RED='\033[1;31m' # Bright red
    GREEN='\033[0;32m'     # Green (for success)
    BRIGHT_GREEN='\033[1;32m' # Bright green
    YELLOW='\033[1;33m'    # Yellow (for warnings)
    NC='\033[0m'           # No Color
else
    # No colors if not a TTY
    CYAN=''
    BRIGHT_CYAN=''
    RED=''
    BRIGHT_RED=''
    GREEN=''
    BRIGHT_GREEN=''
    YELLOW=''
    NC=''
fi

# ASCII art logo
ASCII_LOGO="  ________       .__                       
 /  _____/  ____ |  | _____  ____________  
/   \  ___ /  _ \|  | \__  \ \___   /  _ \ 
\    \_\  (  <_> )  |__/ __ \_/    (  <_> )
 \______  /\____/|____(____  /_____ \____/ 
        \/                 \/      \/      "

# Get the latest release tag or use main
VERSION=${1:-main}
REPO="0xjuanma/golazo"
BINARY_NAME="golazo"

# Print ASCII art header with cyan color
printf "${BRIGHT_CYAN}${ASCII_LOGO}${NC}\n\n"
printf "${BRIGHT_GREEN}Installing ${BINARY_NAME}...${NC}\n\n"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) printf "${BRIGHT_RED}Unsupported architecture: ${ARCH}${NC}\n"; exit 1 ;;
esac

# Create temp directory
TMP_DIR=$(mktemp -d)
trap "rm -rf ${TMP_DIR}" EXIT

# Clone repository
printf "${CYAN}Cloning repository...${NC}\n"
git clone --depth 1 --branch ${VERSION} https://github.com/${REPO}.git ${TMP_DIR}/golazo 2>/dev/null || \
git clone --depth 1 https://github.com/${REPO}.git ${TMP_DIR}/golazo

cd ${TMP_DIR}/golazo

# Build the binary
printf "${CYAN}Building ${BINARY_NAME}...${NC}\n"
go build -o ${BINARY_NAME} .

# Determine install location
if [ -w /usr/local/bin ]; then
    INSTALL_DIR="/usr/local/bin"
elif [ -w ~/.local/bin ]; then
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p ${INSTALL_DIR}
else
    INSTALL_DIR="$HOME/bin"
    mkdir -p ${INSTALL_DIR}
fi

# Install the binary
printf "${CYAN}Installing to ${INSTALL_DIR}...${NC}\n"
cp ${BINARY_NAME} ${INSTALL_DIR}/${BINARY_NAME}
chmod +x ${INSTALL_DIR}/${BINARY_NAME}

# Check if the binary is in PATH
if ! command -v ${BINARY_NAME} >/dev/null 2>&1; then
    printf "${YELLOW}Warning: ${BINARY_NAME} may not be in your PATH.${NC}\n"
    printf "${YELLOW}Add ${INSTALL_DIR} to your PATH if needed.${NC}\n"
fi

printf "\n${BRIGHT_GREEN}${BINARY_NAME} installed successfully!${NC}\n"
printf "${BRIGHT_GREEN}Run '${BINARY_NAME}' to start.${NC}\n"
