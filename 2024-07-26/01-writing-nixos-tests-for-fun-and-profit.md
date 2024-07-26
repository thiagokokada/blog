# Writing NixOS tests for fun and profit

I recently started a [new side
project](https://github.com/thiagokokada/hyprland-go) writing an IPC library in
Go for [Hyprland](https://hyprland.org/), a Window Manager for Wayland.

Once I got past the Work-in-Progress phase, I realise I had an issue: I wrote
some tests, but I was running then inside my system running Hyprland. And the
tests themselves were annoying: since they send commands to the current running
Hyprland instance, I was having programs being opened and settings being
changed, because this was the only way to have a reasonable good confidence
that what I was doing was correct. So I need to do like any good developer and
implement a CI, but how?

One approach would be to create something like a mock client and test against
my mock. Since this mock wouldn't need a running Hyprland instance the tests
could run everywhere (even in non-Linux systems!), but they wouldn't be much
useful. Mocks are great for testing business logic, but not really for making
sure everything is working correctly.

I need something more akin to an integration test, but this is tricky. It is
not like I am doing integration with e.g.: PostgreSQL that has thousands of
libraries available to make integration tests easier, I am doing integration
with a Window Manager that is a moving target with multiple breaking changes in
each release. And this is where NixOS tests enter, a way to run tests inside
Virtual Machines configured in Nix.

I am a long time NixOS user and commiter, but I never wrote a NixOS test
outside of [nixpkgs](https://github.com/NixOS/nixpkgs) itself. However I knew
it was possible, and after doing a quick reading of the [Wiki
entry](https://wiki.nixos.org/wiki/NixOS_VM_tests) about it, I was ready to
start.

The first part is to call `pkgs.nixosTest` and configure the machine as any
other NixOS system, e.g.:

```nix
{ pkgs, ... }:
pkgs.nixosTest {
  name = "hyprland-go";

  nodes.machine =
    { config, pkgs, lib, ... }:
    {
      # bootloader related configuration
      boot.loader.systemd-boot.enable = true;
      boot.loader.efi.canTouchEfiVariables = true;

      # enable hyprland
      programs.hyprland.enable = true;

      # create a user called alice
      users.users.alice = {
        isNormalUser = true;
      };

      # add some extra packages that we need during tests
      environment.systemPackages = with pkgs; [ go kitty ];

      # auto login as alice
      services.getty.autologinUser = "alice";

      # configure VM, increase memory and CPU and enable OpenGL via LLVMPipe
      virtualisation.qemu = {
        options = [
          "-smp 2"
          "-m 4G"
          "-vga none"
          "-device virtio-gpu-pci"
        ];
      };

      # Start hyprland at login
      programs.bash.loginShellInit = "Hyprland";
    };

  testScript = "start_all()";
}
```

A few details that I want to bring to attention. The first one is how easy it
is to setup things like a normal user account, add some extra packages we need
for testing, add Hyprland itself and configure auto-login. I have no idea how
painful it would be to automatise all those steps in e.g.: Ansible, but here we
are in a few lines of Nix code. This is, of course, thanks to all the
contributors to nixpkgs that implement something that help their own use case,
but once combined make it greater than the sum of the parts.

Second is something that I took a while to figure out: how to enable GPU
acceleration inside the VM. You see, Hyprland, different from other Window
Managers, requires OpenGL support. This is basically why the flag `-device
virtio-gpu-pci` is in `virtualisation.qemu.options`, this enables OpenGL
rendering via LLVMPipe, that while being slow since it is rendered in CPU, is
sufficient for this case.

Putting the above code inside a
[`flake.nix`](https://wiki.nixos.org/wiki/Flakes) for reproducibility, I had
something similar to:

```nix
{
  description = "Hyprland's IPC bindings for Go";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { nixpkgs, ... }:
    {
      checks.x86_64-linux =
        let
          pkgs = import nixpkgs { system = "x86_64-linux"; };
        in
        {
          testVm = pkgs.nixosTest {
            # the code above
          };
        }
    };
}
```

I can now run `nix build .#checks.x86_64-linux.testVm -L` to build and run the
VM. However it is not really useful right now, since we didn't add any useful
code in `testScript`, the core of the NixOS test framework. We can also run
`nix build .#checks.x86_64-linux.testVm.driverInteractive` and
`./result/bin/nixos-test-driver`: this will start a Python console where we can
manually play with the VM (try typing `start_all()` for example).

The `testScript` is a sequence of Python statements that perform various
actions, such as starting VMs, executing commands in the VMs, and so on. More
about it in the official
[documentation](https://nixos.org/manual/nixos/stable/index.html#sec-nixos-tests).
For our case we can start with something like this:

```nix
{
    textScript = /* python */ ''
      start_all()

      machine.wait_for_unit("multi-user.target")
      machine.wait_for_file("/home/alice/test-finished")
    '';
}
```

The first statement, `start_all()`, starts all VMs, in this case we have only
one, called `machine`. We send two further commands to `machine`:
`wait_for_unit("multi-user.target")` and
`wait_for_file("/home/alice/test-finished")`.

The first command waits until systemd's `multi-user.target` is ready, a good
way to ensure that the system is ready for further commands. The second one we
wait for a file called `test-finished` to appear in Alice's `$HOME` (basically,
a canary), but how can we generate this file?

Remember that we added `programs.bash.loginShellInit = "Hyprland"`, that
automatically starts Hyprland when Alice logs in. We need to modify that
command to run the Go tests from our library. The good thing is that Hyprland
configuration file supports a
[`exec-once`](https://wiki.hyprland.org/Configuring/Keywords/#executing)
command that runs a command during Hyprland launch. We can abuse this to launch
a terminal emulator and run our tests:

```nix
{
  programs.bash.loginShellInit =
    let
      testScript = pkgs.writeShellScript "hyprland-go-test" ''
        set -euo pipefail

        trap 'echo $? > $HOME/test-finished' EXIT # creates the canary when the script finishes

        cd ${./.} # go to the library directory
        go test -v ./... > $HOME/test.log 2>&1 # run Go tests
      '';
      hyprlandConf = pkgs.writeText "hyprland.conf" ''
        exec-once = kitty sh -c ${testScript}
      '';
    in ''
      Hyprland --config ${hyprlandConf}
    '';
}
```

So we are basically creating a custom Hyprland config that starts a
[Kitty](https://sw.kovidgoyal.net/kitty/) terminal emulator, that then launches
a shell script that runs the test. Since we have no way to get the results of
the test, we pipe the output to a file that we can collect later (e.g.:
`machine.succeded("cat /home/alice/test.log")`). And once the script exit, we
create the canary file `$HOME/test-finished`, that allows the `testScript`
knows that the test finished and it can destroy the VM safely.

If you want to take a look at the final result, it is
[here](https://github.com/thiagokokada/hyprland-go/blob/v0.0.1/flake.nix). This
tests run in any Linux machine that supports KVM, and also works in [GitHub
Actions](https://github.com/thiagokokada/hyprland-go/actions/workflows/nix.yaml)
thanks to the the
[nix-installer-action](https://github.com/DeterminateSystems/nix-installer-action).

And now I have a proper CI pipeline in a way that I never imagined would be
possible, especially considering how simple it was.
