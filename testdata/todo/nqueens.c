// Copyright (c) 2018 Andrew Chambers
// MIT Licensed
// From github.com/c-testsuite/c-testsuite

int N;
int t[64];

int chk(int x, int y) {
    int i;
    int r;

    for (r=i=0; i<8; i++) {
        r = r + t[x + 8*i];
        r = r + t[i + 8*y];
        if (x+i < 8 & y+i < 8)
            r = r + t[x+i + 8*(y+i)];
        if (x+i < 8 & y-i >= 0)
            r = r + t[x+i + 8*(y-i)];
        if (x-i >= 0 & y+i < 8)
            r = r + t[x-i + 8*(y+i)];
        if (x-i >= 0 & y-i >= 0)
            r = r + t[x-i + 8*(y-i)];
    }
    return r;
}

int go(int n, int x, int y) {
    if (n == 8) {
        N++;
        return 0;
    }
    for (; y<8; y++) {
        for (; x<8; x++) {
            if (chk(x, y) == 0) {
                t[x + 8*y]++;
                go(n+1, x, y);
                t[x + 8*y]--;
            }
        }
        x = 0;
    }
    return 0;
}

int main() {
    go(0, 0, 0);
    return N;
}