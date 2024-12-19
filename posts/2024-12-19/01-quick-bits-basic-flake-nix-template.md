# Quick bits: basic flake.nix template

Sometimes I want a really basic `flake.nix` that has no dependencies except for
`nixpkgs` itself, e.g.: I want to avoid
[flake-utils](https://github.com/numtide/flake-utils) or any other dependency.
So, here you go:

```nix
{
  description = "Description";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs =
    { self, nixpkgs, ... }:
    let
      supportedSystems = [
        "aarch64-linux"
        "x86_64-linux"
        "aarch64-darwin"
        "x86_64-linux"
      ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      packages = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
        in
        rec {
          default = pkgs.hello;
        }
      );
    };
}
```

Not sure where I grabbed the definition for `forAllSystems` and `nixpkgsFor`. I
have the impression it was in a [Julia Evans blog post](https://jvns.ca/), but
I can't find it.

Anyway, it is here for me to remember, and it may help someone else.
