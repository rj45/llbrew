int fibonacci(int n) {
	if (n <= 1) {
		return n;
	}

    int n1 = 0;
	int n2 = 1;
	int n3 = 0;

	for (int i = 2; i < n; i++) {
		n3 = n1+n2;
		n1 = n2;
		n2 = n3;
	}

	return n2 + n1;
}

int main() {
	return fibonacci(7); // 13
}