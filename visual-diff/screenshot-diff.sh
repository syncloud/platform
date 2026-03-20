#!/bin/bash
set -eu

# Compare screenshots in two local directories.
# Usage: screenshot-diff.sh <base-dir> <compare-dir> [diff-dir] [-f filter]
#
# Exits with the number of changed screenshots (0 = all identical).
# Skips files named 'exception*' or containing '_unstable'.

PIXEL_THRESHOLD=100
FILTER=""

BASE_DIR=""
CMP_DIR=""
DIFF_DIR=""

while [ $# -gt 0 ]; do
    case "$1" in
        -f|--filter) FILTER="$2"; shift 2 ;;
        *)
            if [ -z "$BASE_DIR" ]; then
                BASE_DIR="$1"
            elif [ -z "$CMP_DIR" ]; then
                CMP_DIR="$1"
            elif [ -z "$DIFF_DIR" ]; then
                DIFF_DIR="$1"
            fi
            shift ;;
    esac
done

if [ -z "$BASE_DIR" ] || [ -z "$CMP_DIR" ]; then
    echo "Usage: $0 <base-dir> <compare-dir> [diff-dir] [-f filter]"
    exit 1
fi

if [ -z "$DIFF_DIR" ]; then
    DIFF_DIR="${CMP_DIR}/../diff"
fi

VIEWS="desktop mobile"
TOTAL_CHANGED=0

for view in $VIEWS; do
    BASE_VIEW="${BASE_DIR}/${view}"
    CMP_VIEW="${CMP_DIR}/${view}"
    DIFF_VIEW="${DIFF_DIR}/${view}"

    if [ ! -d "$BASE_VIEW" ] || [ ! -d "$CMP_VIEW" ]; then
        continue
    fi

    mkdir -p "$DIFF_VIEW"

    echo ""
    echo "=== ${view} ==="

    changed=0
    identical=0
    missing=0

    for img in "${BASE_VIEW}/"*.png; do
        [ -f "$img" ] || continue
        name=$(basename "$img")

        # Skip exception and unstable screenshots
        case "$name" in exception*|*_unstable*) continue ;; esac

        if [ -n "$FILTER" ] && echo "$name" | grep -qv "$FILTER"; then
            continue
        fi

        if [ ! -f "${CMP_VIEW}/${name}" ]; then
            echo "  MISSING: $name"
            missing=$((missing + 1))
            continue
        fi

        if command -v magick >/dev/null 2>&1; then
            METRIC=$(magick compare -fuzz 5% -metric AE "${BASE_VIEW}/${name}" "${CMP_VIEW}/${name}" "${DIFF_VIEW}/${name}" 2>&1 || true)
            PIXELS=$(echo "$METRIC" | grep -oE '^[0-9]+')
            if [ "${PIXELS:-0}" -le "$PIXEL_THRESHOLD" ]; then
                identical=$((identical + 1))
                rm -f "${DIFF_VIEW}/${name}"
            else
                echo "  CHANGED: $name (${PIXELS} pixels differ)"
                changed=$((changed + 1))
            fi
        elif command -v compare >/dev/null 2>&1; then
            METRIC=$(compare -fuzz 5% -metric AE "${BASE_VIEW}/${name}" "${CMP_VIEW}/${name}" "${DIFF_VIEW}/${name}" 2>&1 || true)
            PIXELS=$(echo "$METRIC" | grep -oE '^[0-9]+')
            if [ "${PIXELS:-0}" -le "$PIXEL_THRESHOLD" ]; then
                identical=$((identical + 1))
                rm -f "${DIFF_VIEW}/${name}"
            else
                echo "  CHANGED: $name (${PIXELS} pixels differ)"
                changed=$((changed + 1))
            fi
        else
            if cmp -s "${BASE_VIEW}/${name}" "${CMP_VIEW}/${name}"; then
                identical=$((identical + 1))
            else
                echo "  CHANGED: $name (binary diff, install ImageMagick for pixel diff)"
                cp "${CMP_VIEW}/${name}" "${DIFF_VIEW}/${name}"
                changed=$((changed + 1))
            fi
        fi
    done

    for img in "${CMP_VIEW}/"*.png; do
        [ -f "$img" ] || continue
        name=$(basename "$img")
        case "$name" in exception*|*_unstable*) continue ;; esac
        if [ -n "$FILTER" ] && echo "$name" | grep -qv "$FILTER"; then
            continue
        fi
        if [ ! -f "${BASE_VIEW}/${name}" ]; then
            echo "  NEW: $name"
            missing=$((missing + 1))
        fi
    done

    echo "  Identical: $identical, Changed: $changed, Missing/New: $missing"
    TOTAL_CHANGED=$((TOTAL_CHANGED + changed))
done

echo ""
if [ "$TOTAL_CHANGED" -gt 0 ]; then
    echo "FAIL: ${TOTAL_CHANGED} screenshots differ"
    exit 1
else
    echo "PASS: All screenshots match"
fi
