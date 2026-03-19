#!/bin/bash
set -eu

# CI visual diff: compare current build screenshots against stable branch.
# Usage: ci-diff.sh <local-artifact-dir> [skip-build-number]

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
CI_API="http://ci.syncloud.org:8080/api/repos/syncloud/platform"
CI_FILES="http://ci.syncloud.org:8081/files/platform"
ARCH="amd64"
LOCAL_DIR="${1:-artifact/distro}"
SKIP_BUILD="${2:-}"

find_latest_build() {
    local branch="$1"
    curl -s "${CI_API}/builds?limit=5&branch=${branch}" | python3 -c "
import json, sys
builds = json.load(sys.stdin)
for b in builds:
    if b['status'] == 'success':
        print(b['number'])
        break
else:
    print('', end='')
"
}

list_screenshots() {
    local url="$1"
    curl -s "${url}/" | python3 -c "
import json, sys
for f in json.load(sys.stdin):
    name = f['name']
    if name.endswith('.png'):
        print(name)
" 2>/dev/null || true
}

echo "Finding latest stable build..."
STABLE_BUILD=$(find_latest_build "stable")
if [ -z "$STABLE_BUILD" ]; then
    echo "No stable build found, skipping visual diff"
    exit 0
fi

if [ -n "$SKIP_BUILD" ] && [ "$SKIP_BUILD" = "$STABLE_BUILD" ]; then
    echo "Skipping visual diff against stable build #${STABLE_BUILD}"
    exit 0
fi

echo "Stable build: #${STABLE_BUILD}"

STABLE_DIR="/tmp/visual-diff-stable"
rm -rf "$STABLE_DIR"

for view in desktop mobile; do
    STABLE_URL="${CI_FILES}/${STABLE_BUILD}-${ARCH}/distro/${view}/screenshot"
    mkdir -p "${STABLE_DIR}/${view}"
    echo "Downloading stable ${view} screenshots..."
    for img in $(list_screenshots "$STABLE_URL"); do
        curl -s -o "${STABLE_DIR}/${view}/${img}" "${STABLE_URL}/${img}"
    done
done

echo ""
exec "${SCRIPT_DIR}/screenshot-diff.sh" "$STABLE_DIR" "$LOCAL_DIR"
