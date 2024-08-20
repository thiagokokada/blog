# Things I don't like in my Chromebook Duet 3

So this is kind of a continuation from my [previous
post](/2024-08-05/01-my-favorite-device-is-a-chromebook.md) talking why my
favorite device is a Chromebook. In this post I want to talk about what makes
me this device unhappy, and comment about things that if changed would make it
a much better device.

But before talking about the negative aspects, let me talk about a positive
aspect that I just briefly talked in the previous post: the screen. It is a
HiDPI screen (2000x1200 resolution in 10.95''), that is unexpected bright (400
nits according to the
[specs](https://www.lenovo.com/us/en/p/laptops/lenovo/lenovo-edu-chromebooks/ideapad-duet-3-chromebook-11-inch,-qlc/len101i0034)).
It is difficult to find laptops at the same price with a screen that good. At
10.95'' in its default resolution I find it too small (like 1250x750), but I
find the font size acceptable at 115% scale (1087x652). Yes, it result in a
small workspace, but this is not a big issue for what I do in this device. It
is also only 60Hz, but I thought I would miss high refresh rate more than I
actually miss in this device.

Update: I forgot to say one thing about the screen: it scratches really easy. I
got my screen scratched after the first day of usage, and considering the price
I don't think the screen has a hardened glass. I bought a cheap glass screen
protector and this did the trick though, even hiding the previous scratch, and
I have zero issues with the screen now.

Now the first aspect that I don't like: the speakers. They sound tiny and even
at maximum volume it is not really loud. The speakers is the only reason why I
still keep my [Xiaomi Pad 5](https://www.gsmarena.com/xiaomi_pad_5-11042.php),
because I like to watch animes/videos before sleep and having good speakers is
a must.

The keyboard has that issue that I mentioned in the previous post: sometimes
the key get stuck, and I get duplicated characters. But it also has some minor
issues that I didn't talked about: the first one is the UK layout that has some
extra keys that I have no use for, but this also makes the keys that I use
smaller. Very much a "me" problem, since if I had got a US version I wouldn't
have those issues, but an issue nonetheless that gets worse considering how
small the keyboard is. I am actually suprised how fast I can type considering
how many issues this keyboard has, so maybe this is a testament that this
keyboard is not actually that bad.

The other keyboard issue is a problem that affects all Chromebooks: its custom
layout. Google replaced a few keys like Fn keys with shortcuts and replaced the
Caps Lock with a
["Everything"](https://chromeunboxed.com/chromebook-launcher-now-everything-button)
key (that is similar to the Windows Key), while removing Windows Key from its
place. I actually have less issue with this than I initially though: I don't
care too much about Fn keys (except when using IntelliJ, but that is something
that I only use at `$CURRENT_JOB`), and ChromeOS is surprisingly powerful in
its customisation, allowing you to swap key functionality. I remap Everything
key with Esc, and Esc for the Everything key, and I can get productive in my
`neovim` setup.

And finally, let me talk more about the performance: yes, it is bad, but
bearable once you get used to. The issue is both the CPU and IO. While the CPU,
a [Snapdragon 7c Gen
2](https://www.qualcomm.com/products/mobile/snapdragon/laptops-and-tablets/snapdragon-mobile-compute-platforms/snapdragon-7c-gen-2-compute-platform)
is octa-core, it has only 2 high performance CPU cores vs 6 low performance
ones (2xARM Cortex A76 vs 6xARM Cortex A55). If it was something like 4x4, it
would be much better. The fact that the cores are old doesn't help either.

But the worst part is the IO. Not only it uses a eMMC module, it is slow:

[![CPDT Benchmark results from Chromebook Duet 3.](/2024-08-12/Screenshot_2024-08-12_20.50.42.png)](/2024-08-12/Screenshot_2024-08-12_20.50.42.png)

I don't know how much more expensive it would be to put a
[UFS](https://en.wikipedia.org/wiki/Universal_Flash_Storage) instead of eMMC in
this device, but this is probably the choice that would most increase
performance in this device, especially considering how aggressive Chromebooks
use (z)swap.

Update 2: I forgot to talk about the fact that the exterior of the device is
covered in cloth. I thought I would hate this at first, but nowadays I kind
like it. And it is also nice that it will never get scratched, I don't care too
much about the exterior of this device and it is the only device that I have at
home that doesn't have additional protection (except the screen protector
mentioned above).
