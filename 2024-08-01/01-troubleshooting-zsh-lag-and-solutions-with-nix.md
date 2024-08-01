# Troubleshoting: ZSH lag and solutions with Nix

Inspired by this [blog post from
Tavis](https://lock.cmpxchg8b.com/slowterm.html), I decided to document my own
recent journey of reducing terminal (ZSH) lag startup. This post is way less
interesting than the one from Tavis that uses a debugger to patch applications
on the fly, but should still be interesting for some. And it also shows how
powerful Nix can be for some things.

For context, I have basically 3 systems where I interact with terminal
frequently:

- Thinkpad P14s Gen 1 running NixOS, with a reasonable fast CPU and disk
- [MacBook Pro "M1
  Pro"](https://everymac.com/systems/apple/macbook_pro/specs/macbook-pro-m1-pro-10-core-cpu-16-core-gpu-16-2021-specs.html)
  (what an awful name scheme Apple) with a really fast CPU and disk, but of
  course running macOS
  + Sadly this is being phased-out since this is a job owned machine and I am
    changing jobs right now, but should be replaced with another one soon™
- [Chromebook Duet
  3](https://chromeunboxed.com/lenovo-chromebook-duet-3-review-perfect-sequel)
  running ChromeOS, with slow CPU and disk

My experience is similar to Tavis, at around 300ms of startup time I don't care
too much, but around 500ms+ is where I start to notice. I never had any issues
with startup time in NixOS itself (I had issues with macOS before, but it was
not actually the fault of macOS), but in the Chromebook it was awful: 600ms+
with [hot
start](https://www.instabug.com/blog/understanding-cold-hot-and-warm-app-launch-time),
while cold start it could take multiple seconds.

We can check how long ZSH takes to start by using:

```console
$ time zsh -ci exit
zsh -ic exit  0.04s user 0.10s system 100% cpu 0.143 total
```

The `-i` flag here is important, because we are interested in the interactive
use of ZSH. Without this flag ZSH will ignore your `~/.zshrc` file, and the
results will be meaningless.

To do a more interesting benchmark, we can use
[`hyperfine`](https://github.com/sharkdp/hyperfine):

```console
$ hyperfine "zsh -ic exit"
Benchmark 1: zsh -ic exit
  Time (mean ± σ):     145.4 ms ±   4.2 ms    [User: 49.8 ms, System: 97.3 ms]
  Range (min … max):   138.6 ms … 155.3 ms    19 runs
```

Hyperfine will run the command multiple times and take care of things like
shell startup time. A really great tool to have in your toolbox by the way, but
I digress.

So let's do a little time travelling. Going back to commit
[`b12757f`](https://github.com/thiagokokada/nix-configs/tree/b12757f90889653e359a1ab0a8cfd2f90cfabf31)
from [nix-configs](https://github.com/thiagokokada/nix-configs/). Running
`hyperfine` like above from my NixOS laptop, we have:

```console
$ hyperfine "zsh -ic exit"
Benchmark 1: zsh -ic exit
  Time (mean ± σ):     218.6 ms ±   5.1 ms    [User: 70.6 ms, System: 151.5 ms]
  Range (min … max):   210.3 ms … 227.0 ms    13 runs
```

This doesn't look that bad, but let's see the same commit in my Chromebook:

```console
$ hyperfine "zsh -ic exit"
Benchmark 1: zsh -ic exit
  Time (mean ± σ):     679.7 ms ±  40.2 ms    [User: 230.8 ms, System: 448.5 ms]
  Range (min … max):   607.3 ms … 737.0 ms    10 runs
```

Yikes, this is much worse. And those are the results after I retried the
benchmark (so it is a hot start). The cold start times were above 3s. So let's
investigate what is happening here. We can profile what is taking time during
the startup of ZSH using [zprof](https://www.bigbinary.com/blog/zsh-profiling).
You can add the following in your `~/.zshrc`:

```bash
# At the top of your ~/.zshrc file
zmodload zsh/zprof

# ...

# At the end of your ~/.zshrc file
zprof
```

Or if using Home-Manager, use the
[`programs.zsh.zprof.enable`](https://nix-community.github.io/home-manager/options.xhtml#opt-programs.zsh.zprof.enable)
option. Once we restart ZSH, we will have something like:

```console
num  calls                time                       self            name
-----------------------------------------------------------------------------------
 1)    1          36.91    36.91   34.29%     30.47    30.47   28.31%  (anon) [/home/thiagoko/.zsh/plugins/zim-completion/init.zsh:13]
 2)    1          25.43    25.43   23.63%     25.43    25.43   23.63%  (anon) [/home/thiagoko/.zsh/plugins/zim-ssh/init.zsh:6]
 3)    1          22.00    22.00   20.45%     21.92    21.92   20.36%  _zsh_highlight_load_highlighters
 4)    1          12.32    12.32   11.45%     12.32    12.32   11.45%  autopair-init
 5)    1           6.44     6.44    5.98%      6.44     6.44    5.98%  compinit
 6)    1           3.56     3.56    3.31%      3.48     3.48    3.23%  prompt_pure_state_setup
 7)    2           3.79     1.89    3.52%      2.85     1.43    2.65%  async
 8)    1           0.93     0.93    0.87%      0.93     0.93    0.87%  async_init
 9)    6           0.93     0.15    0.86%      0.93     0.15    0.86%  is-at-least
10)    6           0.67     0.11    0.63%      0.67     0.11    0.63%  add-zle-hook-widget
11)    1           8.25     8.25    7.66%      0.61     0.61    0.57%  prompt_pure_setup
12)    1           0.40     0.40    0.37%      0.40     0.40    0.37%  (anon) [/nix/store/p1zqypy7600fvfyl1v571bljx2l8zhay-zsh-autosuggestions-0.7.0/share/zsh-autosuggestions/zsh-autosuggestions.zsh:458]
13)    5           0.31     0.06    0.29%      0.31     0.06    0.29%  add-zsh-hook
14)    1           0.60     0.60    0.56%      0.29     0.29    0.27%  (anon) [/home/thiagoko/.zsh/plugins/zim-input/init.zsh:5]
15)    1           0.21     0.21    0.20%      0.21     0.21    0.20%  compdef
16)    1           0.10     0.10    0.09%      0.10     0.10    0.09%  _zsh_highlight__function_is_autoload_stub_p
17)    1           0.26     0.26    0.24%      0.08     0.08    0.08%  _zsh_highlight__function_callable_p
18)    1           0.08     0.08    0.08%      0.08     0.08    0.08%  prompt_pure_is_inside_container
19)    1           0.07     0.07    0.07%      0.07     0.07    0.07%  _zsh_highlight__is_function_p
20)    1           0.01     0.01    0.01%      0.01     0.01    0.01%  __wezterm_install_bash_prexec
21)    1           0.00     0.00    0.00%      0.00     0.00    0.00%  _zsh_highlight_bind_widgets
# ...
```

I ommited some output for brevit. The first 2 things that shows are from the
[zimfw](https://github.com/zimfw/zimfw), the framework that I use to configure
my ZSH (similar to Oh-My-Zsh). I actually don't use `zimfw` directly, instead I
just load some modules that I find useful, like the `zim-completion` and
`zim-ssh` that we can see above. By the way, Zim is generally really well
optimised for startup time, but those 2 modules are kind slow.

For [`zim-completion`](https://github.com/zimfw/completion), after taking a
look at it, there isn't much I could do. It seems that the reason
`zim-completion` takes so long during startup is because it is trying to decide
if it needs to recompile the completions (and replacing it with just a naive
`autoload -U compinit && compinit` is even worse for startup performance). I
may eventually replace it for something else, but I really like what Zim brings
here, so I decided to not touch it for now.

However [`zim-ssh`](https://github.com/zimfw/ssh) is another history. The only
reason I used it is to start a `ssh-agent` and keep it between multiple ZSH
sessions. It shouldn't have that much influence in startup time. So I took a
look the code (since it is small, I am reproducing it here):

```bash
#
# Set up ssh-agent
#

# Don't do anything unless we can actually use ssh-agent
(( ${+commands[ssh-agent]} )) && () {
  ssh-add -l &>/dev/null
  if (( ? == 2 )); then
    # Unable to contact the authentication agent

    # Load stored agent connection info
    local -r ssh_env=${HOME}/.ssh-agent
    if [[ -r ${ssh_env} ]] source ${ssh_env} >/dev/null

    ssh-add -l &>/dev/null
    if (( ? == 2 )); then
        # Start agent and store agent connection info
        (umask 066; ssh-agent >! ${ssh_env})
        source ${ssh_env} >/dev/null
    fi
  fi

  # Load identities
  ssh-add -l &>/dev/null
  if (( ? == 1 )); then
    local -a zssh_ids
    zstyle -a ':zim:ssh' ids 'zssh_ids'
    if (( ${#zssh_ids} )); then
      ssh-add ${HOME}/.ssh/${^zssh_ids} 2>/dev/null
    else
      ssh-add 2>/dev/null
    fi
  fi
}
```

Well, this is bad. Let's assume the common path, where the `ssh-agent` is
already running but you open a new shell instance (that doesn't have the
connection info yet so it will need to load). This will run `ssh-add` at 4
times. How long does `ssh-add` takes to run?

```console
$ hyperfine -Ni "ssh-add -l"
Benchmark 1: ssh-add -l
  Time (mean ± σ):       4.6 ms ±   1.1 ms    [User: 2.0 ms, System: 2.0 ms]
  Range (min … max):     3.4 ms …   8.7 ms    619 runs

  Warning: Ignoring non-zero exit code.
```

For those curious, `-N` disables the Shell usage, that works better when the
command being tested is too fast.

In average we have 4x4ms=16ms of startup time. But keep in mind the worst case
can be much worse. The question is, how can we improve the situation here?

After taking a look, I decided to write my own code, based in some ideas stolen
from [Oh-My-Zsh ssh-agent
plugin](https://github.com/ohmyzsh/ohmyzsh/blob/67581c53c6458566e174620361e84b364b9034d2/plugins/ssh-agent/ssh-agent.plugin.zsh).
Here is final version of my
[code](https://github.com/thiagokokada/nix-configs/blob/e45a888f2bf3ce5644c3966f0b6371414d0291e2/home-manager/cli/ssh/ssh-agent.zsh):

```bash
zmodload zsh/net/socket

_check_agent(){
  if [[ -S "$SSH_AUTH_SOCK" ]] && zsocket "$SSH_AUTH_SOCK" 2>/dev/null; then
    return 0
  fi
  return 1
}

_start_agent() {
  # Test if $SSH_AUTH_SOCK is visible, in case we start e.g.: ssh-agent via
  # systemd service
  if _check_agent; then
    return 0
  fi

  # Get the filename to store/lookup the environment from
  local -r ssh_env_cache="$HOME/.ssh-agent"

  # Check if ssh-agent is already running
  if [[ -f "$ssh_env_cache" ]]; then
    source "$ssh_env_cache" > /dev/null

    # Test if $SSH_AUTH_SOCK is visible, e.g.: the ssh-agent is still alive
    if _check_agent; then
      return 0
    fi
  fi

  # start ssh-agent and setup environment
  (
    umask 066
    ssh-agent -s >! "$ssh_env_cache"
  )
  source "$ssh_env_cache" > /dev/null
}

_start_agent
unfunction _check_agent _start_agent
```

The idea here is simple: using
[`zsocket`](https://zsh.sourceforge.io/Doc/Release/Zsh-Modules.html#The-zsh_002fnet_002fsocket-Module)
module from ZSH itself to check if the `ssh-agent` is working instead of
executing `ssh-add -l`. The only case we run any program now is to start the
agent itself if needed. Let's run `hyperfine` again:

```
$ hyperfine "zsh -ic exit"
Benchmark 1: zsh -ic exit
  Time (mean ± σ):     188.3 ms ±   8.2 ms    [User: 61.1 ms, System: 130.0 ms]
  Range (min … max):   170.9 ms … 198.4 ms    16 runs
```

Got a good improvement here already. Let's see `zprof` again:

```console
num  calls                time                       self            name
-----------------------------------------------------------------------------------
 1)    1          41.23    41.23   48.66%     33.52    33.52   39.56%  (anon) [/home/thiagoko/.zsh/plugins/zim-completion/init.zsh:13]
 2)    1          22.23    22.23   26.24%     22.12    22.12   26.10%  _zsh_highlight_load_highlighters
 3)    1           8.90     8.90   10.51%      8.90     8.90   10.51%  Gautopair-init
 4)    1           7.71     7.71    9.10%      7.71     7.71    9.10%  compinit
 5)    1           5.74     5.74    6.77%      5.60     5.60    6.60%  prompt_pure_state_setup
 6)    6           1.19     0.20    1.41%      1.19     0.20    1.41%  add-zle-hook-widget
 7)    2           1.97     0.99    2.33%      1.14     0.57    1.34%  async
 8)    6           0.87     0.15    1.03%      0.87     0.15    1.03%  is-at-least
 9)    1           0.84     0.84    0.99%      0.84     0.84    0.99%  async_init
10)    1           9.30     9.30   10.97%      0.72     0.72    0.84%  prompt_pure_setup
11)    5           0.63     0.13    0.75%      0.63     0.13    0.75%  add-zsh-hook
12)    1           0.41     0.41    0.48%      0.41     0.41    0.48%  _start_agent
13)    1           0.31     0.31    0.37%      0.31     0.31    0.37%  (anon) [/nix/store/p1zqypy7600fvfyl1v571bljx2l8zhay-zsh-autosuggestions-0.7.0/share/zsh-autosuggestions/zsh-autosuggestions.zsh:458]
14)    1           0.55     0.55    0.64%      0.24     0.24    0.28%  (anon) [/home/thiagoko/.zsh/plugins/zim-input/init.zsh:5]
15)    1           0.14     0.14    0.16%      0.14     0.14    0.16%  prompt_pure_is_inside_container
16)    1           0.14     0.14    0.16%      0.14     0.14    0.16%  compdef
17)    1           0.09     0.09    0.11%      0.09     0.09    0.11%  _zsh_highlight__function_is_autoload_stub_p
18)    1           0.25     0.25    0.29%      0.08     0.08    0.09%  _zsh_highlight__function_callable_p
19)    1           0.07     0.07    0.09%      0.07     0.07    0.09%  _zsh_highlight__is_function_p
20)    1           0.01     0.01    0.01%      0.01     0.01    0.01%  __wezterm_install_bash_prexec
21)    1           0.01     0.01    0.01%      0.01     0.01    0.01%  _zsh_highlight_bind_widgets
# ...
```

Well, there is nothing interesting here anymore. I mean, `zim-completion` is
still the main culprit, but nothing to do for now. Instead of looking at
`zproof`, let's take a look at my `~/.zshrc` instead:

```bash
# ...
if [[ $options[zle] = on ]]; then
  eval "$(/nix/store/sk6wsgp4h477baxypksz9rl8ldwwh9yg-fzf-0.54.0/bin/fzf --zsh)"
fi

# ...
/nix/store/x3yblr73r5x76dmaanjk3333mvzxc49r-any-nix-shell-1.2.1/bin/any-nix-shell zsh | source /dev/stdin

# ...
eval "$(/nix/store/330d6k81flfs6w46b44afmncxk57qggv-zoxide-0.9.4/bin/zoxide init zsh )"

# ...
eval "$(/nix/store/8l9j9kdv9m0z0s30lp4yvrc9s5bcbgmx-direnv-2.34.0/bin/direnv hook zsh)"
```

So you see, starting all those programs during ZSH startup can hurt the shell
startup considerable. Not necessary for commands fast like `fzf` (that is
written in Go), but let's see
[`any-nix-shell`](https://github.com/haslersn/any-nix-shell), that is written
in shell script:

```console
$ hyperfine "any-nix-shell zsh"
Benchmark 1: any-nix-shell zsh
  Time (mean ± σ):      16.0 ms ±   1.8 ms    [User: 5.6 ms, System: 10.5 ms]
  Range (min … max):    11.3 ms …  20.3 ms    143 runs
```

This is bad, consistently bad actually. Even for commands that are fast, keep
in mind that there is a difference between the cold and hot start again. For
example, `fzf`:

```console
$ hyperfine -N "fzf --zsh"
Benchmark 1: fzf --zsh
  Time (mean ± σ):       2.9 ms ±   0.9 ms    [User: 0.6 ms, System: 2.3 ms]
  Range (min … max):     1.7 ms …   6.8 ms    1113 runs
```

See the range? While 1.7ms is something that is probably difficult to notice,
6.8ms can be noticiable, especially if this accumulates with other slow
starting apps.

And the thing is, all those commands are just generating in the end a fixed
output, at least for the current version of the program. Can we pre-generate
them instead? If using Nix, of course we can:

```nix
programs.zsh.initExtra =
  # bash
  ''
    # any-nix-shell
    source ${
      pkgs.runCommand "any-nix-shell-zsh" { } ''
        ${lib.getExe pkgs.any-nix-shell} zsh > $out
      ''
    }

    # fzf
    source ${config.programs.fzf.package}/share/fzf/completion.zsh
    source ${config.programs.fzf.package}/share/fzf/key-bindings.zsh

    # zoxide
    source ${
      pkgs.runCommand "zoxide-init-zsh" { } ''
        ${lib.getExe config.programs.zoxide.package} init zsh > $out
      ''
    }

    # direnv
    source ${
      pkgs.runCommand "direnv-hook-zsh" { } ''
        ${lib.getExe config.programs.direnv.package} hook zsh > $out
      ''
    }
  '';
```

So we can use `pkgs.runCommand` to run those commands during build time and
`source` the result. `fzf` actually doesn't need this since we have the files
already generated in the package. I think this is one of those things that
really shows the power of Nix: I wouldn't do something similar if I didn't use
Nix because the risk of breaking something later is big (e.g.: forgetting to
update the generated files), but Nix makes those things trivial.

Let's run `hyperfine` again:

```
$ hyperfine "zsh -ic exit"
Benchmark 1: zsh -ic exit
  Time (mean ± σ):     162.3 ms ±   4.9 ms    [User: 52.7 ms, System: 111.1 ms]
  Range (min … max):   153.0 ms … 173.4 ms    19 runs
```

Another good improvement. The last change I did is switching between
[`zsh-syntax-highlighting`](https://github.com/zsh-users/zsh-syntax-highlighting)
to
[`zsh-fast-syntax-highlighting`](https://github.com/zdharma-continuum/fast-syntax-highlighting),
that is supposed to be faster and have better highlighting too. I got that from
`_zsh_highlight_load_highlighters` using 26% of the time from my `zprof` above.
And for the final `hyperfine` in my laptop:

```
$ hyperfine "zsh -ic exit"
Benchmark 1: zsh -ic exit
  Time (mean ± σ):     138.3 ms ±   7.1 ms    [User: 47.5 ms, System: 91.9 ms]
  Range (min … max):   123.8 ms … 157.9 ms    21 runs
```

A ~36% improvement, not bad. Let's see how it fares in my Chromebook:

```
$ hyperfine "zsh -ic exit"
Benchmark 1: zsh -ic exit
  Time (mean ± σ):     278.2 ms ±  46.9 ms    [User: 88.0 ms, System: 184.8 ms]
  Range (min … max):   204.7 ms … 368.5 ms    11 runs
```

An even more impressive ~59% improvement. And yes, the shell startup now feels
much better.
