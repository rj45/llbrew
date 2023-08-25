
#define N 64
unsigned A[N][N];

void init() {
    for (unsigned i = 0; i < N; i++) {
        for (unsigned j = 0; j < N; j++) {
            A[i][j] = i + j;
        }
    }
}

unsigned y = 0;

unsigned interchange() {
    unsigned y;
    for (unsigned i = 0; i < N; i++) {
        y = 0;
        for (unsigned j = 0; j < N; j++) {
            A[i][j] += 1;
            y += A[i][j];
        }
    }
    return y;
}

int main () {
    init();
    return interchange() & 127;
}