struct S { int x; int y; };

void setS(struct S *p) {
    p->x = 1;
    p->y = 2;
}

int main() {

	struct S s;

	setS(&s);
	return s.y + s.x - 3;
}
