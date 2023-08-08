package testsuite

import (
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"testing"

	"github.com/rj45/llir2asm/asm"
	"github.com/rj45/llir2asm/compile"
)

type testCase struct {
	name     string
	filename string
	result   string
}

var testcases []testCase

func TestMain(m *testing.M) {
	testdata := path.Join("..", "testdata")
	entries, err := os.ReadDir(testdata)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if path.Ext(entry.Name()) == ".ll" {
			filename := path.Join(testdata, entry.Name())

			testcases = append(testcases, testCase{
				name:     entry.Name(),
				filename: filename,
				result:   "42",
			})
		}
	}

	os.Exit(m.Run())
}

func Test(t *testing.T) {
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c := compile.Compiler{
				OptSize:  1,
				OptSpeed: 1,
				DumpIR:   nopWriteCloser{},
				DumpLL:   nopWriteCloser{},
			}

			prog, err := c.Compile(tc.filename)
			if err != nil {
				t.Fatal(err)
			}

			asm.Emit(io.Discard, asm.CustomASM{}, prog)
		})
	}
}

type nopWriteCloser struct{}

func (nopWriteCloser) Close() error {
	return nil
}

func (n nopWriteCloser) Write(p []byte) (int, error) {
	return len(p), nil
}
