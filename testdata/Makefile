TESTS = 000_simple.ll 001_add.ll 002_fib.ll 003_global.ll 004_ptrtovar.ll 005_struct.ll \
	006_structptr.ll 007_recurptr.ll 008_string.ll

all: $(TESTS)

%.ll: %.c
	clang --target=msp430-unknown-elf -march=msp430 -nostdlib -emit-llvm -S $<

clean:
	rm -rf $(TESTS)