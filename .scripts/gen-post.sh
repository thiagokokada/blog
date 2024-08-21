#!/bin/sh

set -eu

# Expect to be set as environment variables
readonly DATE SLUG TITLE

i=1
while [ "$i" -ne 100 ]; do
	n="$(printf '%.2d' "$i")"
	# Match either normal files or hidden post files
	if ls "$DATE/$n"*".md" >/dev/null 2>&1 || \
		ls "$DATE/.$n"*".md" >/dev/null 2>&1; then
		i=$(( i + 1 ))
		continue
	fi

	file="$DATE/$n-$SLUG.md"
	echo "# ${TITLE}" > "$file"
	echo "$file"
	exit 0
done

exit 1
