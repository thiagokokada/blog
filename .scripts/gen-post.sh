#!/usr/bin/env bash

set -euo pipefail

day="$1"
title="$2"
slug="$(./blog -slugify "$title")"
readonly day title slug

for i in $(seq -f "%02g" 99); do
	# shellcheck disable=SC2144
	if [ -f "$day/$i"*".md" ]; then
		continue
	fi

	file="$day/$i-$slug.md"
	echo "Creating file: $file"
	echo "# ${title}" > "$file"
	exit 0
done

exit 1
