package parsec

import (
	"fmt"
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
		var tempArray = make([][]Res, totalResults, totalResults)
		validCount := 0
		totalCount := 0
		for _, res := range result {
			rest := res.str
			parseResult := res.a
			parser := f(parseResult)
			r1 := parser(rest)
			if len(r1) == 0 {
				continue
			}
			tempArray[validCount] = r1
			validCount++
			totalCount += len(r1)
		}
		var finalArray = make([]Res, totalCount, totalCount)
		i:=0
		for _, ta := range tempArray {
			for _, r := range ta {
				finalArray[i] = r
				i++
			}
		}
		return finalArray
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
	neWord := Bind(Letter(), func (c Any) Parser {
		return Bind(Word(), func(cs Any) Parser {
			r := c.(rune)
			str := cs.(string)
			return Result(fmt.Sprintf("%c%s",r,str))

		})
	})

	p := Plus (neWord, Result(""))
	return p
}


func Many(p Parser) Parser {
	abc := Zero()
	q := Plus(abc, Result([]Any{}))
	return q
}