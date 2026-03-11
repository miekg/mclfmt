package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/purpleidea/mgmt/lang/parser"
)

func TestPrettyPrint(t *testing.T) {
	dir := "testdata"
	testFiles, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("could not read %s: %q", dir, err)
	}
	for _, f := range testFiles {
		if f.IsDir() {
			continue
		}

		if filepath.Ext(f.Name()) != ".mcl" {
			continue
		}
		buf, err := os.ReadFile(dir + "/" + f.Name())
		if err != nil {
			t.Fatal(err)
		}
		pretty, err := os.ReadFile(dir + "/" + f.Name() + ".pretty")
		if err != nil {
			continue
		}
		t.Run(f.Name(), func(t *testing.T) {
			prog, err := parser.LexParse(bytes.NewReader(buf))
			if err != nil {
				t.Fatal(err)
			}
			lw := &LineWriter{Indent: 0, Start: true, b: &bytes.Buffer{}}
			Print(prog, lw, Option{})

			tr := trimSpace(lw.Bytes())
			tp := trimSpace(pretty)

			if string(tr) != string(tp) {
				t.Errorf("pretty and input, don't match")
				t.Logf("Got\n%s\n", tr)
				t.Logf("Want\n%s\n", tp)
			}
		})
	}
}
func trimSpace(buf []byte) []byte {
	buf = bytes.ReplaceAll(buf, []byte{' '}, []byte{'_'})
	////	buf = bytes.ReplaceAll(buf, []byte{'\n'}, nil)
	//	buf = bytes.ReplaceAll(buf, []byte{'\t'}, nil)
	return buf
}
