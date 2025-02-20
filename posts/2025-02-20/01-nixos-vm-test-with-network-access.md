# NixOS VM test with network access

The [first post of this
blog](/posts/2024-07-26/01-writing-nixos-tests-for-fun-and-profit.md)
described how I wrote integration tests for a hobby project of mine to interact
with Hyprland (a Window Manager for Wayland).

This time, I had another project of mine that desperately needed a better way
to run integrations tests:
[nix-alien](https://github.com/thiagokokada/nix-alien). [I already had
integration
tests](https://github.com/thiagokokada/nix-alien/blob/7e687663d2054fa1708284bd42731c6be62b1667/integration-tests.nix)
that I wrote a few years ago, but they're basically just a bunch of glorified
shell scripts wrapped in Nix for (some) sanity.

But this time I had much better Nix knowledge and I knew about NixOS VM tests,
so why not port the old tests to use it instead? Since [GitHub
Actions](https://github.com/thiagokokada/hyprland-go/actions/workflows/nix.yaml)
and
[nix-installer-action](https://github.com/DeterminateSystems/nix-installer-action)
supports KVM, it means I can even run the tests inside GitHub Actions for free
(since it is an open-source project).

Taking the knowledge from my previous blog post this was mostly a breeze, and I
got a bootable `flake.nix` file really fast. But then I hit a road-block: how
can I give the VM access to internet?

You see, NixOS VM tests are not really different from any other Nix derivation,
so they're just as isolated. This is great to ensure reproducibility, but it is
annoying sometimes. In `nix-alien` case, one test tries to run `nix-shell`
inside the VM, and this ends up trying to download a copy of
[nixpkgs](https://github.com/NixOS/nixpkgs) tarball. I tried as much as I could
to preload the tarball directly to the VM's `/nix/store`, but nothing worked
and I didn't want to leave the test in the previous state (that wasn't even
working anymore in GitHub Actions thanks to some recent changes).

So I decided to be pragmatic: it is better to have an impure NixOS VM test than
to keep the current state. And the easiest way to fix was to give access to VM
to the internet. But how to do so?

The answer is actually simple, but it is kind puzzling if you don't know where
to look for. First, you need to add support for internet inside the VM. DHCP is
the easiest option and will be shown in the example below, but you can
configure it any other way (e.g.: static IP):

```nix
{
  description = "nix-alien";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { nixpkgs, ... }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in
    {
      checks.${system}.it = pkgs.testers.runNixOSTest {
        name = "integration-tests";

        nodes.machine =
          { pkgs, lib, ... }:
          {
            # Not strictly necessary, but if you are depending in Nix inside
            # your VM tests, this has the effect of increasing reproducibility
            nix.nixPath = [ "${nixpkgs}" ];

            # Here is the important bit, starting the DHCP server
            networking.useDHCP = true;

            # Tweak some VM options
            virtualisation = {
              cores = 2;
              memorySize = 2048;
            };
          };

        testScript = # python
          ''
            start_all()

            # Good to make sure that DHCP client started before testing
            machine.wait_for_unit("multi-user.target")

            # Check if we have network
            machine.succeed("ping -c 3 8.8.8.8")
          '';
      };
    };
}
```

If you put the file above in `flake.nix` and run `nix flake check -L`, the test
will eventually fail with an error similar to this one:

```console
vm-test-run-integration-tests> machine: must succeed: ping -c 3 8.8.8.8
vm-test-run-integration-tests> machine: output: PING 8.8.8.8 (8.8.8.8) 56(84) bytes of data.
vm-test-run-integration-tests> From 10.0.2.2 icmp_seq=1 Destination Net Unreachable
vm-test-run-integration-tests> From 10.0.2.2 icmp_seq=2 Destination Net Unreachable
vm-test-run-integration-tests> From 10.0.2.2 icmp_seq=3 Destination Net Unreachable
vm-test-run-integration-tests> --- 8.8.8.8 ping statistics ---
vm-test-run-integration-tests> 3 packets transmitted, 0 received, +3 errors, 100% packet loss, time 2042ms
vm-test-run-integration-tests> cleanup
vm-test-run-integration-tests> kill machine (pid 11)
vm-test-run-integration-tests> qemu-system-x86_64: terminating on signal 15 from pid 8 (/nix/store/0l539chjmcq5kdd43j6dgdjky4sjl7hl-python3-3.12.8/bin/python3.12)
vm-test-run-integration-tests> kill vlan (pid 9)
vm-test-run-integration-tests> (finished: cleanup, in 0.02 seconds)
vm-test-run-integration-tests> Traceback (most recent call last):
vm-test-run-integration-tests>   File "/nix/store/ahpc056hlclhnv4qrdlfb525pk3shnxw-nixos-test-driver-1.1/bin/.nixos-test-driver-wrapped", line 9, in <module>
vm-test-run-integration-tests>     sys.exit(main())
vm-test-run-integration-tests>              ^^^^^^
vm-test-run-integration-tests>   File "/nix/store/ahpc056hlclhnv4qrdlfb525pk3shnxw-nixos-test-driver-1.1/lib/python3.12/site-packages/test_driver/__init__.py", line 146, in main
vm-test-run-integration-tests>     driver.run_tests()
vm-test-run-integration-tests>   File "/nix/store/ahpc056hlclhnv4qrdlfb525pk3shnxw-nixos-test-driver-1.1/lib/python3.12/site-packages/test_driver/driver.py", line 174, in run_tests
```

What gives? Well, this is the Nix sandbox in action. It is how Nix ensure
reproducibility, but we don't want it in this case. While you can disable it in
the Nix configuration, the easiest way is to simple disable it during the `nix`
call:

```console
$ nix flake check -L --option sandbox false
```

And everything works as except:

```console
vm-test-run-integration-tests> machine: must succeed: ping -c 3 8.8.8.8
vm-test-run-integration-tests> (finished: must succeed: ping -c 3 8.8.8.8, in 2.08 seconds)
vm-test-run-integration-tests> (finished: run the VM test script, in 16.77 seconds)
vm-test-run-integration-tests> test script finished in 16.87s
vm-test-run-integration-tests> cleanup
vm-test-run-integration-tests> kill machine (pid 1470953)
vm-test-run-integration-tests> qemu-system-x86_64: terminating on signal 15 from pid 1470949 (/nix/store/0l539chjmcq5kdd43j6dgdjky4sjl7hl-python3-3.12.8/bin/python3.12)
vm-test-run-integration-tests> kill vlan (pid 1470951)
vm-test-run-integration-tests> (finished: cleanup, in 0.01 seconds)
```
