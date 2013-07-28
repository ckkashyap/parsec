package parsec

import (
	"testing"
)

func TResult(test *testing.T) {
	// This should create a praser which when executed produces a
	// channel which when we read from will produce the value 10
	expectedValue := int(10)
	expectedRemainingString := "hello"
	p := Result(expectedValue)

	r1 := p(expectedRemainingString)
	
	ctr := 0
	for i := range r1 {
		if ctr > 0 {
			test.Fatalf("Got more values than expected")
		}
		switch tt := i.(type) {
		case Tup:
			if !tt.Valid {
				test.Fatalf("Expected the parse to succeed")
			}
			if tt.Remaining != expectedRemainingString {
				test.Fatalf("Expected %s as remaining string but got %s", expectedRemainingString, tt.Remaining)
			}
			switch ttt:=tt.Thing.(type) {
			case int:
				if ttt != expectedValue {
					test.Fatalf("Expected %d got %d", expectedValue, ttt)
				}
			default:
				test.Fatalf("Expected int got something else")
			}
		default:
			test.Fatalf("Expected Tup got something else")
		}
		ctr++
	}
}



func Test_Rune(test *testing.T) {
	// Rune is a parser that takes a Rune and returns a parser
	// that matches that rune. If the rune matches then the result
	// is that rune and the remaining string Tup
	p := Rune('A')
	rSuccess := p("ABC")

	var ctr int

	ctr = 0
	for i := range rSuccess {
		if ctr > 0 {
			test.Fatalf("Got more values than expected")
		}
		switch tt := i.(type) {
		case Tup:
			if !tt.Valid {
				test.Fatalf("Expected the parse to succeed")
			}
			if tt.Remaining != "BC" {
				test.Fatalf("Expected BC as remaining string but got %s", tt.Remaining)
			}
			switch ttt:=tt.Thing.(type) {
			case rune:
				if ttt != 'A' {
					test.Fatalf("Expected A got %c", ttt)
				}
			default:
				test.Fatalf("Expected rune got something else")
			}
		default:
			test.Fatalf("Expected Tup got something else")
		}
		ctr++
	}
}

func Test_Runes(test *testing.T) {
	//Create a parser by composing 2 parsers. Using simple parsers
	//that can parse just one rune each, a complex parser that
	//parses both the runes can be built using the Bind function.
	p1:=Rune('A')
	p2:=Rune('B')
	p := Bind(p1,func (x1 A) Parser {
		return Bind(p2, func (x2 A) Parser {
			return func (str string) chan A {
				c := make(chan A)
				r1 := (x1.(Tup).Thing.(rune))
				r2 := (x2.(Tup).Thing.(rune))
				go func () {
					c <- r1
					c <- r2
					close (c)
				}()
				return c
			}
		})
	})

	result := p("ABC")
	runes := []rune{'A', 'B'}

	i := 0
	for r := range result {
		if r != runes[i] {
			test.Fatalf("Expected %c, got %c", runes[i], r)
		}
		i++
	}
}
