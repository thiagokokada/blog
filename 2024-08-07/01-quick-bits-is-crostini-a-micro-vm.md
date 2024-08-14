# Quick bits: is Crostini a micro VM?

**Disclaimer**: I am not a Virtual Machine specialist, this post is mainly the
conclusion I got after some searching.

Someone asked me in my [previous
post](/2024-08-05/01-my-favorite-device-is-a-chromebook.md) about my Chromebook
if Crostini could be considered a micro VM. This was a interesting question, so
I decided to do another post.

To start, it is really difficult to get a good definition of what a "micro VM"
is. [Firecracker](https://firecracker-microvm.github.io/) defines itself as a
micro VM, and describes itself in its website as:

> Firecracker is a virtual machine monitor (VMM) that uses the Linux
> Kernel-based Virtual Machine (KVM) to create and manage microVMs. Firecracker
> has a minimalist design. It excludes unnecessary devices and guest
> functionality to reduce the memory footprint and attack surface area of each
> microVM. This improves security, decreases the startup time, and increases
> hardware utilization.

Now looking at Crostini, its heart is a VMM called
[crosvm](https://crosvm.dev/). It is described in its
[README](https://chromium.googlesource.com/chromiumos/platform/crosvm/+/HEAD/README.md)
as:

> crosvm is a virtual machine monitor (VMM) based on Linux’s KVM hypervisor,
> with a focus on simplicity, security, and speed. crosvm is intended to run
> Linux guests, originally as a security boundary for running native
> applications on the ChromeOS platform. Compared to QEMU, crosvm doesn’t
> emulate architectures or real hardware, instead concentrating on
> paravirtualized devices, such as the virtio standard.

Similar descriptions right? Actually Firecracker website says it "started from
Chromium OS's Virtual Machine Monitor, crosvm, an open source VMM written in
Rust". So I would say it is safe to say crosvm itself is a micro VM.

But
[Crostini](https://www.chromium.org/chromium-os/developer-library/guides/containers/containers-and-vms/)
itself is a combination of virtualization AND containerization. Basically
inside the VM it runs a Linux kernel and [LXC](https://linuxcontainers.org/),
that can start arbitrary containers inside it. From the Crostini documentation
this choice seems to be to keep startup times down, and also to increase
security (e.g.: in case of a security issue inside the container).

This is definitely an interesting choice, since containers allow the overhead
of each distro that you run inside Crostini to be low, and the main VM itself
(called
[Termina](https://chromium.googlesource.com/chromiumos/overlays/board-overlays/+/HEAD/project-termina/))
should have low overhead too thanks to crosvm.

By the way, if you want to learn more how "devices" works inside a micro VM
like crosvm, I recommend [this blog
post](https://prilik.com/blog/post/crosvm-paravirt/) talking about
paravirtualized devices in crosvm.
