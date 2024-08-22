package walker

import (
	"fmt"
	"testing"

	"github.com/devOpifex/vapour/lexer"
	"github.com/devOpifex/vapour/parser"
)

func TestEnvironment(t *testing.T) {
	code := `
let z: int = 1

func addz(n: int, y: int): int | na {
	if(n == 1){
		return na
	}

	return n + y
}

# should fail, this can be na
let result: int = addz(1, 2)

# should fail, comparing wrong types
if (1 == "hello") {
  print("1")
}

# should fail?
if (1 > 2.1) {
  print("1")
}

const y: int = 1

# should fail, is constant
y = 2
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	fmt.Println("----------------------------- Env")
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestInfix(t *testing.T) {
	code := `let x: char = "hello"

# should fail, cannot be NA
x = NA

# should fail, types do not match
let z: char = 1

func add(n: int, y: int): int | na {
	if n == 1 {
		return NA
	}

  return n + y
}

# should fail, this can be na
let result: int = add(1, 2)

# should fail, const must have single type
const v: int | na = 1

v = 2
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	fmt.Println("----------------------------- infix")
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestNamespace(t *testing.T) {
	code := `
# should fail, duplicated params
func foo(x: int, x: int): int {return x + y}

# should fail, duplicated params
func (x: int) bar(x: int): int {return x + y}
`

	l := lexer.NewTest(code)

	l.Run()
	fmt.Println("----------------------------- ns")
	p := parser.New(l)

	prog := p.Run()
	if p.HasError() {
		for _, e := range p.Errors() {
			fmt.Println(e)
		}
		return
	}

	w := New()
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestNumber(t *testing.T) {
	code := `let x: num = 1

x = 1.1

let u: int = 1e10

let integer: int = 1

# should fail, assign num to int
integer = 2.1

let s: int = sum(1,2,3)
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	w.Run(prog)

	fmt.Println("----------------------------- number")
	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestCall(t *testing.T) {
	code := `func foo(x: int, y: char): int {
  print(y)
	return x + 1
}

# should fail, argument does not exist
foo(z = 2)

# should fail, wrong type
foo(
  x = "hello"
)

# should fail, wrong type
foo("hello")

# should fail, too many arguments
foo(1, "hello", 3)

func lg (...: char): null {
  print(...)
}

lg("hello", "world")

# should fail, wrong type
lg("hello", 1)
`

	l := lexer.NewTest(code)

	l.Run()
	fmt.Println("----------------------------- call")
	p := parser.New(l)

	prog := p.Run()

	w := New()
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestMissing(t *testing.T) {
	code := `
# should warn, can be missing
func hello(what: char): char {
  sprintf("hello, %s!", what)
}

hello()
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	fmt.Println("----------------------------- missing")
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestExists(t *testing.T) {
	code := `
# should fail, x does not exist
x = 1

pkg::fn(x = 2)

dplyr::filter(x = 2)

# should fail, y does not exist
func foo(y: int): int {
 return z
}
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	fmt.Println("----------------------------- exists")
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestList(t *testing.T) {
	code := `
type person: list {
	name: char
}

type persons: []person

let peoples: persons = persons(
  person(name = "John"),
  person(name = "Jane")
)

func foo(callback: fn): any {
  return callback()
}

foo((x: int): int => {return x + 1})
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	fmt.Println("----------------------------- decorator")
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestTypeMatch(t *testing.T) {
	code := `
type userid: int

type config: struct {
  char,
	x: int
}

type inline: object {first: int, second: char}

type lst: list {int | num}

lst(2)

config(2, x = 2)

# should fail, must be named
inline(1)

# should fail, first arg of struct cannot be named
config(u = 2, x = 2)

# should fail, struct attribute must be named
config(2, 2)

# should fail, does not exist
inline(
  z = 2
)

type a_function: func(x: int, y: int): int

func foo(callback: a_function, y: int): int {
  return callback(1, y)
}

foo((x: int, y: int): int => {
  return x + y
}, 2)

func bar(x: int, y: int): int {
  return x + y
}

foo(bar, z)
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	fmt.Println("----------------------------- typematch")
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestAnonymous(t *testing.T) {
	code := `
# should fail, returns wrong type
lapply(1..10, (z: int): int => {
  return "hello"
})

type math: func(x: int): int

func apply_math(vector: int, cb: math): int {
  return cb(vector)
}

apply_math((1, 2, 3), (x: int): int => {
  return x * 3
})
`

	l := lexer.NewTest(code)

	l.Run()
	fmt.Println("----------------------------- anon")
	p := parser.New(l)

	prog := p.Run()

	w := New()
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestR(t *testing.T) {
	code := `
# should fail, package not installed
xxx::foo()

dplyr::wrong_function()

let x: int = 1

if x == 1 {
  x <- 2
}

# should fail does not exist
y <- 2
`

	l := lexer.NewTest(code)

	l.Run()
	fmt.Println("----------------------------- R")
	p := parser.New(l)

	prog := p.Run()

	w := New()
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestFunction(t *testing.T) {
	code := `
func foo(n: int): int {
  let x: char = "hello"

  # should fail, returns wrong type
  if (n == 2) {
    return "hello"
  }

  if (n == 3) {
    return 1.2
  }

  # should fail, returns wrong type
  return x

  # should fail, returns does not exist
  return u
}
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	fmt.Println("----------------------------- function")
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestUnused(t *testing.T) {
	code := `
let x: int = 10

let total: int = x + 32

total + 1

# should warn of unused variable
let y: int = 1

type userid: int

type train: list {
  wheels: int
}

let t: train = train(wheels = 256)

# should warn that function foo is not used
func foo(): null {
  print("hello")
}

# should warn that function foo is already declared
func foo(): int {
  return 1
}

as.data.frame(cars)
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	fmt.Println("----------------------------- unused")
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestSquare(t *testing.T) {
	code := `let x: int = (1,2,3)

x[2] = 3

let y: int = list(1,2,3)

y[[1]] = 1

let zz: char = ("hello|world", "hello|again")
let z: char = strsplit(zz[2], "\\|")[[1]]

x[1, 2] = 15

x[[3]] = 15

type xx: dataframe{
  name: int
}

let df: xx = xx(name = 1)

df$name = "hello"

# should fail, not generic
func (p: any) meth(): null {}
`

	l := lexer.NewTest(code)

	l.Run()
	fmt.Println("----------------------------- square")
	p := parser.New(l)

	prog := p.Run()

	w := New()
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestDecorator(t *testing.T) {
	code := `
@class(int, person)
type man: struct {
  int,
	name: char
}

let p: man = man(1)

func (x: person) print_id(): null {
  print(x$int)
}
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	fmt.Println("----------------------------- decorator")
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestBasic(t *testing.T) {
	code := `let x: int | na = 1

x = 2

# should fail, it's already declared
let x: char = "hello"

type id: struct {
	int,
	name: char
}

# should fail, already defined
type id: int

# should fail, cannot find type
let z: undefinedType = "hello"

# should fail, different types
let v: int = (10, "hello", NA)

# should fail, type mismatch
let wrongType: num = "hello"

# should fail, must have a value
const xx: int

if(xx == 1) {
	let x: int = 2
}

# should fail wrong types
let x: int = "hello" + 2

# should fail, x is int, expression coerces to num
x = 1 + 1.2

let Z: num = 1.2 + 3

# should fail, does not exist
x = 2

let uu: int = 2

# should fail, does not exist
uu = "char"
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	fmt.Println("----------------------------- Basic")
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestRecursiveTypes(t *testing.T) {
	code := `
type userid: int

type user: struct {
  userid,
	name: char
}

func create(id: userid): user {
  return user(id)
}

create(2)

type person: struct {
  char,
	name: string
}

# should fail, wrong type
person(2)

@inherits(person)
type man: struct {
  userid
}

man(2)
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	fmt.Println("----------------------------- Recurse types")
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestListTypes(t *testing.T) {
	code := `
type userid: int

type user: object {
	name: char
}

type users: []user

#should fail, wrong type
let z: users = users(
  user(name = "john"),
	4
)
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	fmt.Println("----------------------------- list types")
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestFor(t *testing.T) {
	code := `
type userid: int

let x: userid = 1

for(let i: int in x..10) {
  print(i)
}

type x: struct {
  int
}

let y: x = x(2)

# should fail, range has char..int
for(let i: int in y) {
  print(i)
}
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	fmt.Println("----------------------------- for")
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestMethod(t *testing.T) {
	code := `func (o: obj) add(n: int): char {
  return "hello"
}

type person: struct {
  int,
	name: char
}

# should fail, xxx does not exist
func (p: person) setName(name: char): null {
  p$xxx = 2
}

# should fail, name expects char
func (p: person) setName(name: char): null {
  p$name = 2
}
`

	l := lexer.NewTest(code)

	l.Run()
	fmt.Println("----------------------------- method")
	p := parser.New(l)

	prog := p.Run()

	w := New()
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}
