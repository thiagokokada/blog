#!/usr/bin/env bash

set -euo pipefail

# Expect to be set as environment variables
readonly DATE TITLE

SCRIPT_DIR="$(dirname -- "${BASH_SOURCE[0]}")"
readonly SCRIPT_DIR

SLUG="$(cd "$SCRIPT_DIR" && ../blog -slugify "$TITLE")"
readonly SLUG

for i in $(seq -f "%02g" 99); do
	# Match either normal files or hidden post files
	if compgen -G "$DATE/$i*.md" &>/dev/null || \
		compgen -G "$DATE/.$i*.md" &>/dev/null; then
		continue
	fi

	file="$DATE/$i-$SLUG.md"
	echo "Creating file: $file"
	echo "# ${TITLE}" > "$file"
	exit 0
done

exit 1
