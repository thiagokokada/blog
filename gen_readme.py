#!/usr/bin/env python3

"""\
Usage: gen_readme.py

Needs to be run at the root directory of the repository.

Will output the generated README.md file to the stdout, you can redirect its
contents by:

  $ ./gen_readme.py > README.md\
"""

import re
import os
import sys
from datetime import datetime
from pathlib import Path

README_TEMPLATE = """\
# Blog

Backup of my blog posts in https://kokada.capivaras.dev/.

## Posts

{posts}\
"""


def grab_titles(pwd: Path) -> list[str]:
    titles = []

    for dir in sorted(pwd.iterdir()):
        # Ignore non-directories or hidden files
        if not dir.is_dir() or dir.name[0] == ".":
            continue

        # Try to parse date from directory name
        try:
            date = datetime.strptime(dir.name, "%Y-%m-%d")
        except ValueError:
            print(f"WARN: ignoring non-date directory: {dir}", file=sys.stderr)
            continue

        # Iterate between the files in the date directory
        for post in sorted(dir.iterdir()):
            # Ignore non-markdown files
            if not post.suffix == ".md":
                continue

            # Grab the first H1 section to parse as title
            m = re.match(r"# (?P<title>.*)\r?\n", post.read_text())
            if m and (title := m.groupdict().get("title")):
                link = os.path.join(".", post)  # to format with ./{filepath}
                titles.append(date.strftime(f"- [{title}]({link}) - %Y-%m-%d"))
            else:
                print(f"WARN: did not find title for file: {post}", file=sys.stderr)

    return titles


def main():
    if "-h" in sys.argv:
        print(__doc__, file=sys.stderr)
        sys.exit(0)

    titles = grab_titles(Path())
    print(README_TEMPLATE.format(posts="\n".join(titles)))


if __name__ == "__main__":
    main()
