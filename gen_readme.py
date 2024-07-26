#!/usr/bin/env python3

"""\
Usage: gen_readme.py

Needs to be run at the root directory of the repository.

Will output the generated README.md file to the stdout, you can redirect its
contents by:

  $ ./gen_readme.py > README.md

Will also generate the RSS file in `rss.xml`.\
"""

import os
import re
import sys
import xml.etree.cElementTree as ET
from collections import defaultdict
from datetime import datetime
from pathlib import Path
from urllib.parse import urljoin

README_TEMPLATE = """\
# Blog

Backup of my blog posts in https://kokada.capivaras.dev/.

## Posts

{posts}\
"""
RSS_TITLE = "kokada's blog"
RSS_LINK = "https://github.com/thiagokokada/blog"
RSS_DESCRIPTION = "dd if=/dev/urandom of=/dev/brain0"
RSS_POST_LINK_PREFIX = f"{RSS_LINK}/blob/main/"

Posts = dict[datetime, list[dict[str, str]]]


def grab_titles(pwd: Path) -> Posts:
    posts = defaultdict(list)

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
            # Ignore non-markdown files or hidden files (draft)
            if not post.suffix == ".md" or post.name[0] == ".":
                continue

            # Grab the first H1 section to parse as title
            text = post.read_text()
            mTitle = re.match(r"# (?P<title>.*)\r?\n", text)
            if mTitle and (title := mTitle.groupdict().get("title")):
                posts[date].append({"title": title, "file": post})
            else:
                print(f"WARN: did not find title for file: {post}", file=sys.stderr)

    return posts


def gen_readme(posts: Posts):
    titles = []

    for date, dayPosts in posts.items():
        for post in dayPosts:
            link = os.path.join(".", post["file"])  # to format as ./{filepath}
            title = date.strftime(f"- [{post['title']}]({link}) - %Y-%m-%d")
            titles.append(title)

    print(README_TEMPLATE.format(posts="\n".join(titles)))


def gen_rss(posts: Posts):
    rss = ET.Element("rss", version="2.0")

    channel = ET.SubElement(rss, "channel")
    ET.SubElement(channel, "title").text = RSS_TITLE
    ET.SubElement(channel, "link").text = RSS_LINK
    ET.SubElement(channel, "description").text = RSS_DESCRIPTION

    item = ET.SubElement(channel, "item")
    for date, dayPost in posts.items():
        for post in dayPost:
            ET.SubElement(item, "title").text = post["title"]
            ET.SubElement(item, "link").text = urljoin(
                RSS_POST_LINK_PREFIX, str(post["file"])
            )

    tree = ET.ElementTree(rss)
    ET.indent(tree, space="\t", level=0)
    tree.write("rss.xml", xml_declaration=True, encoding="UTF-8")


def main():
    if "-h" in sys.argv:
        print(__doc__, file=sys.stderr)
        sys.exit(0)
    titles = grab_titles(Path())
    gen_readme(titles)
    gen_rss(titles)


if __name__ == "__main__":
    main()
