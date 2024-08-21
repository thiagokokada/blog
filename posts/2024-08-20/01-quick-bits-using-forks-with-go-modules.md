# Quick bits: using forks with Go modules

There are 2 types of blog posts: the ones you write for others, and the ones
that you write for yourself. This blog post is the latter kind. What I am going
to talk here is probably something most people know, but I didn't, and the
documentation is all scattered in the internet. So I am writing mostly to
remember myself, in future, if I have the same kind of issue what I need to do.

The context: [Mataroa](https://mataroa.blog/), the blog platform
[capivaras.dev](https://capivaras.dev) is hosted on, relies in
[python-markdown](https://python-markdown.github.io/), predates
[CommonMark](https://commonmark.org/) specification and as such, has some
incompatibilities. One of those incompatibilities with CommonMark is the way
sub lists are handled. From the documentation:

> The syntax rules clearly state that when a list item consists of multiple
> paragraphs, “each subsequent paragraph in a list item must be indented by
> either 4 spaces or one tab” (emphasis added). However, many implementations
> do not enforce this rule and allow less than 4 spaces of indentation. The
> implementers of Python-Markdown consider it a bug to not enforce this rule.

CommonMark [relax those
restrictions](https://spec.commonmark.org/0.31.2/#lists), allowing a sublist to
be defined with just 2 spaces of indentation.

So I have automated all posts from this blog [using
Go](/posts/2024-07-29/01-quick-bits-why-you-should-automate-everything.md) and
a CommonMark renderer called [Goldmark](https://github.com/yuin/goldmark/). I
them re-render the Markdown to Mataroa using a [Markdown
renderer](https://github.com/teekennedy/goldmark-markdown) before publising to
[capivaras.dev](https://capivaras.dev), because this allow me to do some
transformations in the original Markdown. It mostly works fine except for sub
lists, thanks to the fact that the Markdown renderer I am using renders sub
lists with 2 spaces.

The only reason sub lists are working right now is because
[@ratsclub](https://gluer.org/) fixed this issue in the fork that
[capivaras.dev](https://capivaras.dev) runs. But I want to be compatible with
the official instance if I ever need to migrate.

The solution? Let's fix this in a
[PR](https://github.com/teekennedy/goldmark-markdown/pull/21). However now that
I have code to fix the issue, how can I use it without waiting upstream to
merge my code?

If you are using Go modules it is easy, you just need to use the [`replace`
directive](https://go.dev/ref/mod#go-mod-file-replace):

```go
module github.com/thiagokokada/blog

go 1.23

require (
	github.com/elliotchance/orderedmap/v2 v2.4.0
	github.com/gorilla/feeds v1.2.0
	github.com/gosimple/slug v1.14.0
	github.com/teekennedy/goldmark-markdown v0.3.0
	github.com/yuin/goldmark v1.7.4
	github.com/yuin/goldmark-highlighting v0.0.0-20220208100518-594be1970594
)

require (
	github.com/alecthomas/chroma v0.10.0 // indirect
	github.com/dlclark/regexp2 v1.11.4 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
)

replace github.com/teekennedy/goldmark-markdown => github.com/thiagokokada/goldmark-markdown v0.0.0-20240820111219-f30775d8ed15
```

This will replace all usages of `github.com/teekennedy/goldmark-markdown` to my
fork in `github.com/thiagokokada/goldmark-markdown`. You even get all the
reproducibility of modules since Go automatically pins the commit.

Since the Go format for versions is quite unique, you can just set to the
desired branch (e.g.: instead of `v0.0.0-20240820111219-f30775d8ed15`, you can
use `add-sublist-length-opt` that is the branch name) and run `go mod tidy` to
fix the format.
