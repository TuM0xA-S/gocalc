# interactive calculator

### features
* [x] float numbers, basic operators(+, -, *, /, (, )) - ~~only binary supported for now~~
* [x] basic error handling
* [x] tokenizer is enough smart to proccess arbitrary formatted expressions
* [x] variables
* [x] parser upgraded to work with unary operators
* [x] simple functions (one line formulas, without subcalls)
* [ ] script mode, options
* [ ] uses readline library(interactive editing)
* [ ] enhanced error handling with indication of problem position in input
* [ ] optimize function expressions

### testing
`go test .`

### building 
`go build -o calc ./cmd`
`./calc` - to run

### quick run
`go run ./cmd`

### syntax
* identifier: starts with letter, can consist of letters and digits(case-sensetive)
* number: floating point number (dot as fraction separator)
* variable:
  * assignment: `var = 2 + 2` variable `var` with value 4
  * usage: `3 * var` => 12
* function:
  * declaration: `@func = (a, b): 2 * a + b` function `@func` with two parameters
  * usage: `@func(3, 5)` => 11
* expression: consists of numbers, operators, function calls, variables
* operators:
  * unary: `+-`
  * binary: `+-/*`
  * parentheses: `()`
* instruction:
  * variable assignment
  * function declaration
  * expression


