# Praise to scdoc to generate man pages

Hey, its been a long time since my [last blog
post](posts/2024-10-07/01-enabling-le-audio-lc3-in-wf-1000xm5.md). It is mostly
because I ran out of things to write, but I expected this. This is probably
more likely how I am actually going to post from now. At least, it shows that
my plan to have a [blog for a long
time](/posts/2024-08-24/01-making-a-blog-for-the-next-10-years.md), that is
easy to go back when I wanted is working fine, but I digress.

Going back to the theme of the today blog post, I needed to write a [man
page](https://en.wikipedia.org/wiki/Man_page) for the first time in years. I
hate [troff](https://en.wikipedia.org/wiki/Troff), the typesetting system used
for man pages (similar to [LaTeX](https://en.wikipedia.org/wiki/LaTeX) for
documents). It is one of the weirdest languages that I ever saw, and even the
example in Wikipedia shows that:

```troff
.ND "January 10, 1993"
.AU "Ms. Jane Smith"
.AT "Upcoming appointment"
.MT 5
.DS
Reference #A12345
.sp 4
Mr. Samuel Jones
Field director, Bureau of Inspections
1010 Government Plaza
Capitoltown, ST
.sp 3
Dear Mr. Jones,
.sp 2
.P
Making reference to the noted obligation to submit for state inspection our newly created production process, we request that you consider the possible inappropriateness of subjecting the  innovative technologies of tomorrow to the largely antiquated requirements of yesterday.  If our great state is to prosper in the twenty-first century, we must take steps
.B now ,
in
.I this
year of
.I this
decade, to prepare our industrial base for the interstate and international competition that is sure to appear.  Our new process does precisely that.  Please do not let it be undone by a regulatory environment that is no longer apt.
.P
Thank you for your consideration of our position.
.FC Sincerely
.SG
```

Keep in mind that the break lines are necessary every time you introduce a
macro, like `.I this` (that I _think_ it is for italics). Yes, this format is
as illegible as hell, and it is worse that the format lacks good tooling (or at
least I didn't find any good ones).

Most people when they need to write a man page nowadays ends up using some
other format that generates a man page. For example, in the past I used
[Pandoc](https://pandoc.org/) to convert Markdown to a man page, but even if
Pandoc is a great project the result is sub-optimal at best: Markdowns are, at
the end, designed for generating HTML (and a subset of it), and not man pages,
so you basically ends up fighting the format for it to do what you want.
Also, Pandoc is a big project, with a ~200MB binary (at least it is the default
Pandoc binary in Nix).

For this specific project I needed something small. I am trying to replace one
of the most essential pieces inside NixOS, `nixos-rebuild`, written in Bash,
with a [full rewritten in
Python](https://discourse.nixos.org/t/nixos-rebuild-ng-a-nixos-rebuild-rewrite/55606/)
(sorry Rust zealots!), called `nixos-rebuild-ng`.

Since this project will eventually (if successful) be in the critical path for
NixOS, I want to reduce the number of dependencies as much as possible, so
something as big as Pandoc is out. I could use
[AsciiDoc](https://asciidoc.org/), but it is a big complicated Python project
(this may seem ironic, but `nixos-rebuild-ng` has only one runtime dependency,
that is optional). And I also hated the last time I tried to use it to generate
man pages: it more flexible than Markdown, but still far from optimal.

Thanks to Drew DeVault (creator of [SwayWM](https://swaywm.org/)) that seems it
had the same issues in the past and created
[`scdoc`](https://drewdevault.com/2018/05/13/scdoc.html), a very simple man
page generator using a DSL inspired in Markdown, but specific to generate man
pages. The binary is written in C (and advantage in this case since it means it
is easier to bootstrap), is small (~1 Kloc) and has no dependencies, so it
fits the requirement.

While the language suffers from being a niche project for a niche segment, the
[man page](https://man.archlinux.org/man/scdoc.5.en) for it is actually really
nice. It is terse though and lacks examples, and this is what this blog post
will try to accomplish.

To start, let's have a quick summary of the syntax, written in `scdoc` as
comments:

```scdoc
; quick summary:
; # new section
; comments starts with ;
; - this is a list
; 	- sub-list
; - *bold*: _underline_, force a line break++
; - [tables], \[ can be used to force an actual [
; . numbered list
; please configure your editor to use hard tabs
; see `man 5 scdoc` for more information about syntax
; or https://man.archlinux.org/man/scdoc.5.en
```

I actually added this summary in the `.scd` (the `scdoc` extension) files that
I wrote, so it is easy for someone that never saw the format to start
collaborating.

And here an example of a (summarised) man page in `.scd` format:

```markdown
nixos-rebuild-ng(8)

# NAME

nixos-rebuild - reconfigure a NixOS machine

# SYNOPSIS

_nixos-rebuild_ \[--upgrade] [--upgrade-all]++
		\[{switch,boot}]

# DESCRIPTION

This command has one required argument, which specifies the desired operation.
It must be one of the following:

*switch*
	Build and activate the new configuration, and make it the boot default.
	That is, the configuration is added to the GRUB boot menu as the
	default menu entry, so that subsequent reboots will boot the system
	into the new configuration. Previous configurations activated with
	nixos-rebuild switch or nixos-rebuild boot remain available in the GRUB
	menu.

*boot*
	Build the new configuration and make it the boot default (as with
	*nixos-rebuild switch*), but do not activate it. That is, the system
	continues to run the previous configuration until the next reboot.

# OPTIONS

*--upgrade, --upgrade-all*
	Update the root user's channel named 'nixos' before rebuilding the
	system.

	In addition to the 'nixos' channel, the root user's channels which have
	a file named '.update-on-nixos-rebuild' in their base directory will
	also be updated.

	Passing *--upgrade-all* updates all of the root user's channels.

See the Nix manual, *nix flake lock --help* or *nix-build --help* for details.

# ENVIRONMENT

NIXOS_CONFIG
	Path to the main NixOS configuration module. Defaults to
	_/etc/nixos/configuration.nix_.

# FILES

/etc/nixos/flake.nix
	If this file exists, then *nixos-rebuild* will use it as if the
	*--flake* option was given. This file may be a symlink to a
	flake.nix in an actual flake; thus _/etc/nixos_ need not be a
	flake.

# AUTHORS

Nixpkgs/NixOS contributors
```

And here is a screenshot of the result:

[![Man page rendered from scd
file](/posts/2024-12-04/2024-12-04-230955_hyprshot.png)](/posts/2024-12-04/2024-12-04-230955_hyprshot.png)

One of nice things that I found is how looking at the plain text looks kind
like the man page result already. And if you know Markdown, you can basically
understand most things that is happening. There are a few differences, like
`*bold*` instead of `**bold**`, and while they're unfortunate they're not the
end of the world.

Now, the format has its quirks. The first line being the name of the program
and section in parenthesis is required, but this makes sense, since you need
this information for the corners. But for one, it requires the usage of hard
tabs to create indentation, and the error messages are awful, in a situation
that kind remembers me of `Makefile`. Also the choice of `[` to start a table
means that the traditional `app [command]` needs in many cases to be escaped as
`app \[command]`. I found this a strange choice since this is supposed to be a
format that is only used for man pages, and using `[command]` to indicate an
optional is common, but at least it is easy to escape.

In the end, I think all that matters is the result. And for the first time for
all those years trying to write a man page, I am satisfied with the result. The
man page looks exactly as I wanted once rendered, and the `.scd` file looks
reasonable good that it can work as a documentation for someone that for one
reason or another can't use the man page (can't say the same for the troff
version). Also, it is really easy for someone to just go there and update the
man page, even without experience in the format (except for maybe the
requirement of tabs). So all in all, I really liked the format, and will use it
again if I need to write another man page in the future.
