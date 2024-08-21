# Quick bits: why you should automate everything

If everything works as expected this blog post should appear in [in my
blog](https://kokada.capivaras.dev/) without I ever touching the
[capivaras.dev](https://capivavas.dev) website. I rewrote my [previous Python
script](/posts/2024-07-26/02-using-github-as-a-bad-blog-platform.md) to Go
([permalink](https://github.com/thiagokokada/blog/blob/3c39e0f7cd58b1af885f69871490b05bf6fc7d99/blog.go))
since my attempt to generate proper description to the RSS feed resulted in
slow startup times (not because of Python, but because of my usage of
`nix-shell` since I didn't want to deal with
[venv](https://docs.python.org/3/library/venv.html) or anything to manage my
Python dependencies).

My previous workflow of this blog already involved me writing the texts in
[neovim](https://neovim.io/), copying and pasting the result in the
[capivaras.dev](https://capivavas.dev) website and publishing. This was not
that bad, except that it seems I have a heavy tendency of editing my posts
multiple times. Copying and pasting data between neovim and the website became
tedious, so I decided to give up and automate the whole process.

[Mataroa](https://mataroa.blog/) (the blog platform
[capivaras.dev](https://capivavas.dev) run) has a reasonable good
[API](https://mataroa.blog/api/docs/), and it only took a few hours to get a
version of publishing working (it would take less if
[Django](https://www.djangoproject.com/), the framework Mataroa is written, did
not have a weird behavior with URLs missing a trailing `/`). An additional few
lines of
[YAML](https://github.com/thiagokokada/blog/blob/51b20612335c7f4312a51a0f436235b4b701ce8b/.github/workflows/go.yml)
to make GitHub Actions trigger a pipeline and now I should never have to
manually update my blog again.

I could have not done this. I mean, I probably wasted more time writing an
automation than I actually wasted publishing manually. But the manual process
is easy to miss, and I already did multiple mistakes publishing in the manual
method. For example, when writing the Markdown files, each post is written in a
particular format, where the first header is considered the title, so I need to
remove it from the contents during publication. But of course, this is easy to
miss, and I had to fix this multiple times already.

So yes, I think this is a good lesson on why you should automate everything. It
is more than just about [time savings](https://xkcd.com/1205/), it is about
reducing mistakes and even documenting (even if poorly) a process. I mean, the
code I wrote is not that great, but I can definitely rely on it in the future
to remember what I need to do. It will be much faster than trying to learn from
scratch again.
