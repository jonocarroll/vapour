package lexer

import (
	"testing"

	"github.com/devOpifex/vapour/token"
)

func TestBasicTypes(t *testing.T) {
	code := `const x: int | na = 1

let y: []int = list(1, 23, 33)

# structure(1..10, name = "", id = 0)
type item struct {
  int
  # attributes
  category string
}

item(1, category = "")

# structure(item,(), name = "", id = 0)
type nested struct {
  item
  # attributes
  name string
  id int
}

nested(
  item(1..10, name = "hello", id = 1),
  category = "test"
)

# data.frame(name = c("a", "z"), id = 1..2)
type df dataframe {
  name string
  id int
}

df(name = "hello", id = 1)

# list(1, 2, 3)
type lst list {
  int
}

lst(1,2,3)

# list(name = "hello", id = 1)
type obj object {
  id int
  n num
}

obj(
  id = 0,
  n = 3.14
)

# list(list(name = "hello", id = 1))
type objs []object

objs(
  obj(),
  obj()
)`

	l := &Lexer{
		Input: code,
	}

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	l.Print()

	expectLength := 25
	if len(l.Items) != expectLength {
		t.Fatalf("Expecting %v tokens, go %v", expectLength, len(l.Items))
	}

	if l.getItem(0).Class != token.ItemConst {
		t.Fatalf("Expecting constant, got %v - %v", l.getItem(0).String(), l.getItem(0).Value)
	}

	if l.getItem(5).Class != token.ItemTypes {
		t.Fatalf("Expecting type, got %v - %v", l.getItem(5).String(), l.getItem(5).Value)
	}
}
