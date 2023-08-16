int len(char *s) {
    int len;
    for (len = 0; ; len++) {
        if (s[len] == 0) {
            return len;
        }
    }
}

int main() {
	char *p;

	p = "hello world!";
	return len(p);
}