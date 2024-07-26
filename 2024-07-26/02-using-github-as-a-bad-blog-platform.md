# Using GitHub as a (bad) blog platform

I finally started a new blog, thanks to the offer of
[@ratsclub](https://gluer.org/) to give me free access to
[capivaras.dev](https://capivaras.dev/). But considering how small this blog
platform is supposed to be, I want to have at least somewhere to have a backup
of my posts. I know [Mataroa](https://mataroa.blog/), the blog platform that
[capivaras.dev](https://capivaras.dev/) runs, has automatic e-mail backups, but
I want something more reliable.

I am writing all my posts in Markdown (the format that [Mataroa
supports](https://1.mataroa.blog/guides/markdown/)) files inside
[neovim](https://neovim.io/) anyway, so why not store all my Markdown files in
Git? So this is what I did, I now have an unofficial mirror in
[GitHub](https://github.com/thiagokokada/blog).

While I am here, why not to overcomplicate? Can I make an usable blog platform
from GitHub? I mean, it already renders Markdown files by default, so no need
to do anything in that space. To reach feature parity with
[capivaras.dev](https://capivaras.dev/), I only need to have an index and RSS
(since comments are not supported anyway). No need for newsletter since GitHub
has a [watch
feature](https://docs.github.com/en/account-and-profile/managing-subscriptions-and-notifications-on-github/managing-subscriptions-for-activity-on-github/viewing-your-subscriptions)
already.

After a couple of hours hacking a Python script, you can see the result of this
monstrosity [here](https://github.com/thiagokokada/blog). The script, called
`gen_blog.py`, is available at the same repository (here is a
[permalink](https://github.com/thiagokokada/blog/blob/c8986d1ab1b94c0986fd814629bb8eb4034fb6e7/gen_blog.py)).
It automatically generates an index at
[`README.md`](https://github.com/thiagokokada/blog/blob/main/README.md) with
each blog post and a
[`rss.xml`](https://raw.githubusercontent.com/thiagokokada/blog/main/rss.xml)
file at the root of the repository.

Instead of trying to explain the code, I am going to explain the general idea,
because I think that if you want to replicate this idea it is better to rewrite
it in a way that you understand. It shouldn't take more than 2 hours in any
decent programming language. But if you really want, the script itself is
licensed in [WTFPL](https://en.wikipedia.org/wiki/WTFPL) license. The code only
uses Python 3's standard library and should work in any relatively recent
version (anything newer than 3.9 should work).

So the idea is basically to organise the repository and the Markdown files in a
easy way that makes it trivial to parse in a deterministic way. For example, my
repository is organised in the following way:

```
root
├── 2024-07-26
│   ├── 01-writing-nixos-tests-for-fun-and-profit.md
│   └── 02-using-github-as-a-bad-blog-platform.md <- this file
├── gen_blog.py
├── README.md
└── rss.xml
```

Each day that you write a new blog post will be on its own directory. This is
nice because Markdown files may include extra files than the posts themselves,
e.g.: images, and this organisation make it trivial to organise everything.

Each post has its own Markdown file. I put a two digit number before each post,
to ensure that when publishing multiple posts at the same day I keep them in
the same order of publishing. But if you don't care about it, you can just name
the files whatever you want.

Also, I am assuming that each Markdown file has a header starting with `# `,
and that is the title of the blog post.

Using the above organisation, I have this function that scrap the repository
and collect the necessary information to generate the index and RSS files:

```python
def grab_posts(pwd: Path) -> Posts:
    posts = defaultdict(list)

    for dir in sorted(pwd.iterdir(), reverse=True):
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
        for post in sorted(dir.iterdir(), reverse=True):
            # Ignore non-markdown files or hidden files (draft)
            if not post.suffix == ".md" or post.name[0] == ".":
                continue

            # Grab the first H1 section to parse as title
            text = post.read_text()
            mTitle = re.match(r"# (?P<title>.*)\r?\n", text)
            if mTitle and (title := mTitle.groupdict().get("title")):
                posts[date].append({"title": title, "file": str(post)})
            else:
                print(f"WARN: did not find title for file: {post}", file=sys.stderr)

    return posts
```

Some interesting tidbits: if a Markdown file has a `.` at the start I assume it
is a draft post, and ignore it from my scrapper. I added a bunch of `WARN`
prints to make sure that the me in the future doesn't do anything dumb. Also,
sorting in reverse since reverse chronological order is the one most people
expect in blogs (i.e.: more recent blog posts at top).

After running the function above, I have a resulting dictionary that I can use
to generate either a `README.md` file or Markdown:

```python
def gen_readme(posts: Posts):
    titles = []

    for date, dayPosts in posts.items():
        for post in dayPosts:
            # This creates a relative link to the Markdown file, .e.g.:
            # ./02-using-github-as-a-bad-blog-platform.md
            link = os.path.join(".", post["file"])
            # This formats the title, e.g.:
            # - [Using GitHub as a (bad) blog platform](./2024-07-26/02-using-github-as-a-bad-blog-platform.md) - 2024-07-26
            title = date.strftime(f"- [{post['title']}]({link}) - %Y-%m-%d")
            # This appends to the list to generate the content later
            titles.append(title)

    # README_TEMPLATE is a string with the static part of the README
    print(README_TEMPLATE.format(posts="\n".join(titles)))


def gen_rss(posts: Posts):
    # Got most of the specification from here:
    # https://www.w3schools.com/XML/xml_rss.asp
    rss = ET.Element("rss", version="2.0")

    # Here are the RSS metadata for the blog itself
    channel = ET.SubElement(rss, "channel")
    ET.SubElement(channel, "title").text = "kokada's blog"
    ET.SubElement(channel, "link").text = "https://github.com/thiagokokada/blog"
    ET.SubElement(channel, "description").text = "dd if=/dev/urandom of=/dev/brain0"

    # You create one item for each blog post
    for date, dayPost in posts.items():
        for post in dayPost:
            item = ET.SubElement(channel, "item")
            link = urljoin(RSS_POST_LINK_PREFIX, post["file"])
            ET.SubElement(item, "title").text = post["title"]
            ET.SubElement(item, "guid").text = link
            ET.SubElement(item, "link").text = link
            ET.SubElement(item, "pubDate").text = date.strftime('%a, %d %b %Y %H:%M:%S GMT')

    # Generate the XML and indent
    tree = ET.ElementTree(rss)
    ET.indent(tree, space="\t", level=0)
    tree.write("rss.xml", xml_declaration=True, encoding="UTF-8")
```

To publish a new Post, a basically write a Markdown file, run
`./gen_readme.py > README.md` at the root of the repository, and see the magic
happens.

It works much better than I initially antecipated. The `README.md` is properly
populated with the titles and links. The RSS is kind empty since it has no
description, but it seems to work fine (at least in
[Inoreader](https://www.inoreader.com/), my RSS reader of choice). I can
probably fill the post description with more information if I really want, but
it is enough for now. Not sure who is that interested in my writing that will
want to use this RSS feed instead the one available in
[capivaras.dev](https://kokada.capivaras.dev/rss/) anyway.

So that is it. I am not saying this is a good idea for your primary blog
platform or whatever, and I still prefer to publish to a platform that doesn't
track users or have tons of JavaScript or whatever. But if you want a backup of
your posts and you are already writing Markdown anyway, well, there are worse
ways to do it I think.
