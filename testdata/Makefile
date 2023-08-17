TESTFILES = $(shell ls *.c)
TESTS = $(TESTFILES:%.c=%.ll)

all: $(TESTS)

%.ll: %.c
	clang --target=msp430-unknown-elf -march=msp430 -nostdlib -emit-llvm -S $<

clean:
	rm -rf $(TESTS)