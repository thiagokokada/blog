# Quick bits: realise Nix symlinks

When you are using Nix, especially with
[Home-Manager](https://github.com/nix-community/home-manager/), there are times
when you want to test something or maybe debug some issue in your
configuration. Those times it would be really convenient if you could avoid a
rebuild of your Home-Manager configuration, since this takes time until
evaluation and activation.

For those times I have this small script in my Nix configuration called
`realise-symlinks`, that is defined as:

```nix
{ pkgs, ... }:
let
  realise-symlink = pkgs.writeShellApplication {
    name = "realise-symlink";
    runtimeInputs = with pkgs; [ coreutils ];
    text = ''
      for file in "$@"; do
        if [[ -L "$file" ]]; then
          if [[ -d "$file" ]]; then
            tmpdir="''${file}.tmp"
            mkdir -p "$tmpdir"
            cp --verbose --recursive "$file"/* "$tmpdir"
            unlink "$file"
            mv "$tmpdir" "$file"
            chmod --changes --recursive +w "$file"
          else
            cp --verbose --remove-destination "$(readlink "$file")" "$file"
            chmod --changes +w "$file"
          fi
        else
          >&2 echo "Not a symlink: $file"
          exit 1
        fi
      done
    '';
  };
in
{
  home.packages = [ realise-symlink ];
}
```

The idea of this script is that you can call it against a symlink against Nix
store and it will realise, e.g.: convert to an "actual" file, e.g.:

```console
$ ls -lah .zshrc
lrwxrwxrwx 1 thiagoko users 69 Aug  1 00:10 .zshrc -> /nix/store/glz018yyh0qfqc9lywx1yhr7c3l96lv7-home-manager-files/.zshrc

$ realise-symlink .zshrc
removed '.zshrc'
'/nix/store/glz018yyh0qfqc9lywx1yhr7c3l96lv7-home-manager-files/.zshrc' -> '.zshrc'
mode of '.zshrc' changed from 0444 (r--r--r--) to 0644 (rw-r--r--)

$ ls -lah .zshrc
-rw-r--r-- 1 thiagoko users 5.8K Aug  1 00:16 .zshrc
```

It also add write permissions to the resulting file, to make it easier to edit.
By the way, it also works with directories:

```console
$ ls -lah zim-completion
lrwxrwxrwx 1 thiagoko users 90 Aug  1 00:10 zim-completion -> /nix/store/glz018yyh0qfqc9lywx1yhr7c3l96lv7-home-manager-files/.zsh/plugins/zim-completion

$ realise-symlink zim-completion
'zim-completion/init.zsh' -> 'zim-completion.tmp/init.zsh'
'zim-completion/init.zsh.zwc' -> 'zim-completion.tmp/init.zsh.zwc'
'zim-completion/LICENSE' -> 'zim-completion.tmp/LICENSE'
'zim-completion/README.md' -> 'zim-completion.tmp/README.md'
mode of 'zim-completion/init.zsh' changed from 0444 (r--r--r--) to 0644 (rw-r--r--)
mode of 'zim-completion/init.zsh.zwc' changed from 0444 (r--r--r--) to 0644 (rw-r--r--)
mode of 'zim-completion/LICENSE' changed from 0444 (r--r--r--) to 0644 (rw-r--r--)
mode of 'zim-completion/README.md' changed from 0444 (r--r--r--) to 0644 (rw-r--r--)

$ ls -lah zim-completion
total 28K
drwxr-xr-x 1 thiagoko users   72 Aug  1 00:18 .
drwxr-xr-x 1 thiagoko users  130 Aug  1 00:18 ..
-rw-r--r-- 1 thiagoko users 5.3K Aug  1 00:18 init.zsh
-rw-r--r-- 1 thiagoko users  12K Aug  1 00:18 init.zsh.zwc
-rw-r--r-- 1 thiagoko users 1.3K Aug  1 00:18 LICENSE
-rw-r--r-- 1 thiagoko users 2.2K Aug  1 00:18 README.md
```

After you finish whatever you are testing, to return to your configuration you
can just delete those files and re-run your Home-Manager activation:

```console
$ rm -rf .zshrc

$ sudo systemctl restart home-manager-<user>.service # or `home-manager switch`

$ ls -lah .zshrc
lrwxrwxrwx 1 thiagoko users 69 Aug  1 00:20 .zshrc -> /nix/store/glz018yyh0qfqc9lywx1yhr7c3l96lv7-home-manager-files/.zshrc
```

It even works with system files:

```console
$ sudo realise-symlink /etc/nix/nix.conf
[sudo] password for thiagoko:
removed 'nix.conf'
'/etc/static/nix/nix.conf' -> 'nix.conf'
mode of 'nix.conf' changed from 0444 (r--r--r--) to 0644 (rw-r--r--)
```

But I never needed for this case since it is more rare to me to experiment with
OS level configuration.
