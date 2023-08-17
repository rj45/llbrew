int add(int a, int b) {
    return a+b;
}

int sub(int a, int b) {
    return a-b;
}

int and(int a, int b) {
    return a & b;
}

int or(int a, int b) {
    return a | b;
}

int xor(int a, int b) {
    return a ^ b;
}

int main() {
    int a = add(123, 456); // 579
    int b = sub(654, 321); // 333
    a = xor(a, b); // 782
    b = or(a, b); // 847
    a = and(b, 481) ; // 321

    return a - 279; // 42
}