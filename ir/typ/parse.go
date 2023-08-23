package typ

import (
	"fmt"
	"strconv"
	"strings"
)

func (types *Types) Parse(str string) (Type, error) {
	switch str[0] {
	case 'i':
		bits, err := strconv.Atoi(str[1:])
		if err != nil {
			return nil, fmt.Errorf("error parsing integer type: %w", err)
		}
		return types.IntegerType(bits), nil
	case '*':
		elem, err := types.Parse(str[1:])
		if err != nil {
			return nil, fmt.Errorf("error parsing pointer type: %w", err)
		}
		return types.PointerType(elem, 0), nil
	case 's':
		if !strings.HasPrefix(str, "struct ") {
			return nil, fmt.Errorf("expected \"struct\" got %q", str)
		}
		panic("todo: parse structs")
	case 'f':
		// todo: add floats too
		if !strings.HasPrefix(str, "func") {
			return nil, fmt.Errorf("expected \"struct\" got %q", str)
		}
		panic("todo: parse funcs")
	case 'w':
		if str == "word" {
			return types.IntegerType(0), nil
		}
	case '[':
		panic("todo: parse arrays")
	}
	return nil, fmt.Errorf("could not parse type %q", str)
}
