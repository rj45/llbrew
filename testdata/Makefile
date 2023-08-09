TESTS = 000_simple.ll 001_add.ll 002_fib.ll

all: $(TESTS)

%.ll: %.c
	clang --target=msp430-unknown-elf -march=msp430 -ffreestanding -fno-builtin -nostdlib -emit-llvm -S $<

clean:
	rm -rf $(TESTS)