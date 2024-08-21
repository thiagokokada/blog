# First impressions: FPGBC

Here is something for nostalgia: I just put together a [Game Boy
Color](https://en.wikipedia.org/wiki/Game_Boy_Color) made of completely new
parts for a friend: here is the
[FPGBC](https://funnyplaying.com/products/fpgbc-kit).

The _FP_ part of the name comes from
[FPGA](https://en.wikipedia.org/wiki/Field-programmable_gate_array), because
instead of software emulation this device use FPGA to reproduce the device.
While I am not convinced that FPGA is necessary more accurate than a good
software emulator, one advantage of FPGA is the (possible) lower input latency
thanks to the avoidance of complexity to handle the user input (e.g.: the
Operational System). A quick playthrough against [Motocross
Maniacs](https://en.wikipedia.org/wiki/Motocross_Maniacs) seems to be fine, but
I can't see much difference from my [Miyoo
Mini+](https://retrogamecorps.com/2022/05/15/miyoo-mini-v2-guide/) (I will do
more comparisons between the two devices later), that is a software emulation
device.

But I think focusing in accuracy is wrong, the main reason of getting a device
like this one is for nostalgia, and this definitely hit the mark. The quality
of the case is as good as I remember the original, and most of the details are
replicate perfectly, including reproduction stickers in the back of the device.
The only differences that I can find is the usage of USB-C port for charging in
place of the barrel jack power adapter (thanks!), and the fact that the screen
bezels are smaller compared to the original (because the screen is bigger) and
doesn't include the Game Boy Color logo (that is fine in my opinion, since it
would look weird in the fine bezels). It even has a supposedly working [Link
Cable](https://en.wikipedia.org/wiki/Game_Link_Cable) (I don't have another
Game Boy to test). Sadly it is missing the infrared sensor, but the usage of
that was pretty limited anyway.

[![FPGBC running Tetris.](/2024-07-30/PXL_20240729_175245569.jpg)](/2024-07-30/PXL_20240729_175245569.jpg)

[![Back of FPGBC. It includes even reproduction stickers of the original.](/2024-07-30/PXL_20240729_175131157.jpg)](/2024-07-30/PXL_20240729_175131157.jpg)

So how well does it work? I can't say for sure. I don't have any original games
with me, so I am relying in backups and a
[flashcard](https://en.wikipedia.org/wiki/Flashcard) for now. Many games that I
tested works fine, a few of them have graphical issues that can be fixed in the
menu (more about it later), and some of them doesn't boot. But I don't know if
the issue with the games not booting are because of the roms, the flashcard or
the device itself.

By the way, the flashcard I am using is a cheap knockoff of an [Everdrive
GB](https://gbatemp.net/review/everdrive-gb.141/). This FPGBC came with
firmware v1.09, while there is an update available for v1.10 in the
[website](https://funnyplaying.com/products/fpgbc-kit). I had an weird issue in
the new firmware where no games would boot with this knockoff Everdrive so I
had to go back to v1.09, but again, I am not sure if the issue was fact that I
am using a knockoff device or this would happen with an original Everdrive GB.
If you are going to buy a proper Everdrive, you probably wouldn't get a
Everdrive GB anyway since it is discontinued, and it seems the [newer
versions](https://www.reddit.com/r/Gameboy/comments/1atwjh3/fpgbc_everdrive_compatibility/)
have better compatibility with FPGBC.

Sadly that the update didn't work, since there is this
[repository](https://github.com/makhowastaken/GWGBC_FW) that patches the
firmware to boot the original logo instead of the ugly FPGBC one. And yes, for
some reason the v1.09 firmware from this repository still doesn't work with my
knockoff Everdrive.

By the way, it seems the device is not easy to brick: I borked the firmware
update process once while trying to downgrade back to v1.09, resulting in a
black screen when I turned on the console. But just connecting the device to
the computer and powering on, I could flash the firmware again and the device
came back to life.

About the features of the device: if you press the volume button (yes, you can
press it now), it opens the following menu:

[![FPGBC menu.](/2024-07-30/PXL_20240729_210604830.jpg)](/2024-07-30/PXL_20240729_210604830.jpg)

The first 2 options are the LCD backlight (`BKLT`) and volume (`VOL`). I didn't talk about
those, but the LCD screen seems to be IPS, and the quality is really good, and
also looks bright enough to play even under bad lightining conditions. And the
speaker has good quality, the sound is better than I remember, but sadly the
maximum volume is kind low. Still should be enough for playing in a quiet room.

`DISPMOD` is probably the most controversial option: it allow you to set which
scale you want. Anything with `EMU` at the end means emulating the original
colors, and as far I remember it gets really close. You can also chose betwen
`X4`, `X4P` and `FUL`, the last one is the one shown in the photos where the
image fills the whole screen at the cost of non-integer scaling. `X4` is
integer scaling, however the image doesn't fill the whole screen. The `X4P`
also includes a pixel effect that makes the image closer than the original
screen. It actually looks good, but the fact that I chose a white border for
this FPGBC makes the border really distracting. Maybe the black one is a better
choice if you want integer scale.

`CORE` is simple: you can choose between `GB` (Game Boy) or `GBC` (Game Boy
Color). For those who don't know, you can run Game Boy games in Game Boy Color
and they will be automatically colorised. Some people don't like this and
prefer the colors of `GB`, so you have this option. The `GB_PALETTE` allows you
to chose the color in GB mode, for example, the green-ish colors from the
original Game Boy or the blue-ish colors from [Game Boy
Light](https://nintendo.fandom.com/wiki/Game_Boy_Light). And yes, you can
choose the color palette for Game Boy games running in `GBC` mode by pressing a
[button combination](https://gbstudiocentral.com/tips/game-boy-color-modes/) at
the boot screen, but it seems not working in my unit and again, not sure if the
fault is my knockoff Everdrive.

`FRAME_MIX` basically is an option that makes some effects, like transparency
in [Wave Race](https://en.wikipedia.org/wiki/Wave_Race), to work at the cost of
introducing blurriness. The reason for this is that those effects depends in
the fact that the Game Boy screen was slow refresh, so you could rely on it by
rapidly changing pixels to create some interesting effects, but sadly those
effects doesn't work well in modern displays.

`GB_CLRFIX` is the option I mentioned before, where some Game Boy games just
get completely wrong colors for some reason, e.g.: [The Addams
Family](https://en.wikipedia.org/wiki/The_Addams_Family_(video_game)). Turning
on fixes those games, but I am not sure if this option breaks other games.

Finally, `SPD` allows you to increase or decrease the CPU clock, slowing or
speeding up the games (including the sound). The result can be hilarious, so I
think this is a nice addition to the features. Sadly you can't know what the
default speed is, so you need to rely on sound to adjust back to the default.

So in the end, can I recommend a FPGBC? I am not sure. If you want a device to
play games, I still think something like a Miyoo Mini+ is a better choice. Not
only you will have access to more games from different platforms, you also
don't need to rely on flashcards or cartridges. Also it has way more features
than FPGBC, like wireless multiplayer,
[RetroArchivements](https://retroachievements.org/) and save states.

But the actual reason to get a FPGBC is nostalgia, and for that I think the
FPGBC is difficult to beat. The price of the [kit to
assemble](https://funnyplaying.com/products/fpgbc-kit) ($69.90) is cheaper than
most Game Boy's in good condition you can find in eBay, and you get for that
price a rechargable battery, an amazing quality screen, the PCB and the
speaker. You need to buy separately the case and the buttons, but in total you
will still end up paying less, and allows you to fully customise your build.
And the result device is not only in mint condition, it is really convenient
too: recharging batteries (via USB-C even) is much more convenient than buying
AA batteries, and the screen not only is better but it even has backlight. You
can also buy a fully built console for
[$99.00](https://funnyplaying.com/products/fpgbc-console), but you have less
options of customisation.

This is the classic case of do what I say, don't do what I do. This FPGBC is a
gift, and I will buy another one soon. Can't wait to play [Pokémon
Gold](https://en.wikipedia.org/wiki/Pok%C3%A9mon_Gold_and_Silver) in (almost)
original hardware again.

[![The kit before assemble.](/2024-07-30/PXL_20240729_123847458.jpg)](/2024-07-30/PXL_20240729_123847458.jpg)
