# Quick bits: nix-shell is cursed

The other day I had to run a PHP project in my machine. I have no idea how PHP
ecosystem work, I just wanted to get it to run.

The easiest way to get a script to run if you use Nix is to use `nix-shell`. As
many of you probably know, you can add `nix-shell` as a shebang in your scripts
to run them as `./script`. This was a PHP script so I wanted to do the same.
Easy right?

```php
#!/use/bin/env nix-shell
#!nix-shell -i php -p php83
<?php
declare(strict_types=1);
```

And:

```console
$ ./index.php
Fatal error: strict_types declaration must be the very first statement in the script in index.php on line 4
```

So it seems that `declare(strict_types=1)` needs to be the first line in a PHP
script if used. I removed `declare(strict_types=1)` and while the script works,
I don't have enough expertise in PHP to know if this would be safe or not.

I decided to try something that initially looked really dumb:

```php
#!/use/bin/env nix-shell
<?php
declare(strict_types=1);
#!nix-shell -i php -p php83
```

And:

```console
$ ./index.php
Works
```

Wat? I mean, it is not dumb if it works, but this at least looks cursed.

Eventually I found this
[comment](https://github.com/NixOS/nix/issues/2570#issuecomment-446220517) in a
Nix issue talking about cases where `nix-shell` shebang doesn't work. It looks
like the classic case of a [bug that becomes a
feature](https://github.com/NixOS/nix/issues/2570#issuecomment-446222206).

Update: after posting this in
[Lobte.rs](https://lobste.rs/s/gkcgza/quick_bits_nix_shell_is_cursed), it seems
someone decided to open a [Pull
Request](https://github.com/NixOS/nix/pull/11202) to document this behavior.
Also the equivalent for the new [nix
CLI](https://nix.dev/manual/nix/2.23/command-ref/new-cli/nix#shebang-interpreter)
explicitly documents this behavior:

> Note that the `#! nix` lines don't need to follow after the first line, to
> accomodate other interpreters.
