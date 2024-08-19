# Quick bits: Go automatically downloads a newer toolchain if needed

I am using
[elliotchance/orderedmap](https://github.com/elliotchance/orderedmap/) as my
choice of ordered maps (since Go [doesn't have
one](/2024-08-17/01-an-unordered-list-of-things-i-miss-in-go.md) in standard
library yet). I recently did a
[PR](https://github.com/elliotchance/orderedmap/pull/41) to implement Go 1.23
iterators, because I find them neat, however I was avoiding to use it in the
code that generates this [blog](https://github.com/thiagokokada/blog) since Go
1.23 was just released and is still not the default Go in
[nixpkgs](https://github.com/NixOS/nixpkgs).

I decided that I would create a
[branch](https://github.com/thiagokokada/blog/pull/2) and leave there for a few
months, until I decided to try to run the code locally and got this:

```console
$ go build
go: downloading go1.23.0 (darwin/arm64)
```

Nice. And before you ask, yes, the compiled binary works perfectly:

```console
$ make
./blog > README.md
./blog -rss > rss.xml
```

So how does this work? Take a look at the documentation in the official [Golang
page](https://tip.golang.org/doc/toolchain):

> Starting in Go 1.21, the Go distribution consists of a go command and a
> bundled Go toolchain, which is the standard library as well as the compiler,
> assembler, and other tools. The go command can use its bundled Go toolchain
> as well as other versions that it finds in the local PATH or downloads as
> needed.

There are a bunch of rules here that I am not going to enter in detail (I
recommend you to read the official documentation), but a quick summary:

- Go will download a toolchain when either `go` or `toolchain` lines `go.mod`
is set to a Go version higher than your current `go` binary
  + But only if your `go` binary is at least version 1.21, since this is the
  version that introduces this behavior
- You can force a specific toolchain with `GOTOOLCHAIN` environment setting,
e.g.: `GOTOOLCHAIN=1.23`
  + The default value for `GOTOOLCHAIN` is `auto`, that basically has the
  behavior described in this post
  + You can also set to `local` to always use the current `go` binary, or the
  previous behaviour pre-1.21 Go
  + There is also `<name>+auto` and `path` options, that can be seen in the
  docs
- The downloaded toolchains go to whatever your `GOPATH` is, inside
`golang.org/toolchain` module, and version `v0.0.1-goVERSION.GOOS-GOARCH`, for
example:

  ```console
  $ ls -lah $GOPATH/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.0.darwin-arm64
  total 64
  dr-xr-xr-x  14 user  staff   448B 19 Aug 12:01 .
  drwxr-xr-x   4 user  staff   128B 19 Aug 12:01 ..
  -r--r--r--   1 user  staff   1.3K 19 Aug 12:01 CONTRIBUTING.md
  -r--r--r--   1 user  staff   1.4K 19 Aug 12:01 LICENSE
  -r--r--r--   1 user  staff   1.3K 19 Aug 12:01 PATENTS
  -r--r--r--   1 user  staff   1.4K 19 Aug 12:01 README.md
  -r--r--r--   1 user  staff   426B 19 Aug 12:01 SECURITY.md
  -r--r--r--   1 user  staff    35B 19 Aug 12:01 VERSION
  dr-xr-xr-x   4 user  staff   128B 19 Aug 12:01 bin
  -r--r--r--   1 user  staff    52B 19 Aug 12:01 codereview.cfg
  -r--r--r--   1 user  staff   505B 19 Aug 12:01 go.env
  dr-xr-xr-x   3 user  staff    96B 19 Aug 12:01 lib
  dr-xr-xr-x   4 user  staff   128B 19 Aug 12:01 pkg
  dr-xr-xr-x  77 user  staff   2.4K 19 Aug 12:02 src
  ```

By the way, this only works well because Go binaries are static, one of the
things that make the language [reasonable
good](/2024-07-29/02-go-a-reasonable-good-language.md).

While I don't like a program downloading random binaries from the internet, I
like what Go is doing here. It makes the whole bootstrapping process for a Go
project much easier: as long as you have a reasonable up-to-date `go` binary in
your `PATH`, you should be ready to go (pun intended). And Go modules are
already reasonable secure, ensuring that each module have a proper checksum. As
long as nobody else can publish modules in `golang.org/toolchain` namespace I
can't see much of a security issue here, but I am not a security expert.

But if you don't like this behavior, you can always disable it by setting
`GOTOOLCHAIN=local`. And just do not forget to set this in your
[CI](https://brandur.org/fragments/go-version-matrix), unless you don't care
about Go versions.
