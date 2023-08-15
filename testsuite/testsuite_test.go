package testsuite

import (
	"fmt"
	"os"
	"path"
	"sort"
	"testing"

	"github.com/rj45/llbrew/compile"
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

func TestOptimized(t *testing.T) {
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c := compile.Compiler{
				OptSize:  1,
				OptSpeed: 1,
			}

			err := c.Compile(tc.filename)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUnoptimized(t *testing.T) {
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c := compile.Compiler{
				OptSize:  0,
				OptSpeed: 0,
			}

			err := c.Compile(tc.filename)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
