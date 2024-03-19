# [Getting started with eBPF](https://isovalent.com/labs/ebpf-getting-started/)

BPF is a revolutionary technology with origins in the Linux kernel that can run sandboxed programs in an operating system kernel. It is used to safely and efficiently extend the capabilities of the kernel without requiring to change kernel source code or load kernel modules.

## 🐧 The JavaScript of the Kernel
eBPF is a kernel technology that allows to dynamically extend the functionalities of the Linux kernel at runtime.

You can think of it as what JavaScript is to the web browser: JavaScript lets you attach callbacks to events in the DOM in order to bring dynamic features to your web page. In a similar fashion, eBPF allows to hook to kernel events and extend their logic when these events are triggered!

![img](https://play.instruqt.com/assets/tracks/foixleio3jpg/da954a22a27ea15d468b0321ff679d08/assets/ebpf_javascript.png)

For example, when a process creates a new process, it calls the execve syscall, which normally results in scheduling the new process execution in the kernel. With eBPF, you can attach a program to that event and use it to act upon it, for example for observability.

## 👮🏻‍♀ Verification and JIT Compilation

Why would eBPF be better than existing solutions to extend the kernel such as kernel modules? Besides the fact that it allows an event-driven approach to kernel development, it is also inherently more secure and safe.

That is because eBPF programs are verified when they are injected into the kernel.

![img](https://play.instruqt.com/assets/tracks/foixleio3jpg/545f3cd7c7ff1d170aa41e82e1cbb7fa/assets/eBPF_animated_bg.gif)

After that first step, the eBPF code can be JIT-compiled into machine code, and then attached to kernel hooks and eBPF maps. When these hooks are triggered by events, the code will be executed.

## 🏛 What You Will Do in This Course
To get first hand experiences with eBPF, in this lab we will:

Build and use opensnoop, an eBPF based tool that reports whenever a file is opened
Use readelf to compare a BPF object file with its source code
Use bpftool to see how your tool is loaded into the kernel
Add additional “hello world”-style tracing to the source code
Build and run it again to see your own customized eBPF tracing The example tool we're using, opensnoop, is one of a collection of eBPF-based tools from the BCC project.