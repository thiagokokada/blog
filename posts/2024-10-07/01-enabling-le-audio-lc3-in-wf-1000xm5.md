# Enabling LE Audio/LC3 in WF-1000XM5

One of things that I hate the most about the fact that we are all using
wireless earbuds instead of wired earphones is the latency: it is bad, getting
up to seconds(!) depending on your particular combination of OS/earbuds/device.

There is a solution though: Bluetooth LE Audio, that is supposed to fix
multiple issues with the original design for Bluetooth Classic Audio, including
a much lower latency, improved efficiency (e.g.: less battery power) and even
multiple streams of audio. LE Audio also includes a new default codec for
improved audio quality, [LC3](https://en.wikipedia.org/wiki/LC3_(codec)), that
replaces the venerable [SBC](https://en.wikipedia.org/wiki/SBC_(codec)) codec
for audio.

However, the standard is a mess right now: a few wireless headphones already
support it, but they're generally disabled by default and it is pretty messy to
enable. And even after enabling it, getting it to work can be a pain.

I have pretty much the best setup to use LE Audio right now: a recently
released Pixel 9 Pro with Sony's
[WF-1000XM5](https://www.sony.ie/headphones/products/wf-1000xm5) earbuds, and
after lots of tries I got it to work. You can see below the versions of
everything I am using:

- Android: 14
- [Sound
  Connect](https://play.google.com/store/apps/details?id=com.sony.songpal.mdr):
  11.0.1
- WM-1000XM5: 4.0.2

The first thing you need to do is enable in "Sound Connect" app "LE Audio
Priority" in "Device Settings -> System":

[![LE Audio option inside Sound
Connect](/posts/2024-10-07/photo_4909454744305642922_y.jpg)](/posts/2024-10-07/photo_4909454744305642922_y.jpg)

After this, you will need to pair your headset with the device again. You can
do this as same as always: press and hold the button in case for a few seconds
until a blue light starts to blink. However, this is where things starts to get
janky: I couldn't get the headset to pair with Android again.

A few of the things that I needed to do (in no specific order):

- Remove the previous paired headset
- Restart the Android
- Clean-up "Sound Connect" storage (Long press the app icon -> "App info" ->
  "Storage and Cache" -> "Clear storage")

If you can get the headset to connect, go to the "Bluetooth" settings in
Android, click in the gear icon for the headset and enable "LE Audio" option:

[![LE Audio option Bluetooth
Settings](/posts/2024-10-07/photo_4909454744305642937_y.jpg)](/posts/2024-10-07/photo_4909454744305642937_y)

If you can't, you may want to [restore the headset to factory
settings](https://helpguide.sony.net/mdr/2963/v1/en/contents/TP1000783925.html)
and try again from the start (that means pairing your device with "Sound
Connect" again, and you may want to try to clear the storage before doing so).

Yes, the process is extremely janky, but I think this is why both "Sound
Connect" and Android marks this feature as beta/experimental. And I still need
to test the latency, but from my initial testing there are some glitches when
the audio is only used for a short period of time (e.g.: Duolingo only enables
the audio when the character is speaking). So I only recommend this if you want
to test how LE Audio will behave, since it is clear that this needs more
polish.
