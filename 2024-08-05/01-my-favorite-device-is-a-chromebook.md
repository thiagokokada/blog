# My favorite device is a Chromebook

Most of the blog posts in this blog (including this one) and most of I would
call "personal computing" that I do nowadays is mostly done in one of the most
unremarkable devices that I own: a cheap [Chromebook Duet
3](https://www.lenovo.com/us/en/p/laptops/lenovo/lenovo-edu-chromebooks/ideapad-duet-3-chromebook-11-inch,-qlc/len101i0034),
that I bought for around EUR300. I was thinking why, because it is woefully
underpowered: a [Snapdragon 7c Gen
2](https://www.anandtech.com/show/16696/qualcomm-announces-snapdragon-7c-gen-2-entrylevel-pc-and-chromebook-refresh),
a CPU that was already considered entry level 3 years ago, coupled with an eMMC
for storage, that is not much fast than a HDD; at least I have the 8GB RAM
version instead of the 4GB one.

It is a hybrid device, one that can be used as either a tablet or laptop, but
is compromised experience in both cases: as a tablet, it lacks the better touch
optimised interface from iOS or Android; as a laptop, you have to depend on the
stand to adjust the screen, and the detachable keyboard is worse then any
laptop I have ever owned: getting keys stucked and characters being duplicated
 as a result is a common occurence.

But I really love this device. It is generally the only device that I bring in
trips nowadays, because while it is compromised it works well enough: I can use
to consume media in tablet mode (the fact that ChromeOS supports Android apps
is a plus in those cases), browse the web and even do Linux stuff (more about
this later). The fact that it is small (the size remembers me of a
[netbook](https://en.wikipedia.org/wiki/Netbook)), lightweight (~1KG, including
the keyboard), has a good screen (that is bright and also HiDPI) and good
battery life (I don't have numbers but I almost never worry about it) is what
makes this device the perfect companion to trips.

Also, it has 2 USB-C ports and supports DisplayPort alt-mode (sadly not at
1440p, the maximum I got was 1080p, altough some people at Reddit [seems to be
successful](https://www.reddit.com/r/chromeos/comments/zh27tg/comment/izku724/?utm_source=share&utm_medium=web3x&utm_name=web3xcss&utm_term=1&utm_content=share_button)
at this resolution; it may be my Dell S3423DWC being Ultrawide or the cable,
who knows?), so it means you can charge it, connect to a external display and
peripherals, all at the same time.

ChromeOS is also really interesting nowadays. Being a Chrome-first OS makes it
a compromised experience, for example, it is the only device that I use Chrome
as my main browser (since I personally prefer Firefox). But having a OS that
boots fast is great: I never worry about OS updates because I know the device
will be ready in seconds after a reboot. And the whole desktop experience
inside the ChromeOS desktop is good, having shortcuts for many operations so
you can get things done fast, and support for virtual desktops (ChromeOS call
it "desks") means you can organise your window as much as you want.

And what I think makes ChromeOS really powerful is
[Crostini](https://chromeos.dev/en/linux), a full Linux VM that you can run
inside ChromeOS. It runs Debian (it seems you can [run other
distros](https://www.reddit.com/r/Crostini/wiki/howto/run-other-distros/)
though) with a deep integration with ChromeOS, so you can run even graphical
programs without issues (even OpenGL works!):

![Fastfetch inside Crostini with gitk running side-by-side.](/2024-08-05/Screenshot_2024-08-05_21.22.29.png)

![Running glxgears inside Crostini.](/2024-08-05/Screenshot_2024-08-05_21.39.58.png)

This is all thanks to
[Sommelier](https://chromium.googlesource.com/chromiumos/platform2/+/HEAD/vm_tools/sommelier/README.md),
a nested Wayland compositor that runs inside Crostini and allow both Wayland
and X11 applications to be forwarded to ChromeOS. The integration is so good
that I can even run Firefox inside Crostini and works well enough, but sadly
Firefox is too slow in this device (I am not sure if the issue is ChromeOS or
Firefox, but I suspect the later since Google does some optimisation per
device).

One interesting tidbit about the OpenGL situation in this device: this seems to
be the first Chromebook to ship with open source drivers, thanks to Freedreno.
There is [this](https://www.youtube.com/watch?v=8mnjSmN03VM) very interesting
presentation done by Rob Clark in XDC 2021, that I recommended anyone
interested in free drivers to watch (the reference design of Duet 3 is
[Strongbad](https://chromeunboxed.com/chromebook-tablet-snapdragon-7c-homestar-coachz-strongbad)).

The Crostini integration is probably the best VM integration with Linux I ever
saw in an OS: you can manage files inside the VM, share directories between the
OS and VM, copy and paste works between the two, GUI applications installed
inside the VM appear in the ChromeOS menu, memory allocation inside the VM is
transparent, etc. Even the themes for applications are customised to match
ChromeOS. It is unironically one of the best Linux desktop experiences I ever
had.

Of course I am using Nix, but since the Crostini integration depends in some
services, I decided to run Nix inside Debian instead of NixOS and run
[Home-Manager
standalone](https://nix-community.github.io/home-manager/index.xhtml#sec-install-standalone).
I recommend checking the official [NixOS Wiki article about
Crostini](https://wiki.nixos.org/wiki/Installing_Nix_on_Crostinihttps://wiki.nixos.org/wiki/Installing_Nix_on_Crostini),
that details how to register applications in ChromeOS (so graphical
applications appear in menu) and also how to use
[nixGL](https://github.com/nix-community/nixGL) to make OpenGL applications
work.

Like I said at the start of the article, the device is woefully slow thanks to
its CPU and eMMC. It does mean that, for example, activating my Home-Manager
configuration takes a while (around 1 minute). But it is much faster than say,
[nix-on-droid](https://github.com/nix-community/nix-on-droid-app), that the
last time I tried in a much more powerful device ([Xiaomi Pad
5](https://www.gsmarena.com/xiaomi_pad_5-11042.php)), took 30 minutes until I
just decided to forget and uninstall it. Having a proper VM instead of
[proot](https://wiki.termux.com/wiki/PRoot) makes all the difference here.

I can even do some light programming here: using my
[ZSH](/2024-08-01/01-troubleshooting-zsh-lag-and-solutions-with-nix.md) and
neovim configuration (including LSP for coding) is reasonable fast. For
example, I did most of the code that [powers this
blog](/2024-07-29/01-quick-bits-why-you-should-automate-everything.md) using
this Chromebook.

Until Google decides to give us a proper VM in Android or release a hybrid
Chromebook device with better specs, this small Chromebook will probably stay
as my travel companion, and one of my favorite devices.
