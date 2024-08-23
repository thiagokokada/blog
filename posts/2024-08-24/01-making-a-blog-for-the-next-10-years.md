# Making a blog for the next 10 years

So one thing that I realise after starting writing this blog is that I care
more about it than some of my other projects. For some reason or another I want
to make sure that this blog will continue with me for a long time. This is one
of the reasons why I use [GitHub as mirror
blog](/posts/2024-07-26/02-using-github-as-a-bad-blog-platform.md) and why I
created a [bunch of
automation](/posts/2024-07-29/01-quick-bits-why-you-should-automate-everything.md)
to make sure I never forget how to maintain this blog.

Still, there are a bunch of dependencies that I need to make sure they're
working so I can publish blog posts:

- Markdown
- A text editor
- Hosting and blog platform
- Shell script and
[Makefile](/posts/2024-08-13/01-abusing-makefiles-for-fun-and-profit.md)
- Go compiler and dependencies

Let's start with the most important one: the texts themselves, they're all
written in [Markdown](https://www.markdownguide.org/). The nice thing about
Markdown is that it is just plain text files with some special notation that
allow you to format text, but the text itself is still legible even if you have
no support to render Markdowns. So it means as long as I can read a plain text
file, I can still read the posts, no issue there. Closely related issue here is
the text editor, but the nice thing about Markdown is that if tomorrow I decide
to change from `neovim` to, say, [Zed](https://zed.dev/), I can still write it
Markdown files without issue. I also use a standardised Markdown implementation
called [CommonMark](https://commonmark.org/), that fixes a bunch of ambiguities
of the original Markdown spec.

The next point is the host ([capivaras.dev](capivaras.dev)) and the blog
platform ([Mataroa](https://github.com/mataroa-blog/mataroa)). One of the nice
things about Mataroa is that it avoids platform lock-in by having multiple ways
to [export your posts](https://mataroa.blog/export/). I could export everything
to [Hugo](https://gohugo.io/), for example, and serve the blog from [GitHub
Pages](https://pages.github.com/).

This is especially nice consider how small [capivaras.dev](capivaras.dev) is,
so it is good to know that if I eventually have issues I could move to
somewhere else. I also have an e-mail backup every month with all posts made
(another [Mataroa
feature](https://hey.mataroa.blog/blog/monthly-auto-exports-via-email/)), and
of course I have a [Git repo](https://github.com/thiagokokada/blog) that also
acts as a [mirror of this
blog](/posts/2024-07-26/02-using-github-as-a-bad-blog-platform.md). So I would
say the chance of losing access to the content is slim.

One other issue is the URL of the posts that are indexed in many different
places, like [Hacker News](https://news.ycombinator.com/),
[Lobte.rs](https://lobste.rs/), etc. This is why I finally decide to bite the
bullet and purchase a proper domain, and this website should now be available
in [kokada.dev](kokada.dev). This means that in my eventual exit from
[capivaras.dev](capivaras.dev), I can just point my new blog location to my own
domain if needed (it is not as easy since I also need to preserve the post
URLs, but shouldn't be difficult to fix this if I ever need to do so).

Now for the tools that I use to publish from the original Markdown files to
everything else. Let's start with shell script(s) and Makefile: I decided that
they're less of an issue if they eventually stop working: they're only used to
make my life easier, but I can still publish files manually if needed. Still, I
tried to rewrite both the
[shell](https://github.com/thiagokokada/blog/commit/a0d421ca90f3da059998295c5e3c6c7a6a3f0688)
and
[Makefile](https://github.com/thiagokokada/blog/commit/074580065b21fbdaf930aa51968e69f015d49505)
to avoid GNUisms, so in the eventual case that I decide to e.g.: stop using a
GNU/Linux system like NixOS and use a *BSD system instead, I am covered.

Go is the more important part: the tooling used to [publish this blog is
written in
Go](/posts/2024-07-29/01-quick-bits-why-you-should-automate-everything.md). Go
is a good language when you want to ensure that things will work for a long
time because of its [backwards compatibility
guarantee](https://go.dev/blog/compat). Also I don't expect Google dropping Go
development soon, but even if this happen (["killed by
Google"](https://killedbygoogle.com/) is a thing after all), it is very likely
some other group or company would adopt its development quickly, considering
[how popular the language](https://www.tiobe.com/tiobe-index/go/) is.

However, the [Go
modules](https://github.com/thiagokokada/blog/blob/main/go.mod) that I depend
are another story:

- [elliotchance/orderedmap](https://github.com/elliotchance/orderedmap/): an
ordered map implementation that I use until Go adds it in the [standard
library](/posts/2024-08-17/01-an-unordered-list-of-things-i-miss-in-go.md)
- [gorilla/feeds](https://github.com/gorilla/feeds): a RSS generator library
- [gosimple/slug](https://github.com/gosimple/slug): a
[slug](https://developer.mozilla.org/en-US/docs/Glossary/Slug) generator
library
- [yuin/goldmark](https://github.com/yuin/goldmark): a CommonMark parser and
renderer
- [teekennedy/goldmark-markdown](https://github.com/teekennedy/goldmark-markdown):
a renderer for Goldmark to render back to Markdown (since Goldmark itself
doesn't have this capacity)

In common for all those modules are that they're clearly small projects
maintained mostly by one developer. They're all very good, don't get me wrong,
but they're still an reliability issue in the future. There is no guarantee
those repositories will not be deleted tomorrow, for example.

Yes, [Go Proxy](https://proxy.golang.org/) exist, but from what I understood
reading its page is that while it caches modules contents, this is not
guarantee:

> proxy.golang.org does not save all modules forever. There are a number of
> reasons for this, but one reason is if proxy.golang.org is not able to detect
> a suitable license. In this case, only a temporarily cached copy of the
> module will be made available, and may become unavailable if it is removed
> from the original source and becomes outdated. The checksums will still
> remain in the checksum database regardless of whether or not they have become
> unavailable in the mirror.

This is why this is the first project that made sense to me to use [`go mod
vendor`](https://go.dev/ref/mod#go-mod-vendor). Now I have a copy of the source
code of all modules inside the
[vendor](https://github.com/thiagokokada/blog/tree/0b97630d6b30551ffe05b5d8124305b1065f729d/vendor)
directory in the repository, avoiding the risks I commented above. This allows
me to ensure that this blog will still be publishable in the future, as long as
I have a working Go toolchain (and Go toolchain makes this
[easy](/posts/2024-08-19/01-quick-bits-go-automatically-downloads-a-newer-toolchain-if-needed.md)).

There is a few other things that can bitrot this blog, for example links going
nowhere. I always try to use
[permalinks](https://en.wikipedia.org/wiki/Permalink) where it makes sense, but
the only actual way to ensure those links would work in the future would be to
point them to [archive.org](https://archive.org/) (but even archive.org may not
exist forever). Maybe something to fix in the future, hope not in the far
future, before things start to break.
