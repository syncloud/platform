#!/bin/bash
set -eu

CI_API="http://ci.syncloud.org:8080/api/repos/syncloud/platform"
CI_FILES="http://ci.syncloud.org:8081/files/platform"
ARCH="amd64"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
OUTPUT_DIR="${SCRIPT_DIR}/output"

CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")

FILTER=""
PIXEL_THRESHOLD=100
BASE_BRANCH="master"
CMP_BRANCH="$CURRENT_BRANCH"

while [ $# -gt 0 ]; do
    case "$1" in
        -f|--filter) FILTER="$2"; shift 2 ;;
        *) if [ "$BASE_BRANCH" = "master" ] && [ $# -ge 1 ]; then
               BASE_BRANCH="$1"; shift
               if [ $# -ge 1 ] && [ "${1:-}" != "-f" ] && [ "${1:-}" != "--filter" ]; then
                   CMP_BRANCH="$1"; shift
               fi
           else
               shift
           fi ;;
    esac
done

if [ -z "$CMP_BRANCH" ]; then
    echo "Usage: $0 [base-branch] [compare-branch] [-f filter]"
    echo "Example: $0 master ui-cleanup -f settings_access"
    echo "Defaults: base=master, compare=current git branch"
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

list_screenshots() {
    local build="$1"
    local view="$2"
    local suffix="-${view}.png"
    local filter="${FILTER}"
    curl -s "${CI_FILES}/${build}-${ARCH}/distro/${view}/screenshot/" | python3 -c "
import json, sys
for f in json.load(sys.stdin):
    name = f['name']
    if name.endswith('${suffix}'):
        if name.startswith('exception'):
            continue
        if '_unstable' in name:
            continue
        if '${filter}' and '${filter}' not in name:
            continue
        print(name)
"
}

compare_screenshots() {
    local base_dir="$1"
    local cmp_dir="$2"
    local diff_dir="$3"
    local suffix="$4"
    local changed=0
    local identical=0
    local missing=0

    for img in "${base_dir}"/*"${suffix}"; do
        [ -f "$img" ] || continue
        name=$(basename "$img")
        if [ -n "$FILTER" ] && echo "$name" | grep -qv "$FILTER"; then
            continue
        fi
        if [ ! -f "${cmp_dir}/${name}" ]; then
            echo "  MISSING: $name (not in $CMP_BRANCH)"
            missing=$((missing + 1))
            continue
        fi

        if command -v magick >/dev/null 2>&1; then
            METRIC=$(magick compare -fuzz 5% -metric AE "${base_dir}/${name}" "${cmp_dir}/${name}" "${diff_dir}/${name}" 2>&1 || true)
            PIXELS=$(echo "$METRIC" | grep -oE '^[0-9]+')
            if [ "${PIXELS:-0}" -le "$PIXEL_THRESHOLD" ]; then
                identical=$((identical + 1))
                rm -f "${diff_dir}/${name}"
            else
                echo "  CHANGED: $name (${PIXELS} pixels differ)"
                changed=$((changed + 1))
            fi
        elif command -v compare >/dev/null 2>&1; then
            METRIC=$(compare -fuzz 5% -metric AE "${base_dir}/${name}" "${cmp_dir}/${name}" "${diff_dir}/${name}" 2>&1 || true)
            PIXELS=$(echo "$METRIC" | grep -oE '^[0-9]+')
            if [ "${PIXELS:-0}" -le "$PIXEL_THRESHOLD" ]; then
                identical=$((identical + 1))
                rm -f "${diff_dir}/${name}"
            else
                echo "  CHANGED: $name (${PIXELS} pixels differ)"
                changed=$((changed + 1))
            fi
        else
            if cmp -s "${base_dir}/${name}" "${cmp_dir}/${name}"; then
                identical=$((identical + 1))
            else
                echo "  CHANGED: $name (binary diff, install ImageMagick for pixel diff)"
                cp "${cmp_dir}/${name}" "${diff_dir}/${name}"
                changed=$((changed + 1))
            fi
        fi
    done

    for img in "${cmp_dir}"/*"${suffix}"; do
        [ -f "$img" ] || continue
        name=$(basename "$img")
        if [ -n "$FILTER" ] && echo "$name" | grep -qv "$FILTER"; then
            continue
        fi
        if [ ! -f "${base_dir}/${name}" ]; then
            echo "  NEW: $name (only in $CMP_BRANCH)"
            missing=$((missing + 1))
        fi
    done

    echo "  Identical: $identical, Changed: $changed, Missing/New: $missing"
    return $changed
}

echo "Finding latest successful builds..."

BASE_BUILD=$(find_latest_build "$BASE_BRANCH")
if [ -z "$BASE_BUILD" ]; then
    echo "ERROR: No successful build found for branch '$BASE_BRANCH'"
    exit 1
fi

CMP_BUILD=$(find_latest_build "$CMP_BRANCH")
if [ -z "$CMP_BUILD" ]; then
    echo "ERROR: No successful build found for branch '$CMP_BRANCH'"
    exit 1
fi

echo "Base:    $BASE_BRANCH (build #$BASE_BUILD)"
echo "Compare: $CMP_BRANCH (build #$CMP_BUILD)"

TOTAL_CHANGED=0

for view in $VIEWS; do
    echo ""
    echo "=== ${view} ==="

    BASE_DIR="${OUTPUT_DIR}/${BASE_BRANCH}/${view}"
    CMP_DIR="${OUTPUT_DIR}/${CMP_BRANCH}/${view}"
    DIFF_DIR="${OUTPUT_DIR}/diff/${view}"

    mkdir -p "$BASE_DIR" "$CMP_DIR" "$DIFF_DIR"

    echo "Downloading base screenshots..."
    for img in $(list_screenshots "$BASE_BUILD" "$view"); do
        if [ ! -f "${BASE_DIR}/${img}" ]; then
            curl -s -o "${BASE_DIR}/${img}" "${CI_FILES}/${BASE_BUILD}-${ARCH}/distro/${view}/screenshot/${img}"
            echo "  $img"
        fi
    done

    echo "Downloading compare screenshots..."
    for img in $(list_screenshots "$CMP_BUILD" "$view"); do
        if [ ! -f "${CMP_DIR}/${img}" ]; then
            curl -s -o "${CMP_DIR}/${img}" "${CI_FILES}/${CMP_BUILD}-${ARCH}/distro/${view}/screenshot/${img}"
            echo "  $img"
        fi
    done

    echo "Comparing..."
    compare_screenshots "$BASE_DIR" "$CMP_DIR" "$DIFF_DIR" "-${view}.png" || TOTAL_CHANGED=$((TOTAL_CHANGED + $?))
done

echo ""
echo "=== Output ==="
echo "  ${OUTPUT_DIR}/"
if [ "$TOTAL_CHANGED" -gt 0 ]; then
    echo "  Diffs: ${OUTPUT_DIR}/diff/"
fi

# Copy to Android Pictures on Termux
if [ -d "$HOME/storage/pictures" ]; then
    PICTURES_DIR="$HOME/storage/pictures/screenshot-diff"
    rm -rf "$PICTURES_DIR"
    mkdir -p "$PICTURES_DIR/base" "$PICTURES_DIR/compare" "$PICTURES_DIR/diff"
    for view in $VIEWS; do
        BASE_DIR="${OUTPUT_DIR}/${BASE_BRANCH}/${view}"
        CMP_DIR="${OUTPUT_DIR}/${CMP_BRANCH}/${view}"
        for subdir_src_dest in "${BASE_DIR}:base" "${CMP_DIR}:compare" "${OUTPUT_DIR}/diff/${view}:diff"; do
            src="${subdir_src_dest%%:*}"
            dest="${subdir_src_dest##*:}"
            for img in "${src}/"*.png; do
                [ -f "$img" ] || continue
                if [ -n "$FILTER" ] && echo "$img" | grep -qv "$FILTER"; then
                    continue
                fi
                cp "$img" "$PICTURES_DIR/${dest}/"
            done
        done
    done
    echo ""
    echo "Pictures/screenshot-diff/"
    echo "  base/    — ${BASE_BRANCH} #${BASE_BUILD}"
    echo "  compare/ — ${CMP_BRANCH} #${CMP_BUILD}"
    echo "  diff/    — changed pixels highlighted in red"
fi
