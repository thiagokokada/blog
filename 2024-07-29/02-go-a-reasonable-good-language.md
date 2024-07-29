# Go, a reasonable good language

Go was one of the languages that I always was interested to learn, but never
got the hang of it. I first got interested in the language when I was in my
first job, between 2016-2018. At the time the language was a completely
different beast: no modules, no generics, no easy way to error wrap yet, etc.

Go forward 2023 (no pun indented), I wrote my [first project in
Go](https://github.com/thiagokokada/twenty-twenty-twenty/), wrote some scripts
at `$CURRENT_JOB` in the language, and now wrote [my first
library](https://github.com/thiagokokada/hyprland-go/). I am also writing more
[scripts](https://github.com/thiagokokada/nix-configs/blob/8c559527ed12e1d4f57a3fc5c72630b956f4c290/home-manager/desktop/wayland/hyprland/hyprtabs/hyprtabs.go)
in the language, where I would prefer to use Bash or Python before. Heck, even
this blog is automatically published with a [Go
script](https://kokada.capivaras.dev/blog/quick-bits-why-you-should-automate-everything/),
that used to be a [Python
one](https://kokada.capivaras.dev/blog/using-github-as-a-bad-blog-platform/)
before. I can say that nowadays it is another language in my toolbox, and while
it is still a love and hate relationship, recently it is more about love and
less about hate.

The points that I love about Go is probably obvious for some, but still
interesting to talk about anyway. The fact that the language generates static
binaries by default and have fast compilation times is something that I
apreciate since I first heard about the language, and now that I am using the
language frequently are points I appreciate even more. Something about getting
almost instant feedback after changing a line of code and running `go run`
(even with its quirks) are great for the developer experience. This is the main
reason why I am using the language more frequently for scripts.

Then we have the fast startup times. I am kind of sensitive to latency,
especially of command line utilities that need to answer fast when I expect
them to be fast (e.g.: `foo --help`). This is one part where I could have
issues in Python, especially for more complex programs, but in Go it is rarely
an issue.

Modules are also fantastic. It is not without its weirdness (like everything in
Go ecossystem), but the fact that it is so easy to add and manage dependencies
in a project using only the `go` CLI is great. I also like that it generates a
hash of every dependency, make it reproducible (well, probably not at Nix
level, but still reproducible).

Since I started to talk about `go` CLI, what a great tool! The fact that you
can manage dependencies, generate documentation, format code, lint, run tests,
etc., all with just the "compiler" for the language is excelent. Still probably
one of the best developer experiences I know in any programming language (maybe
only rivaled by [Zig](https://ziglang.org/)).

Now for the parts that I like less, the test part still quirks me that it is
not based in assertions, but thankfully it is easy to write assertions with
generics nowadays:

```go
func Equal[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got: %#v, want: %#v", got, want)
	}
}

func GreaterOrEqual[T cmp.Ordered](t *testing.T, actual, expected T) {
	t.Helper()
	if actual < expected {
		t.Errorf("got: %v; want: >=%v", actual, expected)
	}
}

// etc...
```

Just one of those things that I end up re-writing in every project. Yes, I know
about [testify](https://github.com/stretchr/testify) and other assertion
libraries, but nowadays I believe if you can avoid importing a library, as long
the code you write is trivial, it is better to duplicate the code than try to
import a dependency.

About another piece of code that generics allows me to write and I always end
up re-writing in every project is the `must*` family of functions:

```go
func must(err error) {
	if err != nil {
		panic(err)
	}
}

func must1[T any](v T, err error) T {
	must(err)
	return v
}

// must2, must3, etc...
```

Those functions are so useful, especially for scripts where I generally don't
want to handle each error: if I have an error, I want the program to halt and
print a stack trace (exactly as I would have with a language with exceptions).
It basically allow me to convert code from:

```go
contents, err := os.ReadFile("file")
if err != nil {
    panic(err)
}
```

To:


```go
contents := must1(os.ReadFile("file"))
```

This brings Go closer to Python to me, and I think for scripts this is
something great.

Finally, for the things that I hate, well the biggest one currently is the lack
of nullability (or in Go terms,
[nillability](https://github.com/golang/go/issues/49202)). After using
languages that has it, like Kotlin, or even something like
[mypy](https://www.mypy-lang.org/), this is one of those things that completely
changes the developer experience. I also still don't like the error handling
(but `must*` goes far by improving the situation, when it is possible to use
it), especially because it is easy to lose context on it:

```go
// bad
func readFileContents(file) ([]byte, error) {
    contents, err := os.ReadFile(file)
    if err != nil {
        return nil, err
    }
    return contents, nil
}

// good
func readFileContents(file) ([]byte, error) {
    contents, err := os.ReadFile(file)
    if err != nil {
        return nil, fmt.Errorf("readFileContents: error while reading a file: %w", err)
    }
    return contents, nil
}
```

I also have some grips about the mutate everything approach of the language. I
prefer immutability by default, but I find that in general as long as you split
your functions at a reasonable size it is generally fine.

I expect to write more Go code going forward. Not because it is the perfect
language or whatever, but just because it a is language that has some really
good qualities that makes the language attractive even with the issues that I
have. That makes it a reasonable good language, and at least for me this is
good enough.
