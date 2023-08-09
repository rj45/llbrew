int fibonacci(int n) {
	if (n <= 1) {
		return n;
	}

	int n2 = 0;
    int n1 = 1;

	for (int i = 2; i < n; i++) {
        n2 = n1;
        n1 = n1+n2;
	}

	return n2 + n1;
}

int main() {
	return fibonacci(7); // 13
}