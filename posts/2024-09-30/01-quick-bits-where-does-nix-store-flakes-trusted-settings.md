# Quick bits: where does Nix store Flake's trusted settings?

Have you ever run a `nix build` command and had this prompt?

```console
$ nix run .#darwinActivations/Sekai-MacBook-Pro
do you want to allow configuration setting 'extra-substituters' to be set to 'https://nix-community.cachix.org https://thiagokokada-nix-configs.cachix.org' (y/N)? y
do you want to permanently mark this value as trusted (y/N)? y
```

And realise that you did/didn't want to mark this value as trusted? But where
is this stored? Well, I had to look at the Nix source code to find the answer,
since I can't find this anywhere in the web or in ChatGPT (but I am sure that
now that I posted this it will eventually be searchable), and the answer can be
found
[here](https://github.com/NixOS/nix/blob/c116030605bf7fecd232d0ff3b6fe066f23e4620/src/libflake/flake/config.cc#L13-L16):

```c++
Path trustedListPath()
{
    return getDataDir() + "/trusted-settings.json";
}
```

Where is `getDataDir()` though? I found the answer
[here](https://github.com/NixOS/nix/blob/c116030605bf7fecd232d0ff3b6fe066f23e4620/src/libutil/users.cc#L52-L65):

```c++
Path getDataDir()
{
    auto dir = getEnv("NIX_DATA_HOME");
    if (dir) {
        return *dir;
    } else {
        auto xdgDir = getEnv("XDG_DATA_HOME");
        if (xdgDir) {
            return *xdgDir + "/nix";
        } else {
            return getHome() + "/.local/share/nix";
        }
    }
}
```

So we solved the mystery:

- If `NIX_DATA_HOME` is set, the file will be in
`$NIX_DATA_HOME/trusted-settings.json`
- If `XDG_DATA_HOME` is set, the file will be in
`$XDG_DATA_HOME/nix/trusted-settings.json`
- Otherwise Nix will fallback to `$HOME/.local/share/nix/trusted-settings.json`

By the way, if you don't know why you got this prompt, if `flake.nix` has a
`nixConfig` attribute inside `outputs` **and** this `nixConfig` is an unsafe
attribute (like `extra-substituters`) you will get this prompt, unless you set
[`accept-flake-config =
true`](https://nix.dev/manual/nix/2.23/command-ref/conf-file#conf-accept-flake-config)
in your Nix configuration (please **do not do this**, it is dangerous because
it may allow running possible unsafe options without asking you first).

You can inspect the JSON file or delete it and Nix will prompt you again the
next time you run a `nix` command. And yes, saving this preference is
controversial considering this is Nix, but having the power to set `nixConfig`
attributes is really powerful, and with great powers comes great
responsibilities.
