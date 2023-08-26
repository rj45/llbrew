// Copyright (c) 2018 Andrew Chambers
// MIT Licensed
// From github.com/c-testsuite/c-testsuite

int fourtytwo() {
	return 42;
}

struct S {
	int (*fourtytwofn)();
} s = { &fourtytwo };

struct S * anon() {
	return &s;
}

typedef struct S * (*fty)();

fty go() {
	return &anon;
}

int main() {
	return go()()->fourtytwofn();
}