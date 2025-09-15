# From Gaming Rig to Personal Computer: My Journey with NixOS and Jovian

I recently built a new gaming PC. The exact configurations don't matter but
if someone is interested you can take a look
[here](https://ie.pcpartpicker.com/list/JdTpv4). The most interesting tidbit
about it is maybe that I am using a AMD Radeon (an AMD Radeon
[9070](https://www.amd.com/en/products/graphics/desktops/radeon/9000-series/amd-radeon-rx-9070.html))
dedicated GPU for the first time in my life, not only because this is probably
the first time that I found an AMD GPU exciting but also because of the better
Linux support (that will come in handy later).

This PC was created as a gaming PC first, so it started its life with a Windows
11 installation. I always wanted it to eventually run Linux though, however I
didn't want it to be complicated to use because this was something that my wife
would also use. In a nutshell, I wanted to lower the barrier of entry for
someone that is unused to Linux desktops, so my window-tiling setup based on
[Sway](https://swaywm.org/) was a no-go.

Three years ago I also bought a Steam Deck. I realized that
[SteamOS](https://store.steampowered.com/steamos) would be a perfect fit for
this project: it is a streamlined experience based in [Steam's Big
Picture](https://store.steampowered.com/bigpicture) mode and a polished
[KDE](https://kde.org/) desktop. SteamOS in Desktop Mode has a single user and
there are no login prompts, making the computer truly personal (remember
Windows 95/98?) It has its limitations but I prefer to deal with its
limitations than my wife having a bad experience.

However, SteamOS is still not ready to be used as a general-purpose operational
system. Sure you can get the recovery image and try to install, but last time I
checked there was no support for my GPU since it is too new (even if AMD GPU
generally have good support on Linux). I was thinking of using something like
[Bazzite](https://bazzite.gg/) or [CachyOS](https://cachyos.org/) instead, but
what about NixOS?

Indeed there is a project to convert NixOS to have as much of a "SteamOS
experience" as possible, [Jovian-NixOS](https://jovian-experiments.github.io/).
As far as I know the project main objective is to be used in Steam Deck only
and other setups are unsupported, but the nice thing about NixOS is that we can
always make it work with other systems thanks to its declarative approach.

I will not go too deep about the setup, but I can recommend this [blog
post](https://ciarandegroot.com/archive/nixos-steam-box/) for a nice tutorial.
In my particular case I wanted support for
[FSR4](https://www.amd.com/en/products/graphics/technologies/fidelityfx/super-resolution.html)
and eventually I figured out that there is built-in support in
[`proton-cachyos`](https://github.com/CachyOS/proton-cachyos), by simply
setting the
[`PROTON_FSR4_UPGRADE=1`](https://github.com/CachyOS/proton-cachyos/blob/683ebf2585e6c43b373021d6586e7f56318b6c78/README.md?plain=1#L344)
environment variable before launching a game. Of course, there is no
`proton-cachyos` in nixpkgs, but
[`chaotic-nyx`](https://github.com/chaotic-cx/nyx) to the rescue:

```nix
{ flake, ... }:
{
  imports = [
    flake.inputs.chaotic-nyx.nixosModules.default
    flake.inputs.jovian-nixos.nixosModules.default
  ];

  {
    # Any recent kernel will do, but since we are already pulling chaotic-nyx
    # why not use CachyOS's kernel?
    boot.kernelPackages = pkgs.linuxPackages_cachyos;

    # Add the newest MESA drivers so we can get the latest performance
    # improvements
    # This will probably break NVIDIA Optimus, do not use if you have a
    # NVIDIA GPU
    chaotic.mesa-git.enable = true;

    jovian = {
      steam = {
        enable = true;
        autoStart = true;
        user = "deck";
        desktopSession = config.services.displayManager.defaultSession;
        # Add custom proton packages to Steam in Gamescope session
        environment = {
          STEAM_EXTRA_COMPAT_TOOLS_PATHS =
            lib.makeSearchPathOutput "steamcompattool" ""
              config.programs.steam.extraCompatPackages;
        };
      };
      hardware.has.amd.gpu = true;
    };

    # Add custom proton packages to Steam in Desktop session
    # They show in the compatibility tools for each game in Steam,
    # (see: https://steamdeckhq.com/tips-and-guides/the-proton-ge-steam-deck-guide/)
    # and setting it to "Proton-CachyOS" or "GE-Proton" and adding
    # `PROTON_FSR4_UPGRADE=1` to the command-line arguments for
    # that particular game I get FSR4 support
    programs.steam.extraCompatPackages = with pkgs; [
      proton-cachyos
      proton-ge-custom
    ];

    # Enable KDE as the desktop environment
    # If you enabled `jovian.steam.autoStart` even things like
    # "Switch to Desktop" works fine
    services = {
      desktopManager.plasma6.enable = true;
      displayManager.defaultSession = "plasma";
    };

    home-manager.users."deck" = {
      # Add "Return to Gaming Mode" desktop shortcut, like in Steam Deck
      home.file."Desktop/Return-to-Gaming-Mode.desktop".source =
        (pkgs.makeDesktopItem {
          desktopName = "Return to Gaming Mode";
          exec = "qdbus org.kde.Shutdown /Shutdown org.kde.Shutdown.logout";
          icon = "steam";
          name = "Return-to-Gaming-Mode";
          startupNotify = false;
          terminal = false;
          type = "Application";
        })
        + "/share/applications/Return-to-Gaming-Mode.desktop";
      # Boot in desktop mode instead of gaming mode by default
      # https://github.com/Jovian-Experiments/Jovian-NixOS/discussions/488
      xdg.stateFile."steamos-session-select" = {
        text = config.jovian.steam.desktopSession;
      };
    };
  };
}
```

By the way, NixOS has the `programs.steam.gamescopeSession` options that allow
you to get a separate [gamescope](https://github.com/ValveSoftware/gamescope)
session starting Steam in Big Picture mode, getting an almost SteamOS
experience without the third-party Jovian repository. However the integration
is partial and there is no support for things like pairing new Bluetooth
controllers in settings.

I am happy with my current setup. Combining the power of NixOS and the
streamlined experience of SteamOS makes for a great experience. I still have
all the power from NixOS to do things like connecting to my Tailscale account
for remote management, or having support for
[ROCm](https://www.amd.com/en/products/software/rocm.html) and
[distrobox](https://distrobox.it/) to play with AI (I probably should nixify
this too but this will be work for future me). Of course I could probably do
all the same with SteamOS or Bazzite, but having access to NixOS sure makes
things easier, especially since I can reuse my [Nix
configuration](https://github.com/thiagokokada/nix-configs/).
