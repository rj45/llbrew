void putc(char c) {
    *((unsigned*)0xFF00) = (unsigned)c;
}

int main() {
    putc('H');
    putc('e');
    putc('l');
    putc('l');
    putc('o');
    putc('r');
    putc('l');
    putc('d');
    putc('!');
    putc('\n');
}