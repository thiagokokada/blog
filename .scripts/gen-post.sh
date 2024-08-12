#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(dirname -- "${BASH_SOURCE[0]}")"
readonly SCRIPT_DIR

DAY="$1"
TITLE="$2"
SLUG="$(cd "$SCRIPT_DIR" && ../blog -slugify "$TITLE")"
readonly DAY TITLE SLUG

for i in $(seq -f "%02g" 99); do
	# Match either normal files or hidden post files
	if compgen -G "$DAY/$i*.md" &>/dev/null || \
		compgen -G "$DAY/.$i*.md" &>/dev/null; then
		continue
	fi

	file="$DAY/$i-$SLUG.md"
	echo "Creating file: $file"
	echo "# ${TITLE}" > "$file"
	exit 0
done

exit 1
