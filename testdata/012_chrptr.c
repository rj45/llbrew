// Copyright (c) 2018 Andrew Chambers
// MIT Licensed
// From github.com/c-testsuite/c-testsuite

int f1(char *p) {
	return *p+1;
}

int main() {
	char s = 1;

	return f1(&s);
}