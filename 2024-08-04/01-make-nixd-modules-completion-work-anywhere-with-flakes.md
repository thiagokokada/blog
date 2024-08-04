# Make nixd module completion to work anywhere (with Flakes)

If you want the TL;DR, go to the bottom of the post (search for "final
result").

I recently switched from [nil](https://github.com/oxalica/nil) to
[nixd](https://github.com/nix-community/nixd) as my LSP of choice for Nix. I
was curious in `nixd` for a long time since the fact that it can eval Nix code
means it can offer much more powerful completion than the `nil`'s static
analysis, however it used to be difficult to setup. Nowadays it is much easier,
basically doing the right thing as long as you have `NIX_PATH` setup, and you
get both package and NixOS modules completion.

Getting Home-Manager modules though needs some setup. The recommended way to
setup accordingly to the [official
documentation](https://github.com/nix-community/nixd/blob/fe202307eaf7e89c4366ed927af761482a6065c8/nixd/docs/configuration.md)
is to use the following for Flake based configurations (using neovim
configuration here, but it should be easy to adapt to other editors):

```lua
{
  nixpkgs = {
    expr = "import <nixpkgs> { }",
  },
  options = {
    nixos = {
      expr = '(builtins.getFlake ("git+file://" + toString ./.)).nixosConfigurations.miku-nixos.options',
    },
    home_manager = {
      expr = '(builtins.getFlake ("git+file://" + toString ./.)).homeConfigurations.home-linux.options',
    },
  },
  -- ...
}
```

This works, but it should be pretty clear the limitations of using `./.`: this
will only work if you open your editor in your [Nix configuration
repository](https://github.com/thiagokokada/nix-configs). For any other
repository, `nixosConfigurations.miku-nixos` or `homeConfigurations.home-linux`
will not exist and the completion will not work.

It may look like this is easy to fix if you have `specialArgs` set to have your
Flakes inputs, but:

```nix
# By the way, ${self} does not exist in the Flake output by default, you need
# to explicit add `inherit self` to your outputs:
# https://discourse.nixos.org/t/who-is-self-in-flake-outputs/31859/4
nix-repl> (builtins.getFlake "git+file://${self}").nixosConfigurations.miku-linux.options
error:
       … while calling the 'getFlake' builtin
         at «string»:1:2:
            1| (builtins.getFlake "git+file://${self}")
             |  ^

       … while evaluating the argument passed to builtins.getFlake

       error: the string 'git+file:///nix/store/avr1lcmznj8ghynh5vj1kakgfdf0zrxx-source' is not allowed to refer to a store path (such as 'avr1lcmznj8ghynh5vj1kakgfdf0zrxx-source')
```

Well, it was worth a try. Another option would be to:

```nix
(builtins.getFlake "github:thiagokokada/nix-configs").nixosConfigurations.miku-linux.options
# Or even something like this
# BTW, using ${rev} means this wouldn't work in Flake repos, since it is not
# set in those cases
(builtins.getFlake "github:thiagokokada/nix-configs/${rev}").nixosConfigurations.miku-linux.options
```

But while it works, it is slow, because it needs network to evaluate (and it is
impure, since there is no `flake.lock`).

The default configuration for `nixd` makes NixOS completion work even outside
of my configuration repo, and it is fast. How? I decided to take a look at the
`nixd` source code and found
[this](https://github.com/nix-community/nixd/blob/d938026c55c7c36a6e79afd9627459160b4924ed/nixd/lib/Controller/LifeTime.cpp#L33C11-L35C76)
(formatted here for legibility):

```nix
(
  let
    pkgs = import <nixpkgs> { };
  in
  (pkgs.lib.evalModules {
    modules = (import <nixpkgs/nixos/modules/module-list.nix>) ++ [
      ({ ... }: { nixpkgs.hostPlatform = builtins.currentSystem; })
    ];
  })
).options
```

Interesting, so they're manually loading the modules using `evalModules`. As I
said above, it depends in `NIX_PATH` being correctly set. Can we fix this to
use our Flake inputs instead? After some tries in the Nix REPL, I got the
following:

```nix
(
  let
    pkgs = import "${inputs.nixpkgs}" { };
  in
  (pkgs.lib.evalModules {
    modules = (import "${inputs.nixpkgs}/nixos/modules/module-list.nix") ++ [
      ({ ... }: { nixpkgs.hostPlatform = builtins.currentSystem; })
    ];
  })
).options
```

So we can adapt this to the neovim configuration:

```lua
{
  options = {
    nixos = {
      expr = '(let pkgs = import "${inputs.nixpkgs}" { }; in (pkgs.lib.evalModules { modules =  (import "${inputs.nixpkgs}/nixos/modules/module-list.nix") ++ [ ({...}: { nixpkgs.hostPlatform = builtins.currentSystem;} ) ] ; })).options',
    },
  },
}
```

This was easy. But the main issue is Home-Manager. How can we fix it? I needed
to take a look at the Home-Manager [source
code](https://github.com/nix-community/home-manager/blob/master/docs/default.nix#L161-L169)
to find the answer:

```nix
(
  let
    pkgs = import "${inputs.nixpkgs}" { };
    lib = import "${inputs.home-manager}/modules/lib/stdlib-extended.nix" pkgs.lib;
  in
  (lib.evalModules {
    modules = (import "${inputs.home-manager}/modules/modules.nix") {
      inherit lib pkgs;
      check = false;
    };
  })
).options
```

The interesting part is: Home-Manager has its own extension of the module
system (including `evalModules`). This includes e.g.: extra types used in
Home-Manager only. Also, we need to disable `checks`, otherwise we will hit
some validations (e.g.: missing `stateVersion`). I am not sure if this causes
any issue for module completion yet, I may set it in the future.

And for the final result:

```lua
{
  nixpkgs = {
    expr = 'import "${flake.inputs.nixpkgs}" { }',
  },
  options = {
    nixos = {
      expr = '(let pkgs = import "${inputs.nixpkgs}" { }; in (pkgs.lib.evalModules { modules =  (import "${inputs.nixpkgs}/nixos/modules/module-list.nix") ++ [ ({...}: { nixpkgs.hostPlatform = builtins.currentSystem;} ) ] ; })).options',
    },
    home_manager = {
      expr = '(let pkgs = import "${inputs.nixpkgs}" { }; lib = import "${inputs.home-manager}/modules/lib/stdlib-extended.nix" pkgs.lib; in (lib.evalModules { modules =  (import "${inputs.home-manager}/modules/modules.nix") { inherit lib pkgs; check = false; }; })).options',
    },
  },
}
```

Yes, it is quite a mouthful, but it makes module completion work in any
repository, as long as you're using Flakes. And it is fast, since it doesn't
need any network access. Since we are already here, let's define `nixpkgs` to
not depend in the `NIX_PATH` being set too.
