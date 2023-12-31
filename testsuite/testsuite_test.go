package testsuite

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
	"testing"

	"github.com/rj45/llbrew/arch"
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
			// if entry.Name() == "011_loopinter.ll" {
			// 	continue
			// }
			filename := path.Join(testdata, entry.Name())

			resultfile := strings.Replace(filename, ".ll", ".txt", 1)
			buf, err := os.ReadFile(resultfile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			result := string(bytes.TrimSpace(buf))

			testcases = append(testcases, testCase{
				name:     entry.Name(),
				filename: filename,
				result:   result,
			})
		}
	}

	arch.SetArch("rj32")

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

			outbuf := &bytes.Buffer{}

			c := compile.Compiler{
				OptSize:  0,
				OptSpeed: 0,
				Run:      true,
				RunWR:    outbuf,
			}

			err := c.Compile(tc.filename)

			// translate \r\n escape codes into \n
			buf := bytes.ReplaceAll(outbuf.Bytes(), []byte("\r\x1bD"), []byte("\n"))

			resultstr := string(bytes.TrimSpace(buf))

			if err != nil {
				if e, ok := err.(*exec.ExitError); ok && resultstr == "" {
					if fmt.Sprintf("%d", e.ExitCode()) != tc.result {
						t.Errorf("Expected run to exit with %s but got %d", tc.result, e.ExitCode())
					}
				} else {
					t.Error(err)
				}
			}
			if err == nil && resultstr == "" && tc.result != "0" {
				t.Error("expecting a non-zero result!")
			}

			if resultstr != "" && resultstr != tc.result {
				t.Errorf("Outputs did not match! expecting <<<<\n%q\n>>>> to match <<<<\n%q\n>>>>", resultstr, tc.result)
			}
		})
	}
}
