# Quick bits: NixOS's boot.kernelPackages is a attrset

I don't know if this is just something that got me by surprise or not, but
[`boot.kernelPackages`](https://github.com/NixOS/nixpkgs/blob/bb16119a4a7639ebbc91ad0f516b324b0f7c9b68/nixos/modules/system/boot/kernel.nix#L40-L71)
option does not receive a derivation like most other packages, but instead
receives a function and returns an attribute set with all packages. Here is the
documentation:

> This option allows you to override the Linux kernel used by NixOS. Since
> things like external kernel module packages are tied to the kernel you’re
> using, it also overrides those. This option is a function that takes Nixpkgs
> as an argument (as a convenience), and returns an attribute set containing at
> the very least an attribute kernel. Additional attributes may be needed
> depending on your configuration. For instance, if you use the NVIDIA X
> driver, then it also needs to contain an attribute `nvidia_x11`.

The kernel package itself is referenced by the
[`kernel`](https://github.com/NixOS/nixpkgs/blob/bb16119a4a7639ebbc91ad0f516b324b0f7c9b68/nixos/modules/system/boot/kernel.nix#L331-L332)
derivation inside this attribute set:

```console
nix-repl> nixosConfigurations.sankyuu-nixos.config.boot.kernelPackages.kernel
«derivation /nix/store/5zyjvf3qgfk52qmgxh36l4dkr9lf100x-linux-6.10.3.drv»
```

The reason for this is because it ensure that things like modules are built
with the same kernel version you are booting.

However one less obvious consequence about this is that if you want packages
that come from `linuxPackages`, say for example
[`cpupower`](https://github.com/NixOS/nixpkgs/blob/nixos-unstable/pkgs/os-specific/linux/cpupower/default.nix)
or
[`turbostat`](https://github.com/NixOS/nixpkgs/blob/nixos-unstable/pkgs/os-specific/linux/cpupower/default.nix),
it is better to do:

```nix
{ config, ... }:
{
  environment.systemPackages = [
    config.boot.kernelPackages.cpupower
    config.boot.kernelPackages.turbostat
  ];
}
```

Instead of:

```nix
{ pkgs, ... }:
{
  environment.systemPackages = with pkgs; [
    linuxPackages.cpupower
    linuxPackages.turbostat
  ];
}
```

Now, I used the later for a long time and never had issues. But technically
those packages depends in a specific kernel version for a reason, so getting
them from `config.boot.kernelPackages` reduces the chance of you having some
compatibility issue in the future.
