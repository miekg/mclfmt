package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/purpleidea/mgmt/lang/parser"
)

func TestPrint(t *testing.T) {
	testcases := []struct {
		name string
		code string
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
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			prog, err := parser.LexParse(strings.NewReader(tc.code))
			if err != nil {
				t.Fatal(err)
			}
			lw := &LineWriter{Indent: 0, Start: true, b: &bytes.Buffer{}}

			Print(prog, lw, Option{})
			fmt.Println(lw.String())
		})
	}
}
