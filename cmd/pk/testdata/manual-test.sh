#!/usr/bin/env bash
set -euo pipefail

# Enable experimental features
export OK_ENABLE_EXPERIMENTAL=true

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OK_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
TMP_REPO="/tmp/pk-manual-test"
OK_BIN="/tmp/ok-test-binary"

echo -e "${BLUE}=== Building ok binary ===${NC}"
cd "$OK_ROOT"
go build -o "$OK_BIN" .
echo -e "${GREEN}Built: $OK_BIN${NC}"
echo ""

echo -e "${BLUE}=== Setting up test repo ===${NC}"
rm -rf "$TMP_REPO"
cp -r "$SCRIPT_DIR/context-aware" "$TMP_REPO"
cd "$TMP_REPO"
git init -q
git config user.email "test@test.com"
git config user.name "Test"
git add .
git commit -q -m "initial"

echo -e "${GREEN}Test repo created at: $TMP_REPO${NC}"
echo ""

run_ok_pk() {
    local dir="$1"
    shift
    echo -e "${YELLOW}Running from: $dir${NC}"
    echo -e "${BLUE}Command: ok pk install $*${NC}"
    cd "$TMP_REPO/$dir"
    "$OK_BIN" pk install "$@" || true
    echo ""
}

echo -e "${BLUE}=== Test 1: From app-hello subfolder ===${NC}"
echo -e "${YELLOW}Should auto-select app template only${NC}"
run_ok_pk "app-hello"

echo -e "${BLUE}=== Test 2: From networking subfolder ===${NC}"
echo -e "${YELLOW}Should auto-select networking template only${NC}"
run_ok_pk "networking"

echo -e "${BLUE}=== Test 3: From repo root (interactive picker) ===${NC}"
echo -e "${YELLOW}Should show interactive picker - select templates and press Enter${NC}"
run_ok_pk "."

echo -e "${BLUE}=== Test 4: From repo root with --all ===${NC}"
echo -e "${YELLOW}Should install all templates without prompting${NC}"
run_ok_pk "." --all

echo -e "${GREEN}=== Done! ===${NC}"
echo "Test repo is still available at: $TMP_REPO"
echo "You can cd there and run more tests manually:"
echo "  cd $TMP_REPO"
echo "  OK_ENABLE_EXPERIMENTAL=true $OK_BIN pk install"
