# llir2asm

A LLVM IR to assembly compiler. In other words, an LLVM backend that produces assembly for hobby CPUs.

Input LLVM IR compiled using gclang, tinygo, etc for a CPU similar to yours (at least needs the same word size and pointer size), and it will output [customasm](https://github.com/hlorenzi/customasm) compatible assembly. Yes that's right, C or Go or Rust for your CPU you designed and built yourself.

At least that's the goal anyway, currently this is very much a work in progress with a long way to go.

If you want to help, please get in touch.

# License

MIT
