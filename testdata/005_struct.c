// Copyright (c) 2018 Andrew Chambers
// MIT Licensed
// From github.com/c-testsuite/c-testsuite

struct { int x; int y; } s;

void setvars() {
    s.x = 3;
	s.y = 5;
}

int main()
{
	setvars();
	return s.y - s.x - 2;
}