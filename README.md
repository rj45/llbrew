# llbrew

An LLVM backend designed specifically to make it easy to retarget to new instruction set architectures / CPUs.

Input LLVM IR compiled using gclang, tinygo, etc for a CPU similar to yours (at least needs the same word size and pointer size), and it will output [customasm](https://github.com/hlorenzi/customasm) compatible assembly. Yes that's right, C or Go or Rust for your CPU you designed and built yourself.

At least that's the goal anyway, currently this is very much a work in progress with a long way to go.

If you want to help, please get in touch.

# License

MIT

The `html` package has a heavily modified version of the Go compiler's ssa.html dumper, and is under a [BSD style license](./html/LICENSE).

Some tests in the test suite were borrowed from other projects, see the header of those files for the license.
