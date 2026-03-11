package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/purpleidea/mgmt/lang/ast"
	"github.com/purpleidea/mgmt/lang/funcs/operators"
	"github.com/purpleidea/mgmt/lang/interfaces"
)

func Print(a any, w *LineWriter, opt Option) {
	switch a := a.(type) {
	case *ast.StmtProg:
		StmtProg(a, w, opt)
	case *ast.StmtRes:
		StmtRes(a, w, opt)
	case *ast.StmtBind:
		StmtBind(a, w, opt)
	case *ast.StmtResField:
		StmtResField(a, w, opt)
	case *ast.StmtEdge:
		StmtEdge(a, w, opt)
	case *ast.StmtEdgeHalf:
		StmtEdgeHalf(a, w, opt)
	case *ast.StmtClass:
		StmtClass(a, w, opt)
	case *ast.StmtInclude:
		StmtInclude(a, w, opt)
	case *ast.StmtImport:
		StmtImport(a, w, opt)
	case *ast.StmtFunc:
		StmtFunc(a, w, opt)
	case *ast.StmtResMeta:
		StmtResMeta(a, w, opt)
	case *ast.StmtResEdge:
		StmtResEdge(a, w, opt)

	case *ast.ExprStruct:
		ExprStruct(a, w, opt)
	case *ast.ExprStructField:
		ExprStructField(a, w, opt)
	case *ast.ExprFunc:
		ExprFunc(a, w, opt)
	case *ast.ExprStr:
		ExprStr(a, w, opt)
	case *ast.ExprInt:
		ExprInt(a, w, opt)
	case *ast.ExprFloat:
		ExprFloat(a, w, opt)
	case *ast.ExprBool:
		ExprBool(a, w, opt)
	case *ast.ExprVar:
		ExprVar(a, w, opt)
	case *ast.ExprCall:
		ExprCall(a, w, opt)
	case *ast.ExprMap:
		ExprMap(a, w, opt)
	case *ast.ExprMapKV:
		ExprMapKV(a, w, opt)
	case *ast.ExprList:
		ExprList(a, w, opt)
	case *ast.ExprIf:
		ExprIf(a, w, opt)

	case *interfaces.Arg:
		InterfacesArg(a, w, opt)

	default:
		panic("mclfmt: unhandled ast " + fmt.Sprintf("%T", a) + fmt.Sprintf(" : %v", a))
	}
}

func StmtProg(a *ast.StmtProg, w *LineWriter, opt Option) {
	for _, b := range a.Body {
		Print(b, w, opt)
	}
}

func StmtBind(a *ast.StmtBind, w *LineWriter, opt Option) {
	if opt.DropSpace {
		fmt.Fprintf(w, "$%s =", a.Ident)
		opt.DropSpace = false
	} else {
		fmt.Fprintf(w, " $%s =", a.Ident)
	}
	Print(a.Value, w, opt)
	io.WriteString(w, "\n") // always closes the line
}

func StmtRes(a *ast.StmtRes, w *LineWriter, opt Option) {
	if opt.DropSpace {
		fmt.Fprintf(w, "%s", a.Kind)
		opt.DropSpace = false
	} else {
		fmt.Fprintf(w, " %s", a.Kind)
	}

	Print(a.Name, w, opt)

	io.WriteString(w, " {") // with no contents this stays on a single line
	w.Indent++
	for i, c := range a.Contents {
		if i == 0 {
			io.WriteString(w, "\n")
		}
		Print(c, w, opt)
		io.WriteString(w, ",\n")
	}
	w.Indent--
	io.WriteString(w, "}\n")
}

func StmtEdge(a *ast.StmtEdge, w *LineWriter, opt Option) {
	// Start with new line?
	fmt.Fprintln(w)
	for i, e := range a.EdgeHalfList {
		Print(e, w, opt)
		if i < len(a.EdgeHalfList)-1 {
			fmt.Fprintf(w, " %s ", "->")
		}
	}
}

func StmtEdgeHalf(a *ast.StmtEdgeHalf, w *LineWriter, opt Option) {
	fmt.Fprintf(w, "%s%s[", strings.ToUpper(a.Kind[:1]), a.Kind[1:])
	Print(a.Name, w, opt)
	if a.SendRecv != "" {
		fmt.Fprintf(w, "].%s", a.SendRecv)
	} else {
		fmt.Fprint(w, "]")
	}
}

func ExprStr(a *ast.ExprStr, w *LineWriter, opt Option) {
	if opt.DropQuote {
		if opt.DropSpace {
			fmt.Fprintf(w, "%s", a.V)
		} else {
			fmt.Fprintf(w, " %s", a.V)
		}
		return
	}

	if opt.DropSpace {
		fmt.Fprintf(w, `"%s"`, a.V)
	} else {
		fmt.Fprintf(w, ` "%s"`, a.V)
	}
}

func ExprInt(a *ast.ExprInt, w *LineWriter, opt Option) {
	if opt.DropSpace {
		fmt.Fprintf(w, "%d", a.V)
	} else {
		fmt.Fprintf(w, " %d", a.V)
	}
}

func ExprFloat(a *ast.ExprFloat, w *LineWriter, opt Option) {
	if opt.DropSpace {
		fmt.Fprintf(w, "%g", a.V)
	} else {
		fmt.Fprintf(w, " %g", a.V)
	}
}

func ExprBool(a *ast.ExprBool, w *LineWriter, opt Option) {
	if opt.DropSpace {
		fmt.Fprintf(w, "%t", a.V)
	} else {
		fmt.Fprintf(w, " %t", a.V)
	}
}

func ExprVar(a *ast.ExprVar, w *LineWriter, opt Option) {
	if opt.DropSpace {
		fmt.Fprintf(w, "$%s", a.Name)
	} else {
		fmt.Fprintf(w, " $%s", a.Name)
	}
}

func StmtResField(a *ast.StmtResField, w *LineWriter, opt Option) {
	if opt.DropSpace {
		fmt.Fprintf(w, "%s =>", a.Field)
	} else {
		fmt.Fprintf(w, " %s =>", a.Field)
	}

	Print(a.Value, w, opt)
}

func StmtClass(a *ast.StmtClass, w *LineWriter, opt Option) {
	if opt.DropSpace {
		fmt.Fprintf(w, "class %s", a.Name)
	} else {
		fmt.Fprintf(w, " class %s", a.Name)
	}

	opt.DropSpace = false
	if a.Args != nil {
		io.WriteString(w, "(")
		for i, arg := range a.Args {
			if i == 0 {
				opt.DropSpace = true
			} else {
				opt.DropSpace = false
			}
			Print(arg, w, opt)
			if i < len(a.Args)-1 {
				fmt.Fprint(w, ",")
			}
		}
		fmt.Fprint(w, ")")
		opt.DropSpace = false

	}
	io.WriteString(w, " {\n")

	w.Indent++
	// always StmtProg?
	prog := a.Body.(*ast.StmtProg)
	for _, b := range prog.Body {
		Print(b, w, opt)
	}
	w.Indent--
	io.WriteString(w, "}\n")
}

func StmtInclude(a *ast.StmtInclude, w *LineWriter, opt Option) {
	fmt.Fprintf(w, "include %s", a.Name)
	opt.DropSpace = false
	if a.Args != nil {
		io.WriteString(w, "(")
		for i, arg := range a.Args {
			if i == 0 {
				opt.DropSpace = true
			} else {
				opt.DropSpace = false
			}
			Print(arg, w, opt)
			if i < len(a.Args)-1 {
				fmt.Fprint(w, ",")
			}
		}
		fmt.Fprint(w, ")")
		opt.DropSpace = false
	}

	if a.Alias != "" {
		fmt.Fprintf(w, " as %s", a.Alias)
	}

	io.WriteString(w, "\n")
}

func StmtImport(a *ast.StmtImport, w *LineWriter, opt Option) {
	fmt.Fprintf(w, "import \"%s\"", a.Name)
	if a.Alias != "" {
		fmt.Fprintf(w, " as %s", a.Alias)
	}
	io.WriteString(w, "\n")
}

func StmtFunc(a *ast.StmtFunc, w *LineWriter, opt Option) {
	fmt.Fprintf(w, "func %s", a.Name)
	opt.Func = a.Name
	Print(a.Func, w, opt)
}

func StmtResMeta(a *ast.StmtResMeta, w *LineWriter, opt Option) {
	if a.Property != "meta" {
		fmt.Fprintf(w, "Meta:%s =>", a.Property)
	} else {
		fmt.Fprint(w, "Meta =>")
	}
	if a.Condition != nil {
		Print(a.Condition, w, opt)
		io.WriteString(w, " ?:")
	}

	Print(a.MetaExpr, w, opt)
}

func StmtResEdge(a *ast.StmtResEdge, w *LineWriter, opt Option) {
	fmt.Fprintf(w, "%s%s => ", strings.ToUpper(a.Property[:1]), a.Property[1:])
	opt.DropSpace = true
	Print(a.EdgeHalf, w, opt)
}

func ExprFunc(a *ast.ExprFunc, w *LineWriter, opt Option) {
	if opt.Func == "" { // No StmtFunc seen
		if opt.DropSpace {
			io.WriteString(w, "func")
		} else {
			io.WriteString(w, " func")
		}
	}
	opt.Func = ""
	opt.DropSpace = false
	io.WriteString(w, "(")

	for i, arg := range a.Args {
		if i == 0 {
			opt.DropSpace = true
		} else {
			opt.DropSpace = false
		}
		Print(arg, w, opt)
		if i < len(a.Args)-1 {
			fmt.Fprint(w, ",")
		}
	}
	fmt.Fprint(w, ")")
	opt.DropSpace = false

	if a.Return != nil {
		fmt.Fprintf(w, " %s", a.Return)
	}

	io.WriteString(w, " {\n")
	w.Indent++
	Print(a.Body, w, opt)
	w.Indent--
	io.WriteString(w, "\n") // TODO(miek): adds extra newline in nested functions
	io.WriteString(w, "}\n")
}

func ExprStruct(a *ast.ExprStruct, w *LineWriter, opt Option) {
	if opt.DropSpace {
		io.WriteString(w, "struct{\n")
		opt.DropSpace = false
	} else {
		io.WriteString(w, " struct{\n")
	}

	w.Indent++
	for _, f := range a.Fields {
		Print(f, w, opt)
		io.WriteString(w, ",\n")
	}
	w.Indent--
	io.WriteString(w, "}")
}

func ExprStructField(a *ast.ExprStructField, w *LineWriter, opt Option) {
	fmt.Fprint(w, a.Name)
	fmt.Fprintf(w, " =>")
	Print(a.Value, w, opt)
}

func ExprCall(a *ast.ExprCall, w *LineWriter, opt Option) {
	switch a.Name {
	case operators.OperatorFuncName:
		Print(a.Args[1], w, opt)

		opt.DropQuote = true
		Print(a.Args[0], w, opt)
		opt.DropQuote = false

		Print(a.Args[2], w, opt)
	default:
		if opt.DropSpace {
			fmt.Fprintf(w, "%s(", a.Name)
		} else {
			fmt.Fprintf(w, " %s(", a.Name)
		}

		for i, arg := range a.Args {
			if i == 0 {
				opt.DropSpace = true
			} else {
				opt.DropSpace = false
			}
			Print(arg, w, opt)
			if i < len(a.Args)-1 {
				fmt.Fprint(w, ",")
			}
		}
		fmt.Fprint(w, ")")
		opt.DropSpace = false
	}
}

func ExprMap(a *ast.ExprMap, w *LineWriter, opt Option) {
	if opt.DropSpace {
		io.WriteString(w, "{\n")
		opt.DropSpace = false
	} else {
		io.WriteString(w, " {\n")
	}

	w.Indent++
	for _, kv := range a.KVs {
		Print(kv, w, opt)
		io.WriteString(w, ",\n")
	}
	w.Indent--
	io.WriteString(w, "}")
}

func ExprMapKV(a *ast.ExprMapKV, w *LineWriter, opt Option) {
	Print(a.Key, w, opt)
	fmt.Fprintf(w, " =>")
	Print(a.Val, w, opt)
}

func ExprList(a *ast.ExprList, w *LineWriter, opt Option) {
	if opt.DropSpace {
		io.WriteString(w, "[")
		opt.DropSpace = false
	} else {
		io.WriteString(w, " [")
	}
	for i, e := range a.Elements {
		if i == 0 {
			opt.DropSpace = true
		} else {
			opt.DropSpace = false
		}
		Print(e, w, opt)
		fmt.Fprint(w, ",")
	}
	opt.DropSpace = false
	io.WriteString(w, "]")
}

func ExprIf(a *ast.ExprIf, w *LineWriter, opt Option) {
	if opt.DropSpace {
		io.WriteString(w, "if")
		opt.DropSpace = false
	} else {
		io.WriteString(w, " if")
	}

	Print(a.Condition, w, opt)

	io.WriteString(w, " {\n")
	w.Indent++
	Print(a.ThenBranch, w, opt)
	w.Indent--
	io.WriteString(w, "\n")
	io.WriteString(w, "}")
	if a.ElseBranch == nil {
		return
	}

	io.WriteString(w, " else {\n")
	w.Indent++
	Print(a.ElseBranch, w, opt)
	w.Indent--
	io.WriteString(w, "\n")
	io.WriteString(w, "}")
}

func InterfacesArg(a *interfaces.Arg, w *LineWriter, opt Option) {
	if opt.DropSpace {
		fmt.Fprintf(w, "$%s", a.Name)
	} else {
		fmt.Fprintf(w, " $%s", a.Name)
	}

	if a.Type != nil {
		fmt.Fprintf(w, " %s", a.Type)
	}
}
