# My coding experience with LLM

I don't want to talk about AI. AI fatigue is real, with (garbage) AI content
everywhere you go. However, I like to talk about tools, and AIs, or more
specifically, [Large Language
Models](https://en.wikipedia.org/wiki/Large_language_model) (LLM), can be
really useful tools. Maybe someday we will have our AI overlords dominating
everything (including my job), but right now I don't feel particularly
threatened given my experience with AIs.

Given that [Google Search is getting
worse](https://mashable.com/article/google-search-low-quality-research), I am
slowly relying in LLMs more and more. This is the focus of this post, basically
[the good, the bad and the
ugly](https://www.rottentomatoes.com/m/the_good_the_bad_and_the_ugly) (by the
way, never saw this film) of LLMs, especially focused for solving coding
problems. I will take a look at the my chat history in
[ChatGPT](https://chatgpt.com/) and see where it was a good tool, where it
shined, and where it just completely wasted my time.

A disclaimer, I never bothered to buy a ChatGPT subscription because I don't
feel that the tool is useful enough (at least yet), so most of the examples
here are either with GPT-3.5 or GPT-4o mini, since those where the models
offered by free in ChatGPT.

An addendum, I have some experience with [GitHub
Copilot](https://github.com/features/copilot) (another LLM, but this one runs
inside your IDE or favorite text editor) in `$PREVIOUS_JOB`. But my experience
with it is even more underwhelming than what I am going to describe in this
post: it was basically really good for generating fixture for tests and between
ok-ish to plain bad for everything else. Yes, that's would fit in a tweet (or
xeet?).

## The good

[Chat history](https://chatgpt.com/share/c58a1432-2761-4954-acac-cade29942c34).

> Is there any way that I can pass a context for a net.conn to cancel its
> Read()?

This is one of the good cases where ChatGPT felt to me that I was talking with
an expert in Go. To start, I didn't even mention in my initial prompt that I
was talking about Go, but it correctly inferred so. Also, the solution overview
was great:

> 1. **Use a separate goroutine** to handle the context cancellation.
> 2. **Set a deadline** on the connection using `SetReadDeadline()` to unblock
>    the `Read()` when the context is canceled.

Before asking ChatGPT about this issue I had already coded my own naive
solution that was something similar to this:

```go
func foo(ctx context.Context) {
    // ...
    done := make(chan struct{})
    var n int
    var err error

    go func() {
        n, err = conn.Read(buf)
        done <-struct{}{}
    }()

    select {
    case <-ctx.Done():
        conn.Close()
        return 0, ctx.Err()
    case <-done:
        return n, err
    }
    // ...
}
```

Compared this solution with the one from ChatGPT:

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

The ChatGPT version is much better because:

1. It doesn't close the connection socket allowing it to be reused
2. It waits until the goroutine is finished to avoid leaks
3. It closes the channel after it finishes

Now, I don't claim to be a Go expert (far from it), but this was one of those
cases where I felt that I was being guided by someone that is coding in Go for
years.

But it didn't end there, I asked:

> If I want to reuse the connection, do I need to reset the SetReadDeadline?

Since in this particular case I wanted the underlying connection to be reused,
and ChatGPT correctly amended the code:

```go
func readWithContext(ctx context.Context, conn net.Conn, buf []byte) (int, error) {
    // Same as before
    select {
    case <-ctx.Done(): // Context was canceled or timed out
        // Set a short deadline to unblock the Read()
        conn.SetReadDeadline(time.Now())
        <-done // Wait for the read to finish

        // Reset the deadline so the connection can be reused
        conn.SetReadDeadline(time.Time{})
        return 0, ctx.Err()
    case <-done: // Read finished successfully
        return n, err
    }
}
```

Also the explanation on why this happens was good:

> When you set a read deadline with `conn.SetReadDeadline(time.Now())` to unblock
> the `Read()` call, this deadline remains in effect for future `Read()`
> operations. This means that if you attempt another read operation after this,
> it could immediately fail with a timeout error if the deadline is still in
> effect.
> To reuse the connection, you'll need to remove the deadline or set a new one.
> You can do this by calling `conn.SetReadDeadline(time.Time{})`, which clears
> any existing deadline.

I had one last question that I kind already knew the answer:

> Do we need to defer the clearing?

And ChatGPT correctly returned:

> Whether you need to defer the clearing of the read deadline depends on how
> your code is structured. If you want to ensure that the read deadline is
> always cleared after a context cancellation, even if the code panics or exits
> early, then deferring the `conn.SetReadDeadline(time.Time{})` call is a good
> idea.

This whole interaction was one of the best ones I ever had in ChatGPT, one that
it gave me an actual functional code with good explanation on what was
happening. But still, it didn't talked at all about the gotcha of setting the
read deadline and reusing the underlying connection, needing another prompt to
give me this information.

For comparison,
[this](https://stackoverflow.com/questions/61637359/how-to-use-context-context-with-tcp-connection-read)
is the best that I got searching in both DuckDuckGo and Google, that is where I
kind based my first solution on. I still not sure where ChatGPT could give me
such a good answer for this question, because I just couldn't find anything
remotely similar during my searches.

## The amazing

[Chat history](https://chatgpt.com/share/0c9a148f-4cfd-4b4e-aa84-7486e185125b).

> I want to write an extension to https://github.com/yuin/goldmark/ that allows
> me to rewrite `[link](../2024-07-30/post.md)` to
> `[link](https://kokada.capivaras.dev/slug)`

This one will need some context. [Goldmark](https://github.com/yuin/goldmark/),
the Markdown renderer that [powers this
blog](/posts/2024-08-24/01-making-a-blog-for-the-next-10-years.md), does not
have good examples or documentation (at least I didn't found neither).

So I was not really hopeful I would get any answers in ChatGPT, but not only it
put me in the correct path, the code was almost correct. It did had some unused
and missing imports, and there were some compilation errors that were easily
fixed. Also, it just registered a `html.WithHardWraps()` renderer for no reason
at all.

But this is one of the experiences in LLM that most impressed me. I had no idea
at all on how to solve this problem, and ChatGPT basically gave me the
solution. A [modified version of the
code](https://github.com/thiagokokada/blog/blob/3a2ea58a0638d1005617d25eabeae0626933db46/link_rewriter.go)
from this chat is still used to power this blog to this day, that also
demonstrates that the code is not garbage (since now I understand what it does
and I can even extend it).

I am still not sure how the license situation of this code is though. Does
licenses even make sense in AI generated code?

## The bad

[Chat history](https://chatgpt.com/share/95fcd777-b03c-4cfe-8d1d-12e53f5a4dfe).

> I want to configure VM options inside IntelliJ by project

> My project uses sbt, and I want to set this in a way that I can share with my
> colleagues

> I am using a build.sbt file, but IntelliJ is not loading it's JVM options and
> I need to set this up manually. I want to automate this

This is an example of my most common interaction with ChatGPT: it generally
gives some options that could work, but not in my specific case. So I start to
ask for more information until it eventually starts to
[hallucinate](https://en.wikipedia.org/wiki/Hallucination_(artificial_intelligence)).

In this particular case, the hallucination was in:

> The sbt-idea-plugin can generate IntelliJ IDEA project files based on your
> SBT configuration, including JVM options. This allows you to ensure that
> IntelliJ loads the correct settings from the SBT configuration.
>
> Steps:
>
> 1. Add the `sbt-idea-plugin` to your `project/plugins.sbt` file:
> ```scala
> addSbtPlugin("com.github.mwz" % "sbt-idea" % "2.4.0")
> ```
> 2. Run the following command in your terminal to regenerate IntelliJ IDEA
>    project files:
> ```console
> sbt gen-idea
> ```

Yes, there is a
[`sbt-idea-plugin`](https://github.com/JetBrains/sbt-idea-plugin), but no, it
doesn't generate project files based in my SBT configuration. It is not even
from the repository that it is supposed to be pointing too.

I will just not say that this interaction was a waste of time because some of
the other ideas can be applied, but this was the kind of case where I ended up
just being disappointed because the thing that seemed to be a possible solution
for my problem ended up not existing at all.

## The ugly

There are lots of examples where ChatGPT just plain disappointed me. Here is a
few:

> What is the language that is a fork of Go that has extra features like map?

[Chat history](https://chatgpt.com/share/2b854bbd-d376-4049-97fe-b47321c1ef07).

This is one of those cases where I saw something once and was trying to
remember. ChatGPT just gave me lots of answers that were not a fork of Go
([V](https://vlang.io/), [Gluon](https://github.com/gluon-lang/gluon),
[Pony](https://www.ponylang.io)), hallucinated a programming language (Gomini),
and said that Zig is focused in STEM.

By the way, I eventually remember what the language was:
[Go+](https://goplus.org/). To be fair with ChatGPT, I misremembered that the
Go+ had map (as the [high-order
function](https://en.wikipedia.org/wiki/Map_(higher-order_function)), not a
hashmap that both Go and Go+ have). But still, Go+ explicit announces that it
is a language for STEM everywhere in their website, and even them ChatGPT
couldn't answer.

> In IntelliJ, I need to pass some VM options to the run configuration for my
> application to run tests, e.g.:
> ```
> --add-exports=java.base/sun.security.x509=ALL-UNNAMED
> --add-opens=java.base/sun.security.ssl=ALL-UNNAMED
> -Dconfig.resource=api.test.conf
> ```
> However, the `api.test.conf` may be `jobs.api.conf` or something else
> depending on the subproject I am running inside Idea. Can I parametrize this
> in some way?

[Chat history](https://chatgpt.com/share/aa9bd051-65de-49d9-8d56-0e28de94722c)

So quick answer: you can't. IntelliJ macros are not powerful enough for this.
ChatGPT itself said so when I asked about the IntelliJ macros
[directly](https://chatgpt.com/share/7a057bf5-bcec-458b-9432-98e3f551ad28).

But in this particular case, ChatGPT happily said "yes, you can parameterize VM
options in IntelliJ IDEA to make them more flexible depending on the subproject
or other context". Thanks for wasting my time, I guess.

## Conclusion

I do think that LLMs, in particular ChatGPT since this is the one I have more
experience with, are improving. But still, if there is anything that I
concluded after putting all those experiences in a post is that you still
someone that knows at least as much about whatever you're asking ChatGPT to
actually validate the answer.

Not only because LLMs hallucinates a lot, but also because even in the "good"
example the LLM needed extra input to figure out about the possible drawback of
the solution. The only case where I felt that ChatGPT outsmarted me was the
"amazing" case, but this was literally the only example that I could find in
all my chat history with ChatGPT that was like this.

So yes, at least for now I still don't think AI will be our overlords. Maybe
someday.
