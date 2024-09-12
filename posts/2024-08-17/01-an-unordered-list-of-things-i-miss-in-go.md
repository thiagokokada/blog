# An unordered list of things I miss in Go

I like Go. I think it is a [reasonable good
language](/posts/2024-07-29/02-go-a-reasonable-good-language.md), and has some
good qualities that makes up for its flaws. However, this doesn't mean I think
the language couldn't be better, far from it.

This blog post is a list of things that I miss from Go from other languages.
Some of the things here could probably be implemented soon, some other would
probably need a major revision of the language. The list is unordered, because
this makes it easier for me to update in the future if I found something else,
but also because I don't want to think too hard about giving each point here a
rank.

With all above, let's start.

## Ordered maps in standard library

When I first learned about
[dictionaries](https://docs.python.org/3/library/stdtypes.html#typesmapping) in
Python it quickly became one of my favorite data structures ever. They're
extremely versatile, and most modern programming languages have something
similar in its standard library. Go isn't different, it has
[`map`](https://go.dev/blog/maps), that is Go implementation of a [hash
table](https://en.wikipedia.org/wiki/Hash_table). However `map` in Go are
quirky, for example:

```go
package main

func main() {
	m := map[string]bool{"foo": true, "bar": false, "baz": true, "qux": false, "quux": true}

	for k := range m {
		println(k)
	}
}
```

```console
$ go run ./test.go
bar
baz
qux
quux
foo

$ go run ./test.go
foo
bar
baz
qux
quux

$ go run ./test.go
qux
quux
foo
bar
baz
```

Now, I don't expect any hash table implementation to keep the order of the
elements, but Go actually [randomise each map
instance](https://victoriametrics.com/blog/go-map/):

> But here’s the deal, while the hash function used for maps in Go is
> consistent across all maps with **the same key type**, the `seed` used by
> that hash function is different for each map instance. So, when you create a
> new map, Go generates a random seed just for that map.

While I understand the reason for this (i.e.: to avoid developers relying in a
specific iteration order), I still find it weird, and I think this is something
unique for Go. This decision means that even if you don't care about a specific
order, you will still need to sort the map before doing something else if you
want reproducibility, something that I care a lot.

The fix for this? Go could offer an ordered map implementation inside the
standard library. An ordered map ensure that the iteration order of the map is
the same as the insertion order (that is, by the way, a powerful property that
allow maps to be used in other contexts, not just my pet peeve above).

Python actually does this for any dictionaries since [Python
3.6](https://stackoverflow.com/a/39980744), but it offered an
[OrderedDict](https://docs.python.org/3/library/collections.html#collections.OrderedDict)
before it (and `OrderedDict` still has some methods that normal `dict` doesn't,
that maybe useful in specific cases).

Before generics it would be impossible to have a type-safe API for such data
structure without introducing a new data type in the language (like `slices`),
but now Go has generics so it is not an issue anymore. The other issue is that
you would be forced to iterate manually in this new data structure, but thanks
to the new [`range-over-func`](https://tip.golang.org/doc/go1.23#language) in
Go 1.23, it means we can iterate in an ordered map as a library almost exactly
like we can do as a `map`:

```go
import "orderedmap"

func main() {
    m := orderedmap.New[string, bool]()
    m.Set("foo", true)
    m.Set("bar", false)
    m.Set("baz", true)

    for k := range m.Iterator() {
        println(k) // Order always will be: foo, bar, baz
    }
}
```

Now, of course the lack of Ordered Map in the standard library can be filled
with third party implementations, e.g.: I am using this
[one](https://github.com/elliotchance/orderedmap) in one of my projects. But
being in standard library reduces the friction: if there was some
implementation in standard library, I would generally prefer it unless I have
some specific needs. However when the standard library doesn't offer what I
need, I need to find it myself a suitable library, and this ends up taking time
since generally there are lots of alternatives.

## Keyword and default arguments for functions

Something that comes straight from Python that I miss sometimes in Go is that
you can do things like this when declaring a function:

```python
def hello(name="World"):
    print(f"Hello, {name}")

hello("Foo") # "normal" function call
hello(name="Bar") # calling with keyword arguments
hello() # calling with default arguments
```

```console
$ python hello.py
Hello, Foo
Hello, Bar
Hello, World
```

The lack of default arguments especially affects even some of the API decisions
for Go standard library, for example, `string.Replace`:

> ```func Replace(s, old, new string, n int) string```
>
> Replace returns a copy of the string s with the first n non-overlapping
> instances of old replaced by new. If old is empty, it matches at the
> beginning of the string and after each UTF-8 sequence, yielding up to k+1
> replacements for a k-rune string. If n < 0, there is no limit on the number
> of replacements.

If Go had default arguments, `Replace` could have e.g.: `func Replace(s, old,
new string, n int = -1)` signature, that would mean we don't need
`strings.ReplaceAll` (since this is basically what this function does, call
`strings.Replace(s, old, new, -1)`).

## Nullability (or nillability)

I talked I little about this in [my previous post about
Go](/posts/2024-07-29/02-go-a-reasonable-good-language.md), but I want to
expand here.

First, I don't think the language needs to support the generic solution for
nullability, that would be either having proper Union or Sum types. Kotlin
AFAIK doesn't support neither, but my 2 years experience with Kotlin showed
that just having nullable types already helped a lot in ensuring type safety.

Second, I do feel that Go has less issues with `nil` values, than say, Java,
because its decision of using zero values instead of `nil` in many cases. So
for example, a string can never be `nil`, however a string pointer can be. This
means that this is fine:

```go
func(s string) {
    // do something with s
}
```
However:

```go
func(s *string) {
    // s maybe nil here, better check first
}
```

Still, I get more `panic` for `nil` pointer deference than I get in other
languages that offer nullables (heck, even Python with
[`mypy`](https://www.mypy-lang.org/) is safer).

Sadly this is the change in this post that is more likely to need a completely
new revision of the language.
[nillability](https://github.com/golang/go/issues/49202) was proposed before,
but it is really unlikely it can be done without breaking backwards
compatibility.

It could be done the Java way by adding a `nullable` type to the standard
library ([JSR305](https://jcp.org/en/jsr/detail?id=305)), but the fact that
[JSR305 is considerd
dead](https://stackoverflow.com/questions/2289694/what-is-the-status-of-jsr-305)
by many shows how difficult it is to do something like this without a major
change in the language. Dart is the only language that I know that [did this
successfully](https://dart.dev/null-safety/understanding-null-safety), but
definitely it was not without its pains. And the fact that most people that
program in Dart probably does because of Flutter (that eventually required
newer versions with null-safety) is not a good sign.

## Lambdas

_Added in 2024-08-18_

Go is a surprising good language for some functional code, thanks to having
first class functions and closures. Sadly the syntax doesn't help, since the
only way you can use anonymous functions in Go is using `func`. Especially if
the types are complex, this can result in some convoluted code. Take the
example from the [`range-over-func`
experiment](https://go.dev/wiki/RangefuncExperiment):

```go
package slices

func Backward[E any](s []E) func(func(int, E) bool) {
    return func(yield func(int, E) bool) {
        for i := len(s)-1; i >= 0; i-- {
            if !yield(i, s[i]) {
                return
            }
        }
    }
}
```

If Go had a syntax for lambdas, especially if we could elide the types, this
could be simplified a lot:

```go
package slices

func Backward[E any](s []E) func(func(int, E) bool) {
    return (yield) => {
        for i := len(s)-1; i >= 0; i-- {
            if !yield(i, s[i]) {
                return
            }
        }
    }
}
```

Or even something like this would already help, no special syntax but allowing
the types to be elided in an unnamed function:

```go
package slices

func Backward[E any](s []E) func(func(int, E) bool) {
    return func(yield) {
        for i := len(s)-1; i >= 0; i-- {
            if !yield(i, s[i]) {
                return
            }
        }
    }
}
```

This feature I am still somewhat hopeful that may become a reality in some
future version of the language, since they didn't close the
[issue](https://github.com/golang/go/issues/21498) yet, and the discussion
about the possibility of this feature is still ongoing.

## Error on unused return values

_Added in 2024-09-12_

In a [recent post about
LLMs](/posts/2024-09-02/01-my-coding-experience-with-llm.md), I talked about
this function that I asked ChatGPT to generate to create a socket that accepted
a [context](https://pkg.go.dev/context):

```go
func readWithContext(ctx context.Context, conn net.Conn, buf []byte) (int, error) {
    done := make(chan struct{})
    var n int
    var err error

    // Start a goroutine to perform the read
    go func() {
        n, err = conn.Read(buf)
        close(done)
    }()

    select {
    case <-ctx.Done(): // Context was canceled or timed out
        // Set a short deadline to unblock the Read()
        conn.SetReadDeadline(time.Now())
        <-done // Wait for the read to finish
        return 0, ctx.Err()
    case <-done: // Read finished successfully
        return n, err
    }
}
```

However while reading [this blog
post](https://joeduffyblog.com/2016/02/07/the-error-model/) about "The Error
Model", I realised that this function is lacking error check during
`conn.SetReadDeadline()` calls:

```go
func readWithContext(ctx context.Context, conn net.Conn, buf []byte) (n int, err error) {
	done := make(chan struct{})

	// Start a goroutine to perform the read
	go func() {
		n, err = conn.Read(buf)
		close(done)
	}()

	select {
	case <-done:
		return n, err
	case <-ctx.Done():
		// Set a short deadline to unblock the Read()
		err = conn.SetReadDeadline(time.Now())
		if err != nil {
			return 0, err
		}
		// Reset read deadline
		defer func() {
			if e := conn.SetReadDeadline(time.Time{}); e != nil {
				err = errors.Join(err, e)
			}
		}()
		// Make sure that the goroutine is done to avoid leaks
		<-done
		return 0, errors.Join(err, ctx.Err())
	}
}
```

I could blame the LLM on this, but this is the type of mistake that I could see
happening even in a Pull Request done by an actual human.

Ignoring errors is bad, especially because once the error is ignored it is
gone, forever. You may have strange issues in the code that are impossible to
debug because you can't know about something that doesn't exist. This is one of
the reasons I think exceptions are better, since this is the kind of error that
would be impossible to ignore if the language had exceptions, and the exception
would eventually propagate in your stack until it is captured or you get a
crash (with a stack trace to debug the issue).

Now, I don't think this is a fault of using error as values instead of
exceptions. Go has the tools to encapsulate the errors and propagate them
properly, like it is shown in the fixed code. But the fact that someone can do
this mistake without any warning or error from the compiler is bad.

If Go had a warning or error for unused return values, this would be different:

```go
func foo(conn *net.Conn) {
  // ...
  _ = conn.SetReadDeadline(time.Now())
}
```

Now in this case it is clear: I am ignoring the error explicitly, probably for
a good reason. I can ask during a Pull Request why the committer is ignoring if
it is lacking appropriate context, and maybe even ask for some comments to be
added why this would be safe. What can't happen is this being ignored by
mistake.

I am not completely sure if I want this only for errors or for any unused value
though. There are lots of cases where ignoring non-error values is ok, but I
also don't see the issue of having something like:

```go
func foo() {
  // ...
  _, _ = FuncThatReturnIntAndError()
}
```

If anything, it is clear for whoever is reading this code later that you're
only calling this function only for its side-effects.

By the way, yes, [errcheck](https://github.com/kisielk/errcheck) linter exists,
but if the language itself doesn't enforce this it means that there are lots of
libraries that are probably mistakenly ignoring errors. And if the library is
ignoring errors by mistake, you can't do much when you actually have them.
