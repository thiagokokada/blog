# My favorite device is a Chromebook, without ChromeOS

I wrote in 2024 a blog post about a how cheap and low powered
[Chromebook](/posts/2024-08-05/01-my-favorite-device-is-a-chromebook.md)
was one of my favorite devices to use day to day. At the time I was
praising how well ChromeOS worked, especially thanks
[Crostini](https://chromeos.dev/en/linux), a full Linux VM that you can
run inside ChromeOS.

Two years later, things changed. I still have the same Chromebook Lenovo Duet
3, and up until now I would sporadically use it. However the experience is
getting worse each day: I need to constantly [resize the VM disk a few
MBs](https://old.reddit.com/r/Crostini/comments/1rl2m2k/if_you_are_not_able_to_start_crostini_after_the/)
because otherwise Crostini would fail to start; I got corrupted characters in
the [Terminal app](https://issues.chromium.org/issues/502262228), because of a
recent introduced bug, and thanks to that I needed to disable the 2D Canvas
Acceleration in `chrome:flags`.

None of those issues are deal breaker since they both have workarounds, but
they're all recent issues. I kind think that this is the result of [Android and
ChromeOS getting
merged](https://www.androidauthority.com/google-combine-chrome-os-android-3577035/).
Google seems to be focusing more in getting Android to work well as a desktop
OS (and this shows, I have a Pixel 9 Pro and its [Linux
Terminal
VM](https://www.zdnet.com/article/how-to-use-the-new-linux-terminal-on-android/)
is getting better) than to improve ChromeOS. That makes sense, in some ways
Android is a much more capable OS than ChromeOS, but on the other hand it is
sad to see ChromeOS to wither in real time.

Since I was kind frustrated with ChromeOS, I decided to take a look at
something that [I knew supported my Lenovo Duet
3](https://wiki.postmarketos.org/wiki/Lenovo_IdeaPad_Duet_3_(google-wormdingler))
for some time: [postmarketOS](https://postmarketos.org/). For those who don't
know, postmarketOS is an [Alpine Linux](https://www.alpinelinux.org/)
based-distro focused in replacing the original OS from old phones (generally
running Android) with a "true" Linux distro. They also seem to support some
Chromebooks because of their unique architecture and, luckily, they support my
device under the [google-trogdor
platform](https://wiki.postmarketos.org/wiki/Google_Trogdor_Chromebook_(google-trogdor)).

I will not go too deep in the install instructions, I recommend following the
[official
instructions](https://wiki.postmarketos.org/wiki/Category:ChromeOS#Installation)
plus this [blog post for Duet
3 itself](https://blog.cyril.by/en/free-speech/tutorial-install-postmarketos-on-ideapad-duet-3-chromebook).
I did go a step further and replaced the `pmbootstrap install` command with:

```console
$ pmbootstrap install --disk=/dev/mmcblk1 --fde --filesystem=btrfs
```

The `--fde` flag enables Full Disk Encryption that I think it is a must have
for a portable device. Keep in mind that this will ask the password at every
boot and the password screen will always show in portrait instead of landscape
orientation, one of the quirks of using a LCD panel for tablets instead of one
for laptops.

The `--filesystem=btrfs` uses the (in)famous
[btrfs](https://en.wikipedia.org/wiki/Btrfs) filesystem instead of the default
ext4. I don't have strong opinions either in favor or against btrfs, but since
this is a more traditional distro (e.g., it is not an immutable distro like
NixOS), I think that having support for snapshots is a must. This
`--filesystem=btrfs` sets a pretty good subvolume layout too:

```console
# sudo btrfs subvolume list /
ID 256 gen 321 top level 5 path @
ID 257 gen 341 top level 5 path @home
ID 258 gen 298 top level 5 path @root
ID 259 gen 111 top level 5 path @snapshots
ID 260 gen 13 top level 5 path @srv
ID 261 gen 9 top level 5 path @tmp
ID 262 gen 341 top level 5 path @var
ID 263 gen 13 top level 262 path @var/lib/portables
ID 264 gen 167 top level 262 path @var/lib/machines

# mount | grep btrfs
/dev/mapper/root on / type btrfs (rw,relatime,compress=zstd:2,ssd,space_cache=v2,subvolid=256,subvol=/@)
/dev/mapper/root on /.snapshots type btrfs (rw,relatime,compress=zstd:2,ssd,space_cache=v2,subvolid=259,subvol=/@snapshots)
/dev/mapper/root on /root type btrfs (rw,relatime,compress=zstd:2,ssd,space_cache=v2,subvolid=258,subvol=/@root)
/dev/mapper/root on /srv type btrfs (rw,relatime,compress=zstd:2,ssd,space_cache=v2,subvolid=260,subvol=/@srv)
/dev/mapper/root on /var type btrfs (rw,relatime,compress=zstd:2,ssd,space_cache=v2,subvolid=262,subvol=/@var)
/dev/mapper/root on /home type btrfs (rw,relatime,compress=zstd:2,ssd,space_cache=v2,subvolid=257,subvol=/@home)
```

To be clear, this `@snapshots` subvol is not used by anything by default. I am
using it with [btrbk](https://github.com/digint/btrbk) to get daily backups
though:

```console
$ cat /etc/btrbk/btrbk.conf
compat busybox

timestamp_format long

snapshot_preserve_min 2d
snapshot_preserve 7d

volume /
  snapshot_dir /.snapshots

  subvolume .

$ systemctl cat btrbk.timer
# /etc/systemd/system/btrbk.timer
[Unit]
Description=Run btrbk daily

[Timer]
OnCalendar=daily
Persistent=true
RandomizedDelaySec=30min

[Install]
WantedBy=timers.target

$ systemctl cat btrbk.service
# /etc/systemd/system/btrbk.service
[Unit]
Description=btrbk Btrfs snapshots
Documentation=man:btrbk(1) man:btrbk.conf(5)
ConditionPathExists=/etc/btrbk/btrbk.conf

[Service]
Type=oneshot
ExecStart=/usr/bin/btrbk run
```

If you're coming from Alpine the commands above seems suspicious. Why systemd?
Well, [postmarketOS](https://postmarketos.org/blog/2024/03/05/adding-systemd/)
supports systemd, mostly because of KDE and GNOME, and I since I am using the
later and prefer systemd anyway, this was more of a blessing than a curse to
me. Having systemd installed by default also made installing Nix a breeze in
postmarketOS, I just needed to do:

```console
# doas-sudo-shim is missing some flags that are needed for the nix-installer,
# so we need to replace it with proper sudo
$ doas apk add curl sudo !doas !doas-sudo-shim
# The new experimental Nix installer: https://github.com/NixOS/nix-installer/
$ curl -sSfL https://artifacts.nixos.org/nix-installer | sh -s -- install --enable-flakes
```

And yes, I decided to use GNOME because I think this form factor (a hybrid
laptop/tablet) works perfectly with it. I always thought that GNOME was a well
polished interface but it is kind limiting for a desktop environment (this is
why I particularly prefer
[KDE](/posts/2025-09-17/01-kde-is-now-my-favorite-desktop.md) for desktops).
But this device is not a desktop and GNOME fits like a glove.

Do I recommend postmarketOS for this device? It depends on your needs. I
recommend reading the
[Wiki](https://wiki.postmarketos.org/wiki/Lenovo_IdeaPad_Duet_3_(google-wormdingler))
entry for the device first. For example, the pen is not working for me (but it
seems there is a workaround), neither docking (if I undock and dock the device
again the device stays in portrait). I don't use the pen so the first issue
doesn't matter to me, and for the second issue I disabled the rotation sensor
like the Wiki says.

Also the external display support is wonky enough that I stopped trying to use
it for now (the device has limited bandwidth anyway so I can't drive my desktop
monitor at full resolution). Another issue is that this device is working well
enough in `v25.12` release, but I tried `edge` once (the rolling release
channel) and my touchpad started to work in absolute instead of relative mode.
Sadly because I wanted to use `edge` to get GNOME 50, but I hope those issues
get ironed out until the next stable release.

But on the other hand, this is an actual Linux distro. Given the limitations
and issues, it is working really well. Before trying postmarketOS I thought the
only reason ChromeOS worked as well in this device as it did is because Google
invested a lot to optimise the OS and browser to run well in low end devices.
And while this is probably true, GNOME is also running really well here, and
GNOME is considered one of the heaviest desktop environments in Linux. I think
the fact that the base system is lean and the usage of `musl` instead of
`glibc` helps here, but I don't have any scientific measurement of this.

Also I am really happy to use Firefox as my browser instead of Chrome, since I
used Firefox everywhere else except in ChromeOS. I was worried that Firefox
would run too slow here since I tried it once via Crostini, but running it
natively in postmarketOS feels relatively snappy (considering the hardware
constraints).

While Crostini was really good, there was still some overhead of running a
Linux VM in a already constrained device. For example, cloning
[nixpkgs](https://github.com/NixOS/nixpkgs/) in this device took "only" ~15
minutes, while it took more than half an hour to clone in ChromeOS.

So for me, postmarketOS is a revelation. It kind remembers me at the time when
I installed Arch Linux in my Netbook and making it usable again after fighting
an old Windows installation. It is making this device usable again for my
needs, and I can say that the Duet 3 is back to becoming my favorite device
again, just this time without ChromeOS.
