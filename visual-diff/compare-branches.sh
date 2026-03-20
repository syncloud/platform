#!/bin/bash
set -eu

# Download screenshots from two branches and compare them.
# Usage: compare-branches.sh [base-branch] [compare-branch] [-f filter]
# Defaults: base=master, compare=current git branch

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
OUTPUT_DIR="${SCRIPT_DIR}/output"

CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")

FILTER=""
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
    echo "Defaults: base=master, compare=current git branch"
    exit 1
fi

# Download both branches
"${SCRIPT_DIR}/download-artifacts.sh" "$BASE_BRANCH"
"${SCRIPT_DIR}/download-artifacts.sh" "$CMP_BRANCH"

BASE_DIR="${OUTPUT_DIR}/${BASE_BRANCH}"
CMP_DIR="${OUTPUT_DIR}/${CMP_BRANCH}"
DIFF_DIR="${OUTPUT_DIR}/diff"

rm -rf "$DIFF_DIR"

# Run the diff
FILTER_ARG=""
if [ -n "$FILTER" ]; then
    FILTER_ARG="-f $FILTER"
fi
"${SCRIPT_DIR}/screenshot-diff.sh" "$BASE_DIR" "$CMP_DIR" "$DIFF_DIR" $FILTER_ARG || true

# Copy to Android Pictures on Termux
if [ -d "$HOME/storage/pictures" ]; then
    PICTURES_DIR="$HOME/storage/pictures/screenshot-diff"
    rm -rf "$PICTURES_DIR"
    mkdir -p "$PICTURES_DIR/base" "$PICTURES_DIR/compare" "$PICTURES_DIR/diff"
    for view in desktop mobile; do
        for subdir_src_dest in "${BASE_DIR}/${view}:base" "${CMP_DIR}/${view}:compare" "${DIFF_DIR}/${view}:diff"; do
            src="${subdir_src_dest%%:*}"
            dest="${subdir_src_dest##*:}"
            [ -d "$src" ] || continue
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
    echo "  base/    — ${BASE_BRANCH}"
    echo "  compare/ — ${CMP_BRANCH}"
    echo "  diff/    — changed pixels highlighted in red"
fi
