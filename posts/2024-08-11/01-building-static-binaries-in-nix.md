# Building static binaries in Nix

I recently had laser eye surgery for my near-sightedness, and while evaluating
if I could have the surgery I discovered that I was suffering from dry eyes.
Thanks to this, my ophthalmologist recommended that every 20 minutes while
using screens, I look somewhere 20 feet away for 20 seconds, a technique known
as [20-20-20 rule](https://www.healthline.com/health/eye-health/20-20-20-rule).

I had issues following this rule because I never remembered to do the pauses. I
initially tried to setup an alarm, but this became annoying, so I decided to
try to find a program. I found
[this](https://tonyh4156.github.io/20-20-20.github.io/) one for macOS that
worked fine, but it bothered me that it was compiled for x86_64 while I was
using a M1 MacBook Pro at the time, and also I needed something that worked in
Linux.

Since I never found a good multi-platform alternative, I decided to write my
own. This became
[twenty-twenty-twenty](https://github.com/thiagokokada/twenty-twenty-twenty/),
the first [Go](/posts/2024-07-29/02-go-a-reasonable-good-language.md) program
that I ever wrote. I wrote it in Go because I wanted to learn the language, but
also because Go made it easy to build static binaries. And the first version I
could build static binaries without issues because I was using
[beeep](https://github.com/gen2brain/beeep), that uses pure Go code in all
supported platforms. However, it also meant that the notifications in macOS
looked ugly, since it used
[osascript](https://github.com/gen2brain/beeep/blob/master/beep_darwin.go#L20).

I wanted better integration with macOS, so this meant switching libraries.
After searching for a while, the
[notify](https://pkg.go.dev/gioui.org/x/notify) library from
[GioUI](https://gioui.org/) is the one that seemed to work better. It
implements notification in macOS using its native framework, so it works much
better, but sadly it meant losing static binaries because it depends in CGO.

Not a big loss initially, because I am only depending in Foundation inside
macOS (that should always be available), and in Linux I could still statically
compile. However I eventually added more features like sound
(via [beep](https://github.com/gopxl/beep)) and tray icon (via
[systray](https://github.com/fyne-io/systray)), that meant I needed CGO in both
macOS and Linux.

Losing static binaries in Linux is a much bigger deal, since Linux is a moving
target. The general recommendation for building CGO binaries statically is
using
[musl](https://eli.thegreenplace.net/2024/building-static-binaries-with-go-on-linux/),
but this also means building all dependencies that we need using musl (e.g.:
[`ALSA`](https://github.com/ebitengine/oto?tab=readme-ov-file#linux) for
[beep/oto]). This generally means pain, but Nix makes it easy.

Let's start by creating a [Nix
file](https://github.com/thiagokokada/twenty-twenty-twenty/blob/main/twenty-twenty-twenty.nix)
that builds our Go module (simplified below for brevity):

```nix
{ lib
, stdenv
, alsa-lib
, buildGoModule
, pkg-config
, withStatic ? false
}:

buildGoModule {
  pname = "twenty-twenty-twenty";
  version = "1.0.0";
  src = lib.cleanSource ./.;
  vendorHash = "sha256-NzDhpJRogIfL2IYoqAUHoPh/ZdNnvnhEQ+kn8A+ZyBw=";

  CGO_ENABLED = "1";

  nativeBuildInputs = lib.optionals (stdenv.hostPlatform.isLinux) [
    pkg-config
  ];

  buildInputs = lib.optionals (stdenv.hostPlatform.isLinux) [
    alsa-lib
  ];

  ldflags = [ "-X=main.version=${version}" "-s" "-w" ]
    ++ lib.optionals withStatic [ "-linkmode external" ''-extldflags "-static"'' ];

  meta = with lib; {
    description = "Alerts every 20 minutes to look something at 20 feet away for 20 seconds";
    homepage = "https://github.com/thiagokokada/twenty-twenty-twenty";
    license = licenses.mit;
    mainProgram = "twenty-twenty-twenty";
  };
}
```

And we can build it with the following `flake.nix`:

```nix
{
  description = "twenty-twenty-twenty";

  inputs = {
    nixpkgs.url = "nixpkgs/nixpkgs-unstable";
    flake-compat.url = "github:edolstra/flake-compat";
  };

  outputs = { self, nixpkgs, ... }:
    let
      supportedSystems = [ "x86_64-linux" "aarch64-linux" ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      packages = forAllSystems (system:
        let pkgs = nixpkgsFor.${system};
        in {
          default = self.packages.${system}.twenty-twenty-twenty;
          twenty-twenty-twenty = pkgs.callPackage ./twenty-twenty-twenty.nix { };
          twenty-twenty-twenty-static = pkgs.pkgsStatic.callPackage ./twenty-twenty-twenty.nix {
            withStatic = true;
          };
        });
    };
}
```

I think this shows how powerful Nix is: the only difference between the normal
build and a static build the usage of `pkgs.pkgsStatic` instead of `pkgs`. This
automatically builds all packages statically with `musl`. Also we need pass
some [extra
flags](https://honnef.co/articles/statically-compiled-go-programs-always-even-with-cgo-using-musl/)
to the Go compiler (i.e.: `-linkmode external -extldflags "-static"`), but this
is a requirement from Go.

So, does it work? Let's test:

```console
$ nix build .#twenty-twenty-twenty-static

$ file result/bin/twenty-twenty-twenty
result/bin/twenty-twenty-twenty: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, stripped

$ ./result/bin/twenty-twenty-twenty
ALSA lib pcm.c:2712:(snd_pcm_open_conf) Either /nix/store/dhn51w2km4fyf9ivi00rz03qs8q4mpng-pipewire-1.2.1/lib/alsa-lib/libasound_module_pcm_pipewire.so cannot be opened or _snd_pcm_pipewire_open was not defined inside
ALSA lib pcm.c:2712:(snd_pcm_open_conf) Either /nix/store/ly9d7llymzjyf6gi1455qzayqipk2kab-pipewire-1.2.1/lib/alsa-lib/libasound_module_pcm_pipewire.so cannot be opened or _snd_pcm_pipewire_open was not defined inside
ALSA lib pcm.c:2712:(snd_pcm_open_conf) Either /nix/store/dhn51w2km4fyf9ivi00rz03qs8q4mpng-pipewire-1.2.1/lib/alsa-lib/libasound_module_pcm_pipewire.so cannot be opened or _snd_pcm_pipewire_open was not defined inside
ALSA lib pcm.c:2712:(snd_pcm_open_conf) Either /nix/store/ly9d7llymzjyf6gi1455qzayqipk2kab-pipewire-1.2.1/lib/alsa-lib/libasound_module_pcm_pipewire.so cannot be opened or _snd_pcm_pipewire_open was not defined inside
ALSA lib pcm.c:2712:(snd_pcm_open_conf) Either /nix/store/dhn51w2km4fyf9ivi00rz03qs8q4mpng-pipewire-1.2.1/lib/alsa-lib/libasound_module_pcm_pipewire.so cannot be opened or _snd_pcm_pipewire_open was not defined inside
ALSA lib pcm.c:2712:(snd_pcm_open_conf) Either /nix/store/ly9d7llymzjyf6gi1455qzayqipk2kab-pipewire-1.2.1/lib/alsa-lib/libasound_module_pcm_pipewire.so cannot be opened or _snd_pcm_pipewire_open was not defined inside
2024-08-11T19:26:33+01:00 INFO Running twenty-twenty-twenty every 20.0 minute(s), with 20 second(s) duration and sound set to true
```

There are some warns and sadly the sound doesn't work. I think the issue is
related because of my usage of PipeWire and the binary may work in a pure ALSA
system, but I don't have access to one. Maybe adding `pipewire` to
`buildInputs` would fix this issue, but I can't get `pipewire` to be compiled
statically (because of its dependencies). I think this is a good show how easy
it is to statically compilation is in Nix, but also how complex static binaries
are to get correctly.

Bonus points for
[cross-compilation](https://nix.dev/tutorials/cross-compilation.html). We can
easily cross-compile by using `pkgsCross`:

```nix
{
  # ...
  outputs = { self, nixpkgs, ... }:
    let
      # ...
    in
    {
      packages = forAllSystems (system:
        let pkgs = nixpkgsFor.${system};
        in {
          twenty-twenty-twenty-static-aarch64 = pkgs.pkgsCross.aarch64-multiplatform.pkgsStatic.callPackage ./twenty-twenty-twenty.nix {
            withStatic = true;
          };
        });
    };
}
```

The idea of `pkgsCross` is to select a target platform (e.g.:
`aarch64-multiplatform`) and use it as any other `pkgs`. We can even chain
`pkgsStatic` to statically cross compile binaries:

```console
$ nix build .#twenty-twenty-twenty-static-aarch64

$ file result/bin/twenty-twenty-twenty
result/bin/twenty-twenty-twenty: ELF 64-bit LSB executable, ARM aarch64, version 1 (SYSV), statically linked, stripped
```
