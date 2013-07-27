parsec
======

This is my attempt at translating the Parsec from this paper - Monadic
Parser Combinators [1]

The beauty of this idea (and of combinators in general) is that way in
which we can build complex things out of really simple building
blocks.

For example the paper starts off with a general description of what a
parser is - essentially a function that takes a string and transforms
it into an AST

So something like

    type Parser = String -> Tree

Ofcourse, a parser needs to also consider a parse failure and also
return the remaining string that was not parsed - so the definition is
changed to

    type Parser a = String -> [(a, String)]

The above can be read as - Parser is a parameterized type - as in
Parser Char is a Parser that when executed produces a list of 2 tuple
containing the parsed Character and the remaining String.






[1] http://eprints.nottingham.ac.uk/237/1/monparsing.pdf
