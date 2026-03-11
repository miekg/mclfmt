package main

import (
	"bytes"
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/purpleidea/mgmt/lang/parser"
)

func TestPrint(t *testing.T) {
	testcases := []struct {
		name   string
		code   string
		pretty string
	}{
		{
			name: "addition",
			code: `
		test "t1" {
			int64ptr => 13 + 42,
		}
		`,
		},
		{
			name: "multiple float addition",
			code: `
			test "t1" {
				float32 => -25.38789 + 32.6 + 13.7,
			}
			`,
		},
		{
			name: "func call dotted 1",
			code: `
			$x1 = pkg.foo1()
			`,
		},
		{
			name: "func call dotted 2",
			code: `
			$x1 = pkg.foo1(true, "hello")
			`,
		},
		{
			name: "edge stmt",
			code: `
			test "t1" {
				int64ptr => 42,
			}
			test "t2" {
				int64ptr => 13,
			}

			Test["t1"].foosend -> Test["t2"].barrecv
			`,
		},
		{
			name: "one map",
			code: `
			$somemap = {
				"foo" => "foo1",
				"bar" => "bar1",
			}
			`,
		},
		{
			name: "simple dotted class 1",
			code: `
			# a dotted identifier only occurs via an imported class
			class c1 {
				test "t1" {
					stringptr => "hello",
				}
			}
			# a dotted identifier is allowed here if it is imported
			include pkg.c1
			`,
		},
		{
			name: "simple import 1",
			code: `
			import "foo1"
			`,
		},
		{
			name: "simple import 2",
			code: `
			import "foo1" as bar
			`,
		},
		{
			name: "simple function stmt 1",
			code: `
			func f1() {
				42
			}
			`,
		},
		{
			name: "simple function stmt 3",
			code: `
			func f3($a int, $b) int {
				$a + $b
			}
			`,
		},
		{
			name: "iter.map",
			code: `
			import "iter"
			$fn = func($x) { $x * $x }
			$out = iter.map([1,2,3,], $fn)
			`,
			pretty: `import "iter"

			$fn = func($x) {
				$x * $x
			}
			$out = iter.map([1, 2, 3,], $fn)
			`,
		},
		{
			name: "simple nested function 1",
			code: `
			func funcgen() {	# returns a function expression
				func() {
					"hello"
				}
			}
			$fn = funcgen()
			$foo = $fn()	# hello
			`,
		},
		{
			name: "res meta stmt",
			code: `
			test "t1" {
				Meta:noop => true,
				Meta:delay => true ?: 42,
			}
			test "t2" {
				Meta:limit => 0.45,
				Meta:burst => 4,
			}
			test "t3" {
				Meta:noop => true, # meta params can be combined
				Meta => struct{
					poll => 5,
					retry => 3,
					sema => ["foo:1", "bar:3",],
				},
			}
			`,
		},
		{
			name: "res edge stmt",
			code: `
			tar "/tmp/foo.tar" {
				inputs => [
						"/tmp/tar/",
						"/tmp/standalone",
				],
				format => $const.res.tar.format.gnu,
				Depend => File["/tmp/tar/"], # TODO: add autoedges
			}
			`,
		},
	}

	re := regexp.MustCompile(`(?m)^[ \t]{3}`)
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			prog, err := parser.LexParse(strings.NewReader(tc.code))
			if err != nil {
				t.Fatal(err)
			}
			lw := &LineWriter{Indent: 0, Start: true, b: &bytes.Buffer{}}

			Print(prog, lw, Option{})

			// t.Log("*in*", re.ReplaceAllString(tc.code, ""))
			// t.Log("*out*\n", lw.String())

			if tc.pretty == "" {
				return
			}

			tcode := re.ReplaceAllString(tc.pretty, "")
			got := lw.String()
			diff := cmp.Diff(tcode, got)
			if diff != "" {
				t.Log("*in*", re.ReplaceAllString(tc.code, ""))
				t.Log("*out*\n", lw.String())

				t.Fatal(diff)
			}
		})
	}
}

func TestLineWriter(t *testing.T) {
	lw := &LineWriter{Indent: 0, Start: true, b: &bytes.Buffer{}}
	io.WriteString(lw, "bla")
	if lw.String() != "bla" {
		t.Fatalf("expected %s, got %s", "bla", lw.String())
	}
}
