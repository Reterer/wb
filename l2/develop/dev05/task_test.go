package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestNewGrep(t *testing.T) {
	cfg := &config{}
	_, err := NewGrep(cfg)
	if err != nil {
		t.Errorf("err should be nil but %s", err)
	}
}

func getInput(lines []string) io.Reader {
	return bytes.NewBufferString(strings.Join(lines, "\n"))
}

func TestGrepio(t *testing.T) {
	testCases := []struct {
		testName string
		grep     Grep
		filename string
		input    []string
		want     []string
	}{
		{}, // empty case
		{ // empty pattern
			testName: "empty pattern",
			grep:     Grep{pattern: ""},
			input: []string{
				"",
				"hello",
				"something hello",
				"HELLO",
				"heh",
			},
			want: []string{
				"",
				"hello",
				"something hello",
				"HELLO",
				"heh",
			},
		},
		{ // find just lines
			testName: "find just lines",
			grep:     Grep{pattern: "hello"},
			input: []string{
				"",
				"hello",
				"something hello",
				"HELLO",
				"heh",
			},
			want: []string{
				"hello",
				"something hello",
			},
		},
		// IGNORE CASE
		{ // find just lines
			testName: "IGNORE CASE find just lines",
			grep: Grep{
				pattern: "hello",
				opts: grepOpts{
					ignoreCase: true,
				},
			},
			input: []string{
				"",
				"hello",
				"something hello",
				"HELLO",
				"heh",
			},
			want: []string{
				"hello",
				"something hello",
				"HELLO",
			},
		},
		// INVERT
		{ // find just lines
			testName: "INVERT find just lines",
			grep: Grep{
				pattern: "hello",
				opts: grepOpts{
					invert: true,
				},
			},
			input: []string{
				"",
				"hello",
				"something hello",
				"HELLO",
				"heh",
			},
			want: []string{
				"",
				"HELLO",
				"heh",
			},
		},
		// COUNT
		{ // just count
			testName: "COUNT just count",
			grep: Grep{
				pattern: "hello",
				opts: grepOpts{
					count: true,
				}},
			input: []string{
				"",
				"hello",
				"something hello",
				"HELLO",
				"heh",
			},
			want: []string{
				"2",
			},
		},
		{ // count with invert
			testName: "COUNT count with invert",
			grep: Grep{
				pattern: "hello",
				opts: grepOpts{
					count:  true,
					invert: true,
				}},
			input: []string{
				"",
				"hello",
				"something hello",
				"HELLO",
				"heh",
			},
			want: []string{
				"3",
			},
		},
		// FIXED
		{ // find just lines
			testName: "FIXED find just lines",
			grep: Grep{
				pattern: "hello",
				opts: grepOpts{
					fixed: true,
				}},
			input: []string{
				"",
				"hello",
				"something hello",
				"HELLO",
				"heh",
			},
			want: []string{
				"hello",
			},
		},
		{ // fixed with ignore case
			testName: "FIXED fixed with ignore case",
			grep: Grep{
				pattern: "hello",
				opts: grepOpts{
					fixed:      true,
					ignoreCase: true,
				}},
			input: []string{
				"",
				"hello",
				"something hello",
				"HELLO",
				"heh",
			},
			want: []string{
				"hello",
				"HELLO",
			},
		},
		// PRINT FORMAT
		{ // filename
			testName: "PRINT FORMAT filename",
			grep: Grep{
				pattern: "hello",
				opts:    grepOpts{}},
			filename: "test",
			input: []string{
				"",
				"hello",
				"something hello",
				"HELLO",
				"heh",
			},
			want: []string{
				"test:hello",
				"test:something hello",
			},
		},
		{ // line num
			testName: "PRINT FORMAT line num",
			grep: Grep{
				pattern: "hello",
				opts: grepOpts{
					lineNum: true,
				}},
			filename: "test",
			input: []string{
				"",
				"hello",
				"something hello",
				"HELLO",
				"heh",
			},
			want: []string{
				"test:2:hello",
				"test:3:something hello",
			},
		},
		// AFTER
		{ // AFTER find just lines
			testName: "AFTER find ",
			grep: Grep{
				pattern: "something hello",
				opts: grepOpts{
					after: 2,
				}},
			input: []string{
				"",
				"something hello", // 1
				"hello",           // |
				"something hello", // 2
				"123",             // |
				"heh",             // |
				"string",
				"abc",
				"abc2",
				"something hello", // 3
				"abc3",            // |
			},
			want: []string{
				"something hello", // 1
				"hello",           // |
				"something hello", // 2
				"123",             // |
				"heh",             // |
				"something hello", // 3
				"abc3",            // |
			},
		},
		{ // AFTER find with filename
			testName: "AFTER find with filename",
			grep: Grep{
				pattern: "something hello",
				opts: grepOpts{
					after: 2,
				}},
			filename: "test",
			input: []string{
				"",
				"hello",
				"something hello",
				"HELLO",
				"heh",
				"string",
			},
			want: []string{
				"test:something hello",
				"test-HELLO",
				"test-heh",
			},
		},
		// BEFORE
		{ // find with filename
			testName: "BEFORE find with filename",
			grep: Grep{
				pattern: "something hello",
				opts: grepOpts{
					before: 2,
				}},
			filename: "test",
			input: []string{
				"",                // |
				"something hello", // 1
				"hello",           // |
				"something hello", // 2
				"123",
				"heh",
				"string",
				"abc",             // |
				"abc2",            // |
				"something hello", // 3
				"abc3",
			},
			want: []string{
				"test-",                // |
				"test:something hello", // 1
				"test-hello",           // |
				"test:something hello", // 2
				"test-abc",             // |
				"test-abc2",            // |
				"test:something hello", // 3
			},
		},
	}

	for _, tc := range testCases {
		reader := getInput(tc.input)
		var buf bytes.Buffer
		err := tc.grep.grepio(tc.filename, reader, &buf)
		if err != nil {
			t.Errorf("err should be nil but %s", err)
		}

		got := buf.String()
		wants := strings.Join(tc.want, "\n")
		// fmt.Println(buf.Bytes())
		// fmt.Println([]byte(wants))
		if got != wants {
			t.Errorf("tc:'%s'\nGOT:\n%s\n\nWANT:\n%s", tc.testName, got, wants)
		}
	}

}
