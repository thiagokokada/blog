# An unordered list of hidden gems inside NixOS

After using [NixOS](https://nixos.org/) for the last 5+ years as my main OS, I
end up with a [configuration](https://github.com/thiagokokada/nix-configs/)
with many things that are interesting for one reason or another, but it is not
listed anywhere (well, except if you are the kind of person that reads `man 5
configuration.nix` or the release notes in every release).

So kind in the same spirit as my [list of things that I miss in
Go](/2024-08-17/01-an-unordered-list-of-things-i-miss-in-go.md), here is a list
of modules that I find neat in NixOS and are not default already. Again, the
list is unordered since this makes it easier to update in the future if I find
something else, but also I don't want to think too hard about an order here.

With all above, let's start.

##
[`networking.nftables`](https://github.com/NixOS/nixpkgs/blob/6afb255d976f85f3359e4929abd6f5149c323a02/nixos/modules/services/networking/nftables.nix)

[nftables](https://www.nftables.org/) is, accordingly to Wikipedia:

> nftables is a subsystem of the Linux kernel providing filtering and
> classification of network packets/datagrams/frames.

It is basically a replacement of the venerable
[iptables](https://en.wikipedia.org/wiki/Iptables), that still exist and is the
default program to configure the famous `networking.firewall`, the declarative
[Firewall](https://wiki.nixos.org/wiki/Firewall) that NixOS enable by default.

To enable, it is simple, just add to your configuration:

```nix
{
  networking.nftables.enable = true;
}
```

And thanks to the fact that NixOS's Firewall is declarative, everything should
still work as expect: any service that you set `openFirewall = true` will still
have its ports open, if you set `networking.firewall.allowPing = false` it will
still disable pings like before, etc.

If you look at the documentation of the above option, you will find the
following warning:

> Note that if you have Docker enabled you will not be able to use nftables
> without intervention. Docker uses iptables internally to setup NAT for
> containers. This module disables the ip_tables kernel module, however Docker
> automatically loads the module. Please see
> https://github.com/NixOS/nixpkgs/issues/24318#issuecomment-289216273 for
> more information.

I don't use Docker (switched to Podman instead for quite a long time), so I
don't know how bad the situation is. Also keep in mind that `nftables` does
offer `iptables-compat` for compatibility with old iptables scripts, so it is
most likely Docker doing something weird here.

Now, the actual advantage from the user here is not clear: the main advantage
from my point of view (and the reason I used to use in other distros like Arch)
is the improved syntax, however if you are using the declarative NixOS's
Firewall you are not interacting with either `iptables` or `nftables` directly
anyway. `nftables` is supposed to be more efficient, but not sure most users
will care about this.

However if you are the kind of person that needs custom rules, switching to
`nftables` does bring a few benefits, including
`networking.nftables.checkRuleset` (enabled by default), that checks if your
ruleset has syntax errors during build time. Really valuable to avoid issues
only after switch.

Anyway, this is one of those options that I think it should be the default for
a long time, since most of the new development in NixOS firewall seems to be
focusing `nftables` for a while.

## [`system.switch.enableNg`](https://github.com/NixOS/nixpkgs/blob/877d19523edcac81b167e8fd716ad2658da2adca/nixos/modules/system/activation/switchable-system.nix#L30-L38)

[This one](https://github.com/NixOS/nixpkgs/pull/308801) I just discovered
today, but it has been available for a while (~2 months if you're using
`nixos-unstable`). Finally someone is rewriting
[`switch-to-configuration.pl`](https://github.com/NixOS/nixpkgs/blob/b1eff03c35aa7c90ab3a4d9f6ef297dae5fba37b/nixos/modules/system/activation/switch-to-configuration.pl),
the Perl script that is called everytime you run `nixos-rebuild switch`.

Now, I am not one of those "rewrite in Rust" zealots, but in this case this is
definitely worth it: `switch-to-configuration.pl` is one of those pieces of
code in NixOS that most people avoid touching at the fear of breaking
something. There is a reason why
[`nixos-rebuild`](https://github.com/NixOS/nixpkgs/commit/eeb2588a59c938042b74183ce1da7052a6ef7e59)
is as convoluted as it is, because even if it is a messy shell script, most
people preferred to workaround issues from the `switch-to-configuration.pl`
inside it than trying to understand the mess that `switch-to-configuration.pl`
is.

Trying this one is easy:

```nix
{
  system.switch = {
    enable = false;
    enableNg = true;
  };
}
```

Yes, you need to explicit set `system.switch.enable = false`, since the default
is `true`.

By the way, what is the reason you would want to set `system.switch.enable =
false` before the `enableNg` appeared you ask? For systems that are immutable
and updated by e.g.: image upgrades instead of modifying root.

Enabling `switch-to-configuration-ng` right now is mostly for testing purposes,
but one of the advantages that I saw is that system switches are (slightly)
faster:

```
$ hyperfine "sudo nixos-rebuild switch" # switch-to-configuration.pl
Benchmark 1: sudo nixos-rebuild switch
  Time (mean ± σ):      3.576 s ±  0.035 s    [User: 0.004 s, System: 0.014 s]
  Range (min … max):    3.522 s …  3.645 s    10 runs

$ hyperfine "sudo nixos-rebuild switch" # switch-to-configuration-ng
Benchmark 1: sudo nixos-rebuild switch
  Time (mean ± σ):      3.394 s ±  0.080 s    [User: 0.004 s, System: 0.013 s]
  Range (min … max):    3.325 s …  3.608 s    10 runs
```

But yes, the difference is not enough to make a significant impact, and it is
not the objective anyway. The real reason for the rewrite is to make it easier
to colaborate. I hope one day we also have someone brave enough to rewrite the
`nixos-rebuild` script in something saner.

## [boot.initrd.systemd](https://github.com/NixOS/nixpkgs/blob/cce9aef6fd8f010d288d685b9d2a38f3b6ac47e9/nixos/modules/system/boot/systemd/initrd.nix)

A quick recap on how a modern Linux distro generally boots: the first thing
that the bootloader (say [GRUB](https://www.gnu.org/software/grub/) or
[systemd-boot](https://systemd.io/BOOT/)) loads is `initrd` (_initial
ramdisk_), a small image that runs from RAM and includes the Linux kernel and
some utilities that are responsible for setting up the main system. For
example, one of the responsabilities of the `initrd` is to mount the disks and
start init system (`systemd`).

It may surprising that this `initrd` image does **not** generally include
`systemd`. Traditionally `initrd` is composed by a bunch of shell scripts and a
minimal runtime (e.g.: [busybox](https://www.busybox.net/)), however `systemd`
can also do this job since a long time ago. It is just the paper of the distros
to integrate `systemd` inside the `initrd`.

This is what `boot.initrd.systemd` does: enable `systemd` inside the `initrd`.
It make a few subtle changes:

- If you are using [Full Disk Encryption via
LUKS](https://wiki.nixos.org/wiki/Full_Disk_Encryption), you will get a
different password prompt at login
- You will get `initrd` time information if using `systemd-analyze` to measure
boot time
  + You can get even more information (bootloader) if you also use
  `systemd-boot`
- You will also get `systemd` style status about services during `initrd` (not
only afterwards)

But I think the main reason is that since `systemd` is event-driven, it should
make boot more reliable, especially in challenging situations (like booting
from network). I can't say that I have any system like this to test if it is
actually more reliable or not, but I don't remember having any issues since I
set `boot.initrd.systemd.enable = true`, so there is that.

##
[services.pipewire](https://github.com/NixOS/nixpkgs/blob/b4a09f1f9d1599478afadffa782a02690550447c/pkgs/development/libraries/pipewire/default.nix)

If there is something in that list that has a good chance that you're using
already, it is this one, especially if you're using
[Wayland](https://wayland.freedesktop.org/). Still, I think it is interesting
to include in this list since [PipeWire](https://www.pipewire.org/) is great.

The experience with PipeWire until now for me was seamless: I never had any
issues with it, all my applications still work exactly as it always worked. I
also didn't had any issues with
[PulseAudio](https://www.freedesktop.org/wiki/Software/PulseAudio/) for a
while, but I still remember when I first tried PulseAudio during the 0.x in
Fedora and having tons of issues. So bonus points for PipeWire developers for
polishing the experience of enough that most people will feel no diffference.

To enable PipeWire, I would recommend:

```nix
{
  services.pipewire = {
    enable = true;
    alsa.enable = true;
    pulse.enable = true;
    # jack.enable = true;
  };
  security.rtkit.enable = true;
}
```

This enables both ALSA and PulseAudio emulation support in PipeWire for maximum
compatibility with desktop applications (you can also enable
[`jack`](https://jackaudio.org/) if you use professional audio applications).
It also enables [`rtkit`](https://github.com/heftig/rtkit), allowing PipeWire
to get (soft) realtime, helping avoiding cracks during high CPU load.

I also recommend taking a look at the [Wiki
article](https://wiki.nixos.org/wiki/PipeWire), that has multiple interesting
configurations that can be added for low-latency setups or improved codecs for
Bluetooth devices.

## [`networking.networkmanager.wifi.backend = "iwd"`](https://github.com/NixOS/nixpkgs/blob/c9ec8289781a3c4ac4dd5c42c8d50dd65360e79c/nixos/modules/services/networking/networkmanager.nix#L264-L271)

There is a good change that you're using
[`NetworkManager`](https://www.networkmanager.dev/) to manage network,
especially for Wi-Fi. And if that is the case, I can't recommend enough
changing the backend from the default `wpa_supplicant` to
[`iwd`](https://iwd.wiki.kernel.org/).

If you think that your Wi-Fi takes a long time to connect/re-connect, it may be
because `wpa_supplicant`. `iwd` seems much more optimised in this regard, and
since switching to it I never felt that my Wi-Fi was worse than other OSes (and
generally slightly better than Windows, but keep in mind that this is a
non-scientific comparison).

Not saying that I never had Wi-Fi issues since switching to `iwd`, however
switching back to `wpa_supplicant` in those cases never fixed the issue (it was
the same or worse), so I assume either bad hardware or drivers in those cases.
