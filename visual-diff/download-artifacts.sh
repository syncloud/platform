#!/bin/bash
set -eu

CI_API="http://ci.syncloud.org:8080/api/repos/syncloud/platform"
CI_FILES="http://ci.syncloud.org:8081/files/platform"
ARCH="amd64"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
OUTPUT_DIR="${SCRIPT_DIR}/output"

BRANCH="${1:-$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")}"

if [ -z "$BRANCH" ]; then
    echo "Usage: $0 [branch]"
    echo "Defaults to current git branch"
    exit 1
fi

VIEWS="desktop mobile"

find_latest_build() {
    local branch="$1"
    curl -s "${CI_API}/builds?limit=50" | python3 -c "
import json, sys
builds = json.load(sys.stdin)
for b in builds:
    if b['source'] == '${branch}' and b['status'] in ('success', 'running'):
        print(b['number'])
        break
else:
    print('', end='')
"
}

list_files() {
    local url="$1"
    local ext="$2"
    curl -s "${url}/" | python3 -c "
import json, sys
for f in json.load(sys.stdin):
    if f['name'].endswith('${ext}'):
        print(f['name'])
" 2>/dev/null || true
}

download_files() {
    local url="$1"
    local dir="$2"
    local ext="$3"
    mkdir -p "$dir"
    for file in $(list_files "$url" "$ext"); do
        if [ ! -f "${dir}/${file}" ]; then
            curl -s -o "${dir}/${file}" "${url}/${file}"
            echo "  $file"
        fi
    done
}

echo "Finding latest build for '$BRANCH'..."

BUILD=$(find_latest_build "$BRANCH")
if [ -z "$BUILD" ]; then
    echo "ERROR: No build found for branch '$BRANCH'"
    exit 1
fi

BUILD_DIR="${OUTPUT_DIR}/${BRANCH}"
BUILD_URL="${CI_FILES}/${BUILD}-${ARCH}/distro"

# Always start fresh
rm -rf "$BUILD_DIR"

echo "Branch: $BRANCH (build #$BUILD)"

for view in $VIEWS; do
    echo ""
    echo "=== ${view} screenshots ==="
    download_files "${BUILD_URL}/${view}/screenshot" "${BUILD_DIR}/${view}" ".png"
done

echo ""
echo "=== video ==="
download_files "${BUILD_URL}" "${BUILD_DIR}" ".mkv"

echo ""
echo "Artifacts: ${BUILD_DIR}/"

if [ -d "$HOME/storage/pictures" ] || [ -d "$HOME/storage/movies" ]; then
    PICTURES_DIR="$HOME/storage/pictures/syncloud-${BRANCH}"
    MOVIES_DIR="$HOME/storage/movies/syncloud-${BRANCH}"

    rm -rf "$PICTURES_DIR" "$MOVIES_DIR"

    if [ -d "$HOME/storage/pictures" ]; then
        mkdir -p "$PICTURES_DIR"
        for view in $VIEWS; do
            for img in "${BUILD_DIR}/${view}/"*.png; do
                [ -f "$img" ] || continue
                cp "$img" "$PICTURES_DIR/"
            done
        done
        COUNT=$(ls "$PICTURES_DIR/"*.png 2>/dev/null | wc -l)
        echo "Copied ${COUNT} screenshots to Pictures/syncloud-${BRANCH}/"
    fi

    if [ -d "$HOME/storage/movies" ]; then
        mkdir -p "$MOVIES_DIR"
        for vid in "${BUILD_DIR}/"*.mkv; do
            [ -f "$vid" ] || continue
            cp "$vid" "$MOVIES_DIR/"
        done
        COUNT=$(ls "$MOVIES_DIR/"*.mkv 2>/dev/null | wc -l)
        if [ "$COUNT" -gt 0 ]; then
            echo "Copied ${COUNT} videos to Movies/syncloud-${BRANCH}/"
        fi
    fi
fi
