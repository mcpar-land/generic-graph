# Generic Directed Graph

Implementation of a directed graph in Golang, making use of 2022 generics.

```go

import (
	"fmt"
	graph "github.com/mcpar-land/generic-graph"
)

type Foo struct {
	Value string
}

type Bar struct {
	X string
	Y string
}

func main() {
	// The graph takes two generic types, one for the node and one for the edge.
	g := graph.NewGraph[Foo, Bar]()
	a := g.AddNode(Foo{"This is node a"})
	b := g.AddNode(Foo{"This is node b"})
	c := g.AddNode(Foo{"This is node c"})
	g.AddEdge(a, b, Bar{"from a...", "...to b"})
	g.AddEdge(b, c, Bar{"from b...", "...to c"})
	g.AddEdge(a, c, Bar{"from a...", "...to c"})

	fmt.Println(g.String())

	for _, edge := range g.GetNode(a).OutgoingConnections {
		fmt.Println("node a is connected to: ", edge.To)
	}
}
```

```
$ go run ./examples/basic
Nodes:
  0 = {This is node a}
  1 = {This is node b}
  2 = {This is node c}
Edges:
  0 -> 1 = {from a... ...to b}
  1 -> 2 = {from b... ...to c}
  0 -> 2 = {from a... ...to c}

node a is connected to:  {This is node b}
node a is connected to:  {This is node c}
```
