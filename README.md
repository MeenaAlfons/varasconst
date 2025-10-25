# varasconst

Go Linter to mark vars as const

## Use case

If you want to define a a global var and make sure it's not overwritten later, mark it with `// const`

```go
package p1

// const
var Global_Const_1 = ""

func main() {
    // lint error: "assignment to global variable marked with const"
	Global_Const_1 = "modified" 
}
```

It also works across packages

```go
package p2

import "p1"

func main() {
    // lint error: "assignment to global variable marked with const"
	p1.Global_Const_1 = "modified" 
}
```
