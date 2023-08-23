package rewriter

import (
	"bytes"
	"cmp"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"slices"
)

func Generate(name, rulesfile, outfile, pkg string) {
	rules, err := parse(rulesfile)
	if err != nil {
		panic(err)
	}

	fn, err := gen(name, rules)
	if err != nil {
		log.Fatalf("%s: %s", rulesfile, err)
	}

	out, err := os.Create(outfile)
	if err != nil {
		log.Fatalf("%s: failed to open output file %s", rulesfile, outfile)
	}
	defer out.Close()

	buf := &bytes.Buffer{}

	fmt.Fprintf(buf, "// Code generated from %s using 'go generate'; DO NOT EDIT.\n", rulesfile)
	fmt.Fprintf(buf, "\npackage %s\n\n", pkg)
	fmt.Fprintf(buf, "import (\n")
	fmt.Fprintf(buf, "\t%q\n", "github.com/rj45/llbrew/ir")
	fmt.Fprintf(buf, "\t%q\n", "github.com/rj45/llbrew/ir/op")
	fmt.Fprintf(buf, "\t%q\n", "github.com/rj45/llbrew/ir/reg")
	fmt.Fprintf(buf, ")\n\n")

	format.Node(buf, token.NewFileSet(), fn)

	src, err := format.Source(buf.Bytes())
	if err != nil {
		out.Write(buf.Bytes())
		log.Fatalf("%s: failed to format output for %s: %s", rulesfile, outfile, err)
	}
	_, err = out.Write(src)
	if err != nil {
		log.Fatalf("%s: failed to write output to %s: %s", rulesfile, outfile, err)
	}
}

func gen(name string, rules []rule) (*ast.FuncDecl, error) {
	var order []string
	oprules := map[string][]*rule{}
	for i := range rules {
		rule := &rules[i]
		op := rule.match.opstr
		if _, found := oprules[op]; !found {
			order = append(order, op)
		}
		oprules[op] = append(oprules[op], rule)
	}

	stmts := []ast.Stmt{}

	stmts = append(stmts, declf("instr", "it.Instr()"))
	stmts = append(stmts, declf("blk", "instr.Block()"))
	stmts = append(stmts, declf("fn", "blk.Func()"))
	stmts = append(stmts, declf("types", "fn.Types()"))
	stmts = append(stmts, assignf("_", "types"))
	stmts = append(stmts, assignf("_", "reg.SP"))

	expr := exprf("instr.Op")
	swtch := &ast.SwitchStmt{Tag: expr, Body: &ast.BlockStmt{}}
	cases := swtch.Body

	stmts = append(stmts, swtch)

	for _, op := range order {
		opexpr, err := parser.ParseExpr(op)
		if err != nil {
			return nil, fmt.Errorf("failed to parse op expr: %s", op)
		}
		theCase := &ast.CaseClause{List: []ast.Expr{opexpr}}
		cases.List = append(cases.List, theCase)

		for _, rule := range oprules[op] {
			theCase.Body = append(theCase.Body, &ast.ForStmt{
				Body: &ast.BlockStmt{List: rule.gen()},
			})
		}
	}

	fn := &ast.FuncDecl{
		Name: ast.NewIdent(name),
		Type: &ast.FuncType{Params: &ast.FieldList{List: []*ast.Field{{
			Names: []*ast.Ident{ast.NewIdent("it")},
			Type:  exprf("ir.Iter"),
		}}}},
		Body: &ast.BlockStmt{
			List: stmts,
		},
	}

	return fn, nil
}

func (rule *rule) gen() []ast.Stmt {
	var stmts []ast.Stmt

	var idents []*ast.Ident

	for _, v := range rule.vars {
		idents = append(idents, &ast.Ident{
			Name: v.name,
		})
	}

	slices.SortFunc(idents, func(a, b *ast.Ident) int {
		return cmp.Compare(a.Name, b.Name)
	})

	if len(idents) > 0 {
		stmts = append(stmts, &ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: idents,
						Type:  exprf("*ir.Value"),
					},
				},
			},
		})
	}

	stmts = rule.genMatch("instr", rule.match, stmts)

	if rule.cond != nil {
		stmt := breakf("!(a)")
		stmt.(*ast.IfStmt).Cond.(*ast.UnaryExpr).
			X.(*ast.ParenExpr).X = rule.cond
		stmts = append(stmts, stmt)
	}

	stmts = rule.genRepl(rule.repl, stmts)

	stmts = append(stmts, &ast.BranchStmt{Tok: token.BREAK})

	return stmts
}

func (r *rule) genMatch(vname string, instr *instruction, stmts []ast.Stmt) []ast.Stmt {
	if len(instr.defs) > 1 {
		stmts = append(stmts, breakf("%s.NumDefs() != %d", vname, len(instr.args)))
	}

	for i, def := range instr.defs {
		if def == underscore {
			continue
		}

		if v, ok := def.(*variable); ok {
			stmts = append(stmts, assignf(v.name, "%s.Def(%d)", vname, i))
			continue
		}

		if v, ok := def.(*instruction); ok {
			name := fmt.Sprintf("%s_def%d", vname, i)

			stmts = append(stmts, declf(name, "%s.Def(%d).Def().Instr()", vname, i))

			stmts = r.genMatch(name, v, stmts)
			continue
		}

		panic("unimplemented")
	}

	if len(instr.args) > 0 && instr.args[0] != elipsis {
		stmts = append(stmts, breakf("%s.NumArgs() != %d", vname, len(instr.args)))
	}

	commutative := instr.op.IsCommutative()
	if commutative {
		if instr.args[0] == instr.args[1] {
			// if the args are identical they are interchanable
			// don't generate commutative swap loop
			commutative = false
		}
		// todo: if args[0] && args[1] are identifiers and
		// only used once then commutative loop can be skipped
	}
	_ = commutative
	argstmts := stmts
	// var loop *ast.ForStmt

	// if commutative {
	// 	loop = commutativeLoop()
	// 	stmts = append(stmts, loop)
	// 	argstmts = loop.Body.List
	// }

	for i, arg := range instr.args {
		if arg == underscore || arg == elipsis {
			continue
		}

		if v, ok := arg.(*variable); ok {
			argstmts = append(argstmts, assignf(v.name, "%s.Arg(%d)", vname, i))
			continue
		}

		if v, ok := arg.(constant); ok {
			c := string(v)

			argstmts = append(argstmts, breakf("!%s.Arg(%d).HasConstValue(%s)", vname, i, c))

			continue
		}

		if v, ok := arg.(*instruction); ok {
			name := fmt.Sprintf("%s_arg%d", vname, i)

			argstmts = append(argstmts, declf(name, "%s.Arg(%d).Def().Instr()", vname, i))
			argstmts = append(argstmts, breakf("%s.Op != %s", name, v.opstr))

			argstmts = r.genMatch(name, v, argstmts)
			continue
		}

		panic("unimplemented")
	}

	// if !commutative {
	stmts = argstmts
	// } else {
	// 	loop.Body.List = argstmts
	// }

	return stmts
}

// func commutativeLoop() *ast.ForStmt {
// 	return nil
// }

func (r *rule) genRepl(val value, stmts []ast.Stmt) []ast.Stmt {
	switch v := val.(type) {
	case *variable:
		stmts = append(stmts, stmtf("it.ReplaceWith(%s)", v.name))
	case *instruction:
		if len(v.args) == 1 && v.args[0] == elipsis {
			stmts = append(stmts, stmtf("instr.Op = %s", v.opstr))
			if v.typ != nil {
				stmts = append(stmts, stmtf("instr.Def(0).Type = %s", v.typ.GoString()))
			}
			stmts = append(stmts, stmtf("it.Changed()"))
			break
		}

		args := ""
		cse := make(map[string]string)
		for _, arg := range v.args {
			args += ", "
			var s string
			stmts, s = r.genSubRepl(arg, stmts, cse)
			args += s
		}

		t := "instr.Type()"
		if v.typ != nil {
			t = v.typ.GoString()
		}
		stmts = append(stmts, stmtf("it.Update(%s, %s%s)", v.opstr, t, args))

	default:
		panic("unimplemented")
	}

	return stmts
}

func (r *rule) genSubRepl(val value, stmts []ast.Stmt, cse map[string]string) ([]ast.Stmt, string) {
	strval := val.String()

	if prev, ok := cse[strval]; ok {
		return stmts, prev
	}

	var vname string

	switch v := val.(type) {
	case *variable:
		vname = v.name
	case constant:
		vname = string(v)
	case *instruction:
		args := ""

		firstarg := ""

		for i, arg := range v.args {
			args += ", "
			var s string
			stmts, s = r.genSubRepl(arg, stmts, cse)
			args += s
			if i == 0 {
				firstarg = s
			}
		}

		t := ""
		if v.typ != nil {
			t = v.typ.GoString()
		} else if firstarg != "" {
			t = firstarg + ".Type"
		} else {
			log.Fatalf("%s: type required for sub expressions", r.loc)
		}

		if v.op.IsSink() {
			log.Fatalf("%s: sub expressions can't be sinks", r.loc)
		}

		iname := fmt.Sprintf("i%d", r.instrs)
		r.instrs++
		stmts = append(stmts, declf(iname, "it.Insert(%s, %s%s)", v.opstr, t, args))

		if len(v.defs) == 0 {
			vname = "_"
		} else {
			vname = v.defs[0].(*variable).name
		}

		if vname == "_" {
			vname = fmt.Sprintf("v%d", r.values)
			r.values++
		}

		stmts = append(stmts, declf(vname, "%s.Def(0)", iname))

		for i := 1; i < len(v.defs); i++ {
			n := v.defs[i].(*variable).name
			if n == "_" {
				continue
			}
			stmts = append(stmts, declf(n, "%s.Def(%d)", iname, i))
		}

	default:
		panic("unimplemented")
	}

	cse[strval] = vname

	return stmts, vname
}

func exprf(format string, a ...interface{}) ast.Expr {
	src := fmt.Sprintf(format, a...)
	expr, err := parser.ParseExpr(src)
	if err != nil {
		log.Fatalf("failed to parse expr %q: %s", src, err)
	}
	return expr
}

func stmtf(format string, a ...interface{}) ast.Stmt {
	src := fmt.Sprintf(format, a...)
	fsrc := "package p\nfunc _() {\n" + src + "\n}\n"
	file, err := parser.ParseFile(token.NewFileSet(), "", fsrc, 0)
	if err != nil {
		log.Fatalf("stmt parse error on %q: %v", src, err)
	}
	return file.Decls[0].(*ast.FuncDecl).Body.List[0]
}

func breakf(format string, a ...interface{}) ast.Stmt {
	expr := exprf(format, a...)
	return &ast.IfStmt{
		Cond: expr,
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.BranchStmt{Tok: token.BREAK},
			},
		},
	}
}

func declf(name, format string, a ...interface{}) ast.Stmt {
	lhs := exprf(name)
	rhs := exprf(format, a...)
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{rhs},
	}
}

func assignf(name, format string, a ...interface{}) ast.Stmt {
	lhs := exprf(name)
	rhs := exprf(format, a...)
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{rhs},
	}
}
