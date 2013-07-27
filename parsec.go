package parsec

import (
	"strings"
	"fmt"
)

type A interface{}

type Tup struct {
	Thing A
	Remaining string
	Valid bool
}


func (t Tup)String() string{
	if t.Valid {
		return fmt.Sprintf("Thing=%v, Remaining=%s", t.Thing, t.Remaining)
	} else {
		return "TERM"
	}
}

type Parser func(string) chan A

func Result(a A) Parser {
	var ret = func(str string) chan A {
		c := make(chan A)
		go func () {
			c <- Tup{a, str, true}
			close(c)
		}()
		return c
	}
	return ret
}


func Zero() Parser {
	var ret = func(str string) chan A {
		c := make(chan A)
		go func () {
			close (c)
		}()
		return c
	}
	return ret
}


func Item() Parser {
	var ret = func(str string) chan A {
		c := make(chan A)
		go func () {
			rdr := strings.NewReader(str)
			r, idx, _ := rdr.ReadRune()
			if idx == 0 {
				close(c)
			}else{
				c <- Tup{r, str[idx:], true}
				close(c)
			}
		}()
		return c
	}
	return ret
}


func Bind(p Parser, f func(A) Parser) Parser {
	var ret = func(str string) chan A {
		channel := make(chan A)
		go func () {
			r1 := p(str)
			for eachR1 := range r1 {
				parser := f(eachR1)
				r2 := parser(eachR1.(Tup).Remaining)
				for eachR2 := range r2 {
					channel <-eachR2
				}
			}
			close(channel)
		}()
		return channel
	}
	return ret
}




func Plus(p1, p2 Parser) Parser {
	var ret = func(str string) chan A {
		channel := make(chan A)
		go func () {
			for r := range p1(str) {
				channel<-r
			}
			for r := range p2(str) {
				channel<-r
			}
			close(channel)

		}()
		return channel
	}
	return ret
}


func Satisfy(s func(rune) bool) Parser {

	var f = func(x A) Parser {
		r := x.(Tup).Thing.(rune)
		if s(r) {
			return Result(r)
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

func Letter() Parser {
	p := Plus(Lower(), Upper())
	return p
}

func AlphaNum() Parser {
	p := Plus(Letter(), Digit())
	return p
}

func Many(p Parser) Parser {
	p1 := Bind(p, func (x A) Parser {
		theTuple :=  x.(Tup)
		return Bind(Many(p), func(xs A) Parser {
			c := make(chan A)
			inpChanPrime := xs.(Tup).Thing.(chan A)
			go func(){
				c <- theTuple
				for r := range inpChanPrime {
					r1 := r.(Tup)
					c <- r1
				}
				close(c)
			}()

			return Result(c)
			
		})
	})

	c := make(chan A)
	go func(){
		c <- Tup{Valid:false}
		close(c)
	}()
	
	return Plus(p1, Result(c))
}

func Ident() Parser {
	return Bind (Lower(), func (x A) Parser {
		theTuple := x.(Tup)
		return Bind(Many(AlphaNum()), func(xs A) Parser{
			c := make(chan A)
			inpChanPrime := xs.(Tup).Thing.(chan A)
			go func() {
				c <- theTuple
				for r := range inpChanPrime {
					r1 := r.(Tup)
					c <- r1
				}
				close(c)
			}()
			return Result(c)
		})
	})
}
