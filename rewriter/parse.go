package rewriter

import (
	"bufio"
	"fmt"
	"go/token"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/rj45/llbrew/arch/rj32"
	"github.com/rj45/llbrew/ir/op"
	"github.com/rj45/llbrew/ir/typ"
)

var types = &typ.Types{}

func parse(filename, pkg string) ([]rule, error) {
	var rules []rule

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	rulestr := ""
	ruleline := 0
	lineno := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lineno++

		if i := strings.Index(line, "//"); i >= 0 {
			// skip comments
			line = line[:i]
		}

		rulestr += strings.TrimSpace(rulestr + " " + line)

		if rulestr == "" {
			continue
		}

		if ruleline == 0 {
			ruleline = lineno
		}

		if !strings.Contains(rulestr, "=>") {
			continue
		}
		if strings.HasSuffix(rulestr, "=>") {
			continue
		}

		if parenBalance(rulestr) != 0 {
			continue
		}

		rules = append(rules, rule{loc: fmt.Sprintf("%s:%d", filename, ruleline), pkg: pkg})
		parseRule(rulestr, &rules[len(rules)-1])

		ruleline = 0
		rulestr = ""
	}

	return rules, nil
}

func parenBalance(str string) int {
	balance := 0
	for _, s := range str {
		if s == '(' {
			balance++
		} else if s == ')' {
			balance--
		}
	}
	return balance
}

func parseRule(str string, rule *rule) {
	parts := strings.Split(str, "=>")
	match := parts[0]

	if i := strings.Index(match, "&&"); i >= 0 {
		rule.cond = exprf(match[i+2:])
		match = match[:i]
	}
	rule.match = rule.parseInstr(normalizeSpaces(match))
	rule.repl = rule.parseValue(normalizeSpaces(parts[1]))
}

func normalizeSpaces(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

func splitSexpr(str string) []string {
	var res []string
outer:
	for str != "" {
		depth := 0           // bracket depths
		var open, close byte // open/close bracket char
		hasval := false      // seen a value?
		for i := 0; i < len(str); i++ {
			if depth == 0 {
				switch {
				case str[i] == '(':
					open, close = '(', ')'
					depth++
				case str[i] == '<':
					open, close = '<', '>'
					depth++
				case str[i] == '[':
					open, close = '[', ']'
					depth++
				case str[i] == '{':
					open, close = '{', '}'
					depth++
				case str[i] == ' ':
					if hasval {
						res = append(res, strings.TrimSpace(str[:i]))
						str = str[i:]
						continue outer
					}
				default:
					hasval = true
				}
			} else {
				hasval = true
				switch {
				case str[i] == open:
					depth++
				case str[i] == close:
					depth--
					if depth == 0 {
						res = append(res, strings.TrimSpace(str[:i+1]))
						str = str[i+1:]
						continue outer
					}
				}
			}
		}

		if hasval {
			res = append(res, strings.TrimSpace(str))
		}
		break
	}

	return res
}

func (r *rule) parseValue(str string) value {
	if str == "_" {
		return underscore
	}

	if str == "..." {
		return elipsis
	}

	if token.IsIdentifier(str) {
		if v, found := r.vars[str]; found {
			v.count++
			return v
		}
		if r.vars == nil {
			r.vars = make(map[string]*variable)
		}
		v := &variable{
			name:  str,
			count: 1,
		}
		r.vars[str] = v
		return v
	}

	if str[0] != '(' {
		_, err := strconv.Atoi(str)
		if err == nil || strings.HasPrefix(str, "reg.") {
			return constant(str)
		}
	}

	return r.parseInstr(str)
}

func (r *rule) parseInstr(str string) *instruction {
	instr := &instruction{}

	if str[0] != '(' {
		log.Fatalf("%s: expecting instr, got %q", r.loc, str)
	}

	parts := splitSexpr(str[1 : len(str)-1])

	if parts[0][0] == '[' {
		d := parts[0]
		parts = parts[1:]
		defs := splitSexpr(d[1 : len(d)-1])
		instr.defs = make([]value, len(defs))
		for i, def := range defs {
			instr.defs[i] = r.parseValue(def)
		}
	}

	instr.opstr = parts[0]
	parts = parts[1:]

	opparts := strings.Split(instr.opstr, ".")
	var err error
	switch opparts[0] {
	case "op":
		instr.op, err = op.OpString(opparts[1])
	case "rj32":
		instr.op, err = rj32.OpcodeString(opparts[1])
	}
	if err != nil {
		log.Fatalf("%s: could not parse op %s: %v", r.loc, instr.opstr, err)
	}
	if opparts[0] == r.pkg {
		instr.opstr = opparts[1]
	}

	if parts[0][0] == '<' {
		t := parts[0]
		typstr := t[1 : len(t)-1]
		parts = parts[1:]
		instr.typ, err = types.Parse(typstr)
		if err != nil {
			log.Fatalf("could not parse type %s: %v", typstr, err)
		}
	}

	instr.args = make([]value, len(parts))
	for i, str := range parts {
		instr.args[i] = r.parseValue(str)
	}

	return instr
}
