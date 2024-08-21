# Abusing Makefiles for fun and profit

If you are following this blog for a while, it should be no surprise that most
of the workflow in this blog is [automated using
Go](/posts/2024-07-29/01-quick-bits-why-you-should-automate-everything.md). I
basically write Markdown files with some special rules inside the
[repository](https://github.com/thiagokokada/blog), commit and push it. In
seconds, the CI (currently [GitHub
Actions](https://github.com/thiagokokada/blog/blob/4e3f25485c6682f3e066b219df2290934bc0d256/.github/workflows/go.yml))
will take the latest commit, generate some files (since I use the [repository
itself](/posts/2024-07-26/02-using-github-as-a-bad-blog-platform.md) as a
backup blog) and publish to the [capivaras.dev
website](https://kokada.capivaras.dev/).

Now, considering how much about [Nix](https://nixos.org/) I talk in this blog,
it should be a surprise that the workflow above has **zero** Nix code inside
it. I am not saying this blog will never have it, but I am only going to add if
this is necessary, for example if I start using a tool to build this blog that
I generally don't expect it to be installed by the machine I am currently
using. Go is an exception of this rule since it is relatively straightfoward to
install (just download the [binary](https://go.dev/doc/install)) and because
its [stability guarantee](https://go.dev/doc/go1compat) means (hopefully) no
breakage. But most other things I consider moving targets, and I wouldn't be
comfortable to use unless I have Nix to ensure reproducibility.

This is why the other tool that this blog (ab)uses during its workflow is
[`Make`](https://en.wikipedia.org/wiki/Make_(software)), one of the oldest
build automation tool that exist. It is basically available in any *nix (do not
confuse with [Nix](https://nixos.org/)) system, from most Linux distros to
macOS, by default. So it is the tool I choose to automatise some tasks in this
blog, even if I consider writing a `Makefile` (the domain-specific language
that `Make` uses) kind of a lost, dark art.

To be clear, the idea of this post is not to be a `Makefile` tutorial. I will
explain some basic concepts, but if you want an actual tutorial a good one can
be found [here](https://makefiletutorial.com/). Also, while I am using `Make`
thanks to the reasons above, you can use many other tools for a similar
objective, like [Justfiles](https://github.com/casey/just),
[Taskfiles](https://taskfile.dev/) (sadly it uses
[YAML](/posts/2024-07-31/01-generating-yaml-files-with-nix.md)), or even a
small script written in any language you want. The reason that I am writing
this post is why you should do it, not how.

A quick recap on how this blog works: inside the
[repository](https://github.com/thiagokokada/blog), a post is basically a
Markdown post following the directory structure below
([permalink](https://github.com/thiagokokada/blog/tree/894a388c61ca3a38dfc9d4cbe88dc684fd964bb7)
for the current version of this blog):

```console
.
<...>
├── 2024-08-07
│   ├── 01-quick-bits-is-crostini-a-microvm.md
│   └── 02-meta-are-quick-bits-really-quick.md
├── 2024-08-11
│   └── 01-building-static-binaries-in-nix.md
├── 2024-08-12
│   ├── 01-things-i-dont-like-in-my-chromebook-duet-3.md
│   └── Screenshot_2024-08-12_20.50.42.png
├── 2024-08-13
│   ├── 01-abusing-makefiles-for-fun-and-profit.md <-- this file
├── .github
│   └── workflows
│       └── go.yml
├── .gitignore
├── go.mod
├── go.sum
├── LICENSE
├── link_rewriter.go
├── Makefile
├── mataroa.go
├── README.md
├── rss.xml
└── .scripts
    └── gen-post.sh
```

So I just create a new Markdown file following the
`YYYY-MM-DD/XX-title-slug.md` format. It **must** start with a `h1` header,
that will be automatically extract to be used as the post title, but otherwise
there is no other formatting rules. It is a highly optionated structure, but
the nice thing about being optionated is that we can extract lots of
information just from how the files are organised in the filesystem.

Most of the magic that converts those Markdown files to actual blog posts are
in the Go files that you can see above: `blog.go` is the main logic that walks
in the repository and extracts the necessary information, `mataroa.go` is
responsible for the [capivaras.dev](https://capivaras.dev/) integration (that
uses [Mataroa](https://mataroa.blog/) platform), while `link_rewriter.go` is
responsible to do some transformations in the Markdown files before posting.

While I could manage everything by just using `go` CLI and a few other *nix
commands, to make it easier to manager everything I have the following
[`Makefile`](https://github.com/thiagokokada/blog/blob/527466a2a7c8baae532281bff5db3f0695f018cb/Makefile):

```Makefile
MARKDOWN := $(wildcard ./**/*.md)

.PHONY: all
all: README.md rss.xml

blog: *.go go.*
	go build

README.md: blog $(MARKDOWN)
	./blog > README.md

rss.xml: blog $(MARKDOWN)
	./blog -rss > rss.xml

.PHONY: publish
publish: blog
	./blog -publish

DAY := $(shell date)
_PARSED_DAY := $(shell date '+%Y-%m-%d' -d '$(DAY)')
.PHONY: day
day:
	mkdir -p '$(_PARSED_DAY)'

TITLE = $(error TITLE is not defined)
.PHONY: post
post: blog day
	./.scripts/gen-post.sh '$(_PARSED_DAY)' '$(TITLE)'

FILE = $(error FILE is not defined)
.PHONY: draft
draft:
	cd '$(dir $(FILE))' && mv '$(notdir $(FILE))' '.$(notdir $(FILE))'

.PHONY: undraft
undraft:
	cd '$(dir $(FILE))' && mv '$(notdir $(FILE))' '$(patsubst .%,%,$(notdir $(FILE)))'

.PHONY: clean
clean:
	rm -rf blog
```

For those unfamiliar with `Makefile`, a quick explanation on how it works from
[Wikipedia](https://en.wikipedia.org/wiki/Make_(software)#Makefile):

> Each rule begins with a _dependency line_ which consists of the rule's target
> name followed by a colon (:) and optionally a list of targets on which the
> rule's target depends, its prerequisites.

So if we look for example at the `blog` binary, the dependencies are all the
`.go` files and Go module files like `go.mod` and `go.sum`. We can make the
`blog` binary by running:

```console
$ make blog
go build
```

One nice thing about `Makefile` is that they track if any of the source files
has a newer timestamp than the target file, and only trigger the build again if
there are changes, for example:

```console
$ make blog
make: 'blog' is up to date.

$ touch blog.go

$ make blog
go build
```

But sometimes this property is undesirable. In those cases we can declare a
target as `.PHONY`, that basically instructs `Makefile` to always make the
target. One classic example is `clean` target, that removes build artifacts:

```console
$ make clean
rm -rf blog

$ make clean
rm -rf blog
```

By the way, it is better to declare a target as `.PHONY` than declaring
dependencies incorrectly, especially in languages that has fast build times
like e.g.: Go. The worst thing that can happen is something not being rebuild
when it needs to. So my recomendation if you are writing your first `Makefile`
is to just declare everything as `.PHONY`. You can always improve it later.

One last basic concept that I want to explain about `Makefile` is the default
target: it is the target that is run if you just run `make` without arguments
inside the directory that contains a `Makefile`. The default target is
generally the first target in the `Makefile`. It is common to have an `all`
target (that is also marked as `.PHONY`) that has as dependencies all the
targets that you want to build by default. In this particular case I declare
the `README.md` and `rss.xml` files to be build by default, and they themselves
depends in `blog` binary being build. So once I run `make` you get as result:

```console
$ make
go build
./blog > README.md
./blog -rss > rss.xml
```

And this result above highlights the first reason I think you should have a
`Makefile` or something similar in your projects: you don't need to remember
the exactly steps that you need to get things working. If I see one project of
mine having a `Makefile`, I can be reasonably confident that I can get it
working by just running `make`.

But now let's focus in the other targets that I have in the `Makefile` that are
not related to the build process but are there to help me manage my blog posts.
Remember the rules I explained above? Maybe not, but it should be no problem,
because:

```
$ make post TITLE="My new blog post"
mkdir -p "2024-08-13"
./.scripts/gen-post.sh "2024-08-13" "My new blog post"
Creating file: 2024-08-13/02-my-new-blog-post.md

$ cat 2024-08-13/02-my-new-blog-post.md
# My new blog post
```

This command, `make post`, is responsible for:

1. Create a new directory for today, if it doesn't exist
2. Run the
   [`gen-post.sh`](https://github.com/thiagokokada/blog/blob/6a3b06970729f7650e5bee5fb0e1f9f2541ffea8/.scripts/gen-post.sh)
script, that:
   1. Enumerates all posts from the day, so we can number the new post correctly
      - We already had this post planned for 2024-08-13, so the new post is 02
   2. Slugify the title, so we can create each Markdown file with the correct
   filename
   3. Creates a new Markdown file with the title as a `h1` header

The steps above may or may not seen trivial, and for a while I was doing them
manually. But not having to think what is the current date or if I already
posted that day or what is the slug is for the title make (pun intended) my
like much easier.

Yes, the code is ugly. The way variables works in `Make` is that you can
declare then inside the `Makefile`, but they can be overwritten in the terminal
if you pass them. I used this to allow `make post` to also work for future
posts:

```console
$ make post TITLE="Another new blog post" DAY=2024-12-12
mkdir -p "2024-12-12"
./.scripts/gen-post.sh "2024-12-12" "Another new blog post"
Creating file: 2024-12-12/01-another-new-blog-post.md
```

So in the above case, `DAY` is filled with the value passed in the terminal
instead of default (that would be the current day), and `_PARSED_DAY` is the
day we use to actually create the directory. We can actually pass any date
format recognised by
[`date`](https://www.gnu.org/software/coreutils/manual/html_node/Examples-of-date.html),
not just `YYYY-MM-DD`.

I have 2 other phony targets that I want to talk, `draft` and `undraft`. They
expect a `FILE` to be passed, and I use them to either hide or unhide a file:

```console
$ make draft FILE=2024-12-12/01-another-new-blog-post.md
mv "2024-12-12/01-another-new-blog-post.md" "2024-12-12/.01-another-new-blog-post.md"

$ make undraft FILE=2024-12-12/.01-another-new-blog-post.md
mv "2024-12-12/.01-another-new-blog-post.md" "2024-12-12/01-another-new-blog-post.md"
```

Why? Because hidden files are [explicit
ignored](https://github.com/thiagokokada/blog/blob/894a388c61ca3a38dfc9d4cbe88dc684fd964bb7/blog.go#L101-L104)
during my directory parser to mean they're a draft post and not ready to be
published. And the reason I created those targets is because I was tired of
trying to hide or unhide a file manually.

So that's it, for the same reason you [should probably automate
everything](/posts/2024-07-29/01-quick-bits-why-you-should-automate-everything.md),
you also need to have some way to automate your tasks. `Makefile` is one way to
do it, maybe not the best way to do it, but it works and it is available
anywhere.
