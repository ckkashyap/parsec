package parsec

import (
	"strings"
)

type Any interface{}

type Res struct {
	a   Any
	str string
}


func (t Res)GetResult() Any {
	return t.a
}
func (t Res)GetRest() string {
	return t.str
}

type Parser func(string) []Res

func Result(a Any) Parser {
	var ret = func(str string) []Res {
		return []Res{{a, str}}
	}
	return ret
}

func Zero() Parser {
	var ret = func(str string) []Res {
		return []Res{}
	}
	return ret
}

func Item() Parser {
	var ret = func(str string) []Res {
		r := strings.NewReader(str)
		c, s, _ := r.ReadRune()
		if s == 0 {
			return []Res{}
		}
		return []Res{{c, str[s:]}}
	}
	return ret
}

func Bind(p Parser, f func(Any) Parser) Parser {
	var ret = func(str string) []Res {
		result := p(str)
		var totalResults = len(result)

		if totalResults == 0 {
			return result
		}

		for _, res := range result {
			rest := res.str
			parseResult := res.a
			parser := f(parseResult)
			r1 := parser(rest)
			if len(r1) == 0 {continue}
			return r1[0:1] //Ignoring other successful results
		}
		return []Res{}
	}
	return ret
}

func Satisfy(s func(rune) bool) Parser {
	var f = func(x Any) Parser {
		r := x.(rune)
		if s(r) {
			return Result(x)
		} else {
			return Zero()
		}
	}
	return Bind(Item(), f)
}

func Rune(r rune) Parser {
	var f = func(x rune) bool {
		return x == r
	}
	return Satisfy(f)
}

func Digit() Parser {
	var f = func(x rune) bool {
		return x >= '0' && x <= '9'
	}
	return Satisfy(f)
}

func Lower() Parser {
	var f = func(x rune) bool {
		return x >= 'a' && x <= 'z'
	}
	return Satisfy(f)

}

func Upper() Parser {
	var f = func(x rune) bool {
		return x >= 'A' && x <= 'Z'
	}
	return Satisfy(f)

}

func Plus(p1, p2 Parser) Parser {
	var ret = func(str string) []Res {
		r1 := p1(str)
		r2 := p2(str)
			lr1 := len(r1)
			lr2 := len(r2)
		slice := make([]Res, lr1 + lr2)
		copy(slice, r1)
		copy(slice[lr1:], r2)
		return slice
	}
	return ret
}



func Letter() Parser {
	p := Plus(Lower(), Upper())
	return p
}

func AlphaNum() Parser {
	p := Plus(Letter(), Digit())
	return p
}

func Word() Parser {
	return Many(Letter())
}


func Many(p Parser) Parser {
	neWord := Bind(p, func (x Any) Parser {
		return Bind(Many(p), func(as Any) Parser {
			xs := as.([]Any)
			slice := make([]Any, len(xs) + 1 )
			slice[0] = x
			copy(slice[1:], xs)
			return Result(slice)

		})
	})
	q := Plus(neWord, Result([]Any{}))
	return q
}