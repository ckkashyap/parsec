package parsec

import (
	"strings"
)

type Something interface {}

type Parser  func (string) ( Something, string, bool )

func Result (a Something) Parser {
	var ret = func (str string) (Something, string, bool){
		return a, str, true
	}
	return ret
}

func Zero() Parser {
	var ret = func (str string) (Something, string, bool){
		return 0, str, false// 0? is that right
	}
	return ret
}


func Item() Parser {
	var ret = func (str string) (Something, string, bool) {
		r := strings.NewReader(str)
		c, s, _ := r.ReadRune()
		if s == 0 {return "", str, false}
		return c, str[s:] , true
	}
	return ret
}


func Bind(p Parser, f func(Something) Parser) Parser{
	var ret = func (str string) (Something, string, bool) {
		r, rest, status := p(str)
		if !status {
			return "", str, false
		}
		parser := f(r)
		return parser(rest)
	}
	return ret
}


func Satisfy(s func(rune) bool) Parser {
	var f = func (x Something) Parser {
		r := x.(rune)
		if s(r) {
			return Result(x)
		}else{
			return Zero()
		}
	}
	return Bind(Item(), f) 
}

func Rune (r rune) Parser {
	var f = func (x rune) bool {
		return x == r
	}
	return Satisfy(f)
}

func digit() Parser {
	var f = func (x rune) bool {
		return x >= '0' && x <= '9'
	}
	return Satisfy(f)
}


func lower() Parser {
	var f = func (x rune) bool {
		return x >= 'a' && x <= 'z'
	}
	return Satisfy(f)

}

func upper() Parser {
	var f = func (x rune) bool {
		return x >= 'A' && x <= 'Z'
	}
	return Satisfy(f)

}