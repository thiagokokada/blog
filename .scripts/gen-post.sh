#!/bin/sh

set -eu

# Expect to be set as environment variables
readonly DATE POST_ROOT TITLE SLUG

i=1
while [ "$i" -ne 100 ]; do
	n="$(printf '%.2d' "$i")"
	# Match either normal files or hidden post files
	if ls "$POST_ROOT/$DATE/$n"*".md" >/dev/null 2>&1 || \
		ls "$POST_ROOT/$DATE/.$n"*".md" >/dev/null 2>&1; then
		i=$(( i + 1 ))
		continue
	fi

	readonly file="$POST_ROOT/$DATE/$n-$SLUG.md"
	echo "# ${TITLE}" > "$file"
	echo "$file"
	exit 0
done

exit 1
