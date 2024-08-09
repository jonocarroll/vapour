# Vapour

A typed, imperative, superset of the [R programming language](https://www.r-project.org/)

> [!WARNING]  
> This is a work in progress

```r
type person: object {
  age: int,
  name: char 
}

func create(name: char): person {
  return person(name = name)
}

func(p: person) set_age(age: int): null {
  p$age = age
}

let john: person = create("john")

set_age(john, 36)
```
