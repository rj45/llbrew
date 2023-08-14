# llbrew

An LLVM backend designed specifically to make it easy to retarget to new instruction set architectures / CPUs.

Input LLVM IR compiled using clang, tinygo, etc for a CPU similar to yours (at least needs the same word size and pointer size), and it will output [customasm](https://github.com/hlorenzi/customasm) compatible assembly. Yes that's right, C or Go or Rust for your CPU you designed and built yourself.

At least that's the goal anyway, currently this is very much a work in progress with a long way to go.

If you want to help, please get in touch.

# Status

Currently it can compile [very simple C programs](./testsuite/) into [customasm](https://github.com/hlorenzi/customasm) assembly.

There's a lot of work to do before it can compile all programs.

# Goals

Make it as simple as possible to get a C compiler for a custom new CPU up and running. TinyGo and Rust are also worthy goals.

If anything isn't as simple as it could be, please open an Issue.

# Non-Goals

Fast compilation. Most of the faster algorithms are very complex and would greatly hinder making this as simple as possible. Happy to entertain faster algorithms that happen to also be simpler, however.

# How does it work?

LLVM has an intermediate representation called "LLVM IR". Any language that uses LLVM to compile its code should be able to emit either the textual representation of LLVM IR or "bitcode" which is a binary representation of the same. In fact there's a project [`gllvm`](https://github.com/SRI-CSL/gllvm) that can be used in the place of `clang` in order to compile and link an entire C project into LLVM IR.

`llbrew` takes the LLVM IR, loads it up using the LLVM libraries, optimizes it using LLVM's own optimization passes, then translates it out of LLVM IR into llbrew IR, performs a series of transformation passes on it to "lower" it into the instruction set of the target CPU, does register allocation, some post-allocation passes, then finally emits [customasm](https://github.com/hlorenzi/customasm) assembly.

The llbrew IR is loosely based on CraneLift's IR. It's designed to support both LLVM's instruction set as well as the target instruction set at the same time, so it can be incrementally translated from one form to the next. It also supports register allocation while still in SSA form, so post-allocation transformations work exactly the same as pre-allocation ones. Transforms are written in plain Go.

# Building

Currently [TinyGo's LLVM bindings](https://github.com/tinygo-org/go-llvm) are used to link with a system installed LLVM. This should work on Mac out of the box (though you may need to `brew install llvm@15`), and on linux you will need to install LLVM 13, 14 or 15. AFAIK LLVM 16 is not yet supported.

Then, in theory, `go install github.com/rj45/llbrew` should work.

Getting this to work on Windows may be more involved. If you would like to help with this, it would be much appreciated.

# History

This is the continuation of the [`NanoGo`](https://github.com/rj45/nanogo) project. One day it occurred to me that I could reuse 80% of the work I had done on NanoGo, but also support C, C++ and Rust in addition to Go, if I just take LLVM IR as input instead of only Go. I had already copied a fair bit of [`TinyGo`](https://tinygo.org)'s code just trying to get the basics going in NanoGo, but with this transition I don't need to reinvent everything TinyGo has already invented, I can just take the LLVM IR TinyGo produces and run with that. Plus I also get most of the optimizations of LLVM to boot!

# License

MIT

The `html` package has a heavily modified version of the Go compiler's ssa.html dumper, and is under a [BSD style license](./html/LICENSE).

Some tests in the test suite were borrowed from other projects, see the header of those files for the license.
