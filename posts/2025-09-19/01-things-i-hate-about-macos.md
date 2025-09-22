# Things I hate about macOS

I have a kind of love and hate relationship with my work laptop, a MacBook Pro
M2 Pro (what a weird name). I love almost everything about the hardware, from
the premium materials, to the amazing ProMotion screen, to one of the best
speakers I have in my house (no kidding, those speakers are better than some of
my dedicated Bluetooth speakers) and the best touchpad I ever used, period. But
I hate macOS: after using it for work for the last 3 years I got used to its
quirks, but the experience is just so... bad. I could never really point why I
think so, this is why I decided to write about it.

To start: using terminal on macOS just feels slow. Everything is slower in
macOS compared to my Linux desktop, and this shouldn't be a hardware issue
since my MacBook Pro is probably more powerful and has faster I/O than my
desktop.

From opening a new terminal and even typing seems slower (like there is
something running in background every time I press a key). I am using the
exactly same terminal ([Kitty](https://sw.kovidgoyal.net/kitty/)) and
configuration in both. Just a quick and non-scientific benchmark:

```
$ hyperfine 'zsh -ic exit' # Linux
Benchmark 1: zsh -ic exit
  Time (mean ± σ):      94.1 ms ±   2.4 ms    [User: 60.2 ms, System: 34.5 ms]
  Range (min … max):    90.0 ms …  99.0 ms    29 runs


$ hyperfine 'zsh -ci exit' # macOS
Benchmark 1: zsh -ci exit
  Time (mean ± σ):     233.0 ms ± 180.9 ms    [User: 54.6 ms, System: 51.0 ms]
  Range (min … max):   153.0 ms … 746.5 ms    10 runs

  Warning: The first benchmarking run for this command was significantly slower
  than the rest (746.5 ms). This could be caused by (filesystem) caches that
  were not filled until after the first run. You should consider using the
  '--warmup' option to fill those caches before the actual benchmark.
  Alternatively, use the '--prepare' option to clear the caches before each
  Timing run.

$ hyperfine 'zsh -ic exit' # Chromebook
Benchmark 1: zsh -ic exit
  Time (mean ± σ):     393.1 ms ±  24.7 ms    [User: 136.8 ms, System: 270.8 ms]
  Range (min … max):   357.0 ms … 430.6 ms    10 runs
```

This may look like a unfair comparison because it seems that I run the macOS
tests with cold cache on purpose (and this is why `hyperfine` recommended me to
use `--warmup` flag), while I run the Linux tests with a hot cache. However it
is not, this basically matches my experience with macOS where it seems the file
cache expires much faster than on Linux. So while on Linux I rarely see ZSH
taking time to start, it is a common occurrence in macOS. But even ignoring
this issue macOS in general seems to be much slower, and this is not isolated
to my `zsh`, almost every binary inside my terminal seems to start slower.

I also add the results from my
[Chromebook](/posts/2024-08-05/01-my-favorite-device-is-a-chromebook.md). It is
much slower than both my Linux desktop and my macOS system, but this is
expected considering that both the CPU and I/O is much slower (this device
still uses an [eMMC](https://en.wikipedia.org/wiki/MultiMediaCard#eMMC), that
in some metrics is slower than a HDD). But also the results are much more
consistent, again matching what is my experience with macOS: the system is just
inconsistent slow sometimes.

Now let's look out of the terminal and more for the desktop part. One of my
major grips about the system is the lack of choice. For example, I want to set
my touchpad to use natural (or reverse) scrolling, since well, this is what we
got used after the smartphone boom. But I also want my scroll wheel to use
"normal" scrolling, since this is what years of using a mouse with scroll made
me used to. This is easy to do in any other operating system that it is not
macOS. And the worst thing is that macOS is even deceitful:

[![Mouse](/posts/2025-09-19/Screenshot_2025-09-19_at_13.44.35.png)](/posts/2025-09-19/Screenshot_2025-09-19_at_13.44.35.png)

[![Trackpad](/posts/2025-09-19/Screenshot_2025-09-19_at_13.47.16.png)](/posts/2025-09-19/Screenshot_2025-09-19_at_13.47.16.png)

So you see in the above screenshots that both Mouse and Trakcpad have separate
options for setting "Natural scrolling", but this is a lie: if you change one
of them it changes both, and there is nothing to indicate this. This is so
confusing that before I knew this I would change one option to fix the current
input device that I was using, only later to realise trying to use my other
input device that it would scroll the opposite that I expect, so I would "fix"
again, rise and repeat.

To fix this issue? As far I know, only using an external program. I use [Linear
Mouse](https://linearmouse.app/), that to be clear, it is a great program. It
is just that I shouldn't need to use it, and thanks to the way it works (it
uses Accessibility APIs as far I know) sometimes things get wonky and stops
working.

Another example where macOS refuses to give you choices? Since I have a MacBook
Pro, it has a Touch ID and it works great. Except that I can't use it with a
close lid. No problem, I can just keep the lid of the laptop open. But in macOS
if I keep the lid open I can't turn off the internal display. I don't want that
display to be turned on though, not only it is a waste of energy but also
it means that my mouse can sometimes go to a screen that I am not even
using and this is jarring. The solution? Even another external program:
[BetterDisplay](https://github.com/waydabber/BetterDisplay).

Again, nothing against BetterDisplay that is a really good program. It is just
that I shouldn't need it for something so basic as disabling the internal
screen when I am using an external monitor. BetterDisplay has way more
features, but currently this is the only one I use. The fact that I had to pay
€19.99 for the luxury of turning off the internal display is infuriating.

By the way, talking about multi-monitor support, another grip. I like to use
the dock on the side of the monitor because this makes for better vertical
space (especially good considering that my main monitor is a Ultrawide one, so
I have lots of horizontal space but low amount of vertical space). However, if
I set the dock to the side, it will go to whatever monitor is at that side.
What? Yes, even if my main monitor is setup as the "Main display", if I set my
dock to the left and my laptop is on the left side, the dock will go there.

This is one of the things that I don't have a good solution. My solution was to
eventually just reorganize my whole desk to always ensure that my laptop will
go to the right so I can have the dock on the left side as I want. Yes, instead
of making the operational system works for me, I need to make my desk work with
my laptop.

And of course, there are the bugs. Now to be clear, bugs happens in every
operational system that I know, it is just the way that modern systems works
nowadays: they're too complex, and complexity introduces bugs. But bugs that
completely stop whatever I am doing bother me way more, and macOS seems to have
lots of them. Let me introduce you one of them: the notification of death.

[![I present you the notification of death](/posts/2025-09-19/Screenshot_2025-09-19_at_17.31.34.png)](/posts/2025-09-19/Screenshot_2025-09-19_at_17.31.34.png)

This innocent notification bought me lots of dread. For months every time I
tried to click this "Allow" I would basically lose all my work: video and audio
would continue playing and my mouse would still move, but I couldn't click on
anything or type. The system would respond to a power button press (locking the
screen), but 99% of the time if I unlocked I would go back to the same state.
My only option would be to force turn off the system and restart. It was like
my mouse focus was in a invisible screen that wouldn't want to release the
focus. And yes, I tried everything I could thing, like for example Cmd+Tab to
change the current application focus.

To be clear, it seems that this issue is finally fixed in macOS Tahoe, at least
I couldn't reproduce this issue while writing this blog post. But the reason
this bug ever happened is wild, there is no reason why an application could
steal the focus of the input and not give it back. Also, this was not the only
"stop the world" bugs that I had with macOS (I just had one last week while
doing random things), it is just the one that I knew how to reproduce.

So that is basically why I feel so strong against macOS. Windows 11 is probably
the worse of the two, but since I generally can do what I want it seems that I
feel less strong about the system. And yes,
[KDE](/posts/2025-09-17/01-kde-is-now-my-favorite-desktop.md) is far ahead of
the two as my favorite desktop, since it tries to embrace whatever I want to
do.
