TESTS = simple.ll add.ll

all: $(TESTS)

%.ll: %.c
	clang --target=msp430-unknown-elf -march=msp430 -ffreestanding -fno-builtin -nostdlib -emit-llvm -S $<

clean:
	rm -rf $(TESTS)