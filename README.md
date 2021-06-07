# interactive calculator

### features
* [x] float numbers, basic operators(+, -, *, /, (, ))
* [x] basic error handling
* [x] tokenizer is enough smart to proccess arbitrary formatted expressions
* [x] variables
* [x] parser upgraded to work with unary operators
* [x] simple functions (one line formulas, without subcalls)
* [ ] meta commands(show all variables and functions)
* [ ] script mode, options
* [x] uses readline library(interactive editing)
* [ ] autocomplete
* [x] enhanced error handling with indication of problem position in input
* [ ] optimize(minimize) function expressions
* [ ] api to interact with interpreter objects from go code

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
  * variable_name: identifier
  * assignment: `varable_name = expression` variable `variable_name` with value of `expression`
    * example: `var = 2 * 2` variable `var` with value 4
  * usage: `variable_name` => gives value of variable `variable_name`
    * example: `var` => 6

* function:
  * function_name: @identifier
  * declaration: `function_name = (variable_name [,variable_name]): expression` function with name `function_name` with zero or more parameters(separated with comma), that used for calculate `expression`
    * example: `@foo = (a, b): 2 * a - b`
  * usage: `function_name(expression [,expression])` call function `function_name`
    * example: `@foo(4 - 1, 2)` => 4

* expression: consists of numbers, operators, function calls, variables
  * example `-(a - @bar(1, (2.34 + c) * b)) * 5.1 - d / (100 - 1)`

* operators:
  * unary: `+-`
  * binary: `+-/*`
  * parentheses: `()`

* instruction:
  * variable assignment (create variable)
  * function declaration (create function)
  * expression (calculate and print value)

* interpreter: processes instructions


