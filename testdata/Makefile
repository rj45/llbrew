TESTS = simple.ll add.ll

all: $(TESTS)

%.ll: %.c
	clang --target=riscv32-unknown-elf -march=rv32i -ffreestanding -fno-builtin -nostdlib -emit-llvm -S $<

