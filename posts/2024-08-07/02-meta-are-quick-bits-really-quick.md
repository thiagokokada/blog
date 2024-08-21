# Meta: are quick bits really quick?

When I wrote my first [quick
bits](/posts/2024-07-27/01-quick-bits-nix-shell-is-cursed.md) post in this
blog, I gave that title prefix without much thought: I knew it was supposed to
be a quicker post than my previous one, but I never thought too deeply about
it. But after:

```console
$ ls -lah **/*quick-bits*.md | wc -l
4
```

Well, 4 blog posts starting with the same prefix, I was curious: are those
quick bits really quick, or at least quicker? Let's see:

```
$ wc -w **/*.md
 1107 2024-07-26/01-writing-nixos-tests-for-fun-and-profit.md
 1220 2024-07-26/02-using-github-as-a-bad-blog-platform.md
  286 2024-07-27/01-quick-bits-nix-shell-is-cursed.md
  387 2024-07-29/01-quick-bits-why-you-should-automate-everything.md
 1060 2024-07-29/02-go-a-reasonable-good-language.md
 1380 2024-07-30/01-first-impressions-fpgbc.md
 1238 2024-07-31/01-generating-yaml-files-with-nix.md
 2308 2024-08-01/01-troubleshooting-zsh-lag-and-solutions-with-nix.md
  504 2024-08-01/02-quick-bits-realise-nix-symlinks.md
  834 2024-08-04/01-make-nixd-modules-completion-work-anywhere-with-flakes.md
 1147 2024-08-05/01-my-favorite-device-is-a-chromebook.md
  394 2024-08-07/01-quick-bits-is-crostini-a-microvm.md
  120 README.md
11985 total
```

While using `wc` is probably not the best way to measure word count (especially
in this blog, since I tend to write lots of code snippets), I think this at
least it gives me a good insight: yes, quick bits are quicker, and they're
basically posts with a soft limit around 500 words. So expect in future this
limit to be used.

By the way, at:

```console
$ wc -w 2024-08-07/02-meta-are-quick-bits-really-quick.md
220 2024-08-07/02-meta-are-quick-bits-really-quick.md
```

This post is also technically a quick bits post, but "quick bits meta" would be
too much. And yes, that last block of code is also meta ;).
