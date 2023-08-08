// Copyright 2015 The Go Authors.
// Copyright 2022 Ryan "rj45" Sanche
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"fmt"
	"go/token"
	"html"
	"io"
	"strings"

	"github.com/rj45/llir2asm/ir"
	"github.com/rj45/llir2asm/ir/op"
)

// IRDecorator emits an HTML version of the IR
type IRDecorator struct {
	Asm bool
}

func (dec IRDecorator) Begin(out io.Writer, what interface{}) {
	switch v := what.(type) {
	case *ir.Block:
		fmt.Fprintf(out, "<ul class=\"%s ssa-print-func\">\n", v.String())

	case *ir.Instr:
		fmt.Fprint(out, "<li class=\"ssa-long-value\">\n")
		ids := ""
		for i := 0; i < v.NumDefs(); i++ {
			ids += " " + v.Def(i).IDString()
		}
		fmt.Fprintf(out, "<span class=\"%s ssa-long-value\">\n", ids)

		if v.Pos != token.NoPos {
			fmt.Fprintf(out, "<span class=\"l%v line-number\">(%d)</span>", v.LineNo(), v.LineNo())
		} else {
			fmt.Fprintf(out, "<span class=\"no-line-number\">(?)</span>")
		}
	}
}

func (dec IRDecorator) End(out io.Writer, what interface{}) {
	switch v := what.(type) {
	case *ir.Instr:
		fmt.Fprintf(out, "</span>\n")
		fmt.Fprint(out, "</li>\n")

	case *ir.Block:
		if v.NumInstrs() > 0 { // end list of values
			fmt.Fprint(out, "</ul>")
			fmt.Fprint(out, "</li>")
		}

		fmt.Fprint(out, "<li class=\"ssa-end-block\">")
		if v.NumSuccs() > 0 {
			fmt.Fprint(out, " &#8594;") // right arrow
			for i := 0; i < v.NumSuccs(); i++ {
				fmt.Fprint(out, " "+blockHTML(v.Succ(i)))
			}
		}
		fmt.Fprint(out, "</li>")

		fmt.Fprint(out, "</ul>")
	}
}

func (ss IRDecorator) BeginLabel(out io.Writer, what interface{}) {
	switch what.(type) {
	case *ir.Block:
		fmt.Fprintf(out, "<li class=\"ssa-start-block\">\n")

	}
}

func (ss IRDecorator) EndLabel(out io.Writer, what interface{}) {
	switch v := what.(type) {
	case *ir.Block:
		if v.NumPreds() > 0 {
			fmt.Fprint(out, " &#8592;") // left arrow
			for i := 0; i < v.NumPreds(); i++ {
				pred := v.Pred(i)
				fmt.Fprintf(out, " %s", blockHTML(pred))
			}
		}
		if v.NumInstrs() > 0 {
			fmt.Fprint(out, `<button onclick="hideBlock(this)">-</button>`)
		}
		fmt.Fprint(out, "</li>\n")

		if v.NumInstrs() > 0 { // start list of values
			fmt.Fprint(out, "<li class=\"ssa-value-list\">\n")
			fmt.Fprint(out, "<ul>\n")
		}
	}
}

func (dec IRDecorator) WrapLabel(str string, what interface{}) string {
	switch v := what.(type) {
	case *ir.Block:
		return blockHTML(v)
	case *ir.Value:
		return valueHTML(str, v)
	}
	return str
}

func (dec IRDecorator) WrapRef(str string, what interface{}) string {
	switch v := what.(type) {
	case *ir.Value:
		return valueHTML(str, v)
	case *ir.Block:
		return blockHTML(v)
	}
	return str
}

func (dec IRDecorator) WrapType(str string) string {
	return "<span class=\"ssa-instr-type\">" + html.EscapeString(str) + "</span>"
}

func (dec IRDecorator) WrapOp(str string, vop ir.Op) string {
	opstr := fmt.Sprintf("<span class=\"ssa-instr\">%s</span>", vop.String())
	if vop.IsCopy() {
		opstr = fmt.Sprintf("<span class=\"ssa-instr-copy\">%s</span>", vop.String())
	}
	if vop == op.Call {
		opstr = fmt.Sprintf("<span class=\"ssa-instr-call\">%s</span>", vop.String())
	}
	return opstr
}

func (dec IRDecorator) SSAForm() bool {
	return !dec.Asm
}

func valueHTML(str string, v *ir.Value) string {
	if v.IsConst() && v.Const().Kind() == ir.IntConst {
		return fmt.Sprintf("<span class=\"ssa-value-const-num\">%s</span>", str)
	}

	if v.IsConst() {
		return fmt.Sprintf("<span class=\"ssa-value-const\">%s</span>", str)
	}

	id := html.EscapeString(str)
	s := ""
	if strings.Contains(id, "_") {
		parts := strings.Split(id, "_")
		id = parts[0]
		s = strings.Join(parts[1:], "_")

		s = fmt.Sprintf("_<span class=\"%s ssa-value-reg\">%s</span>", s, s)
	}

	return fmt.Sprintf("<span class=\"%s ssa-value\"><span class=\"ssa-value-id\">%s</span>%s</span>", id, id, s)
}

func blockHTML(b *ir.Block) string {
	// TODO: Using the value ID as the class ignores the fact
	// that value IDs get recycled and that some values
	// are transmuted into other values.
	s := html.EscapeString(b.String())
	return fmt.Sprintf("<span class=\"%s ssa-block\">%s</span>", s, s)
}
