package graph_test

import (
	"testing"

	graph "github.com/mcpar-land/generic-graph"
)

type Foo struct {
	Value string
}

type Bar struct {
	X string
	Y string
}

func TestGraph(t *testing.T) {
	t.Log("Creating dag")
	g := graph.NewGraph[Foo, Bar]()
	t.Log("Adding dag nodes")
	a := g.AddNode(Foo{"This is node a"})
	b := g.AddNode(Foo{"This is node b"})
	c := g.AddNode(Foo{"This is node c"})

	t.Log("Adding dag edges")
	_, err := g.AddEdge(a, b, Bar{"from a...", "...to b"})
	if err != nil {
		t.Error(err)
	}
	_, err = g.AddEdge(b, c, Bar{"from b...", "...to c"})
	if err != nil {
		t.Error(err)
	}
	_, err = g.AddEdge(a, c, Bar{"from a...", "...to c"})
	if err != nil {
		t.Error(err)
	}
	t.Log("Done")

	t.Log(g.String())

	aNode := g.GetNode(a)
	if aNode == nil {
		t.Fail()
	}
	for _, out := range aNode.OutgoingConnections {
		t.Log(out, "-->", out.To)
	}
	if len(aNode.OutgoingConnections) != 2 {
		t.Fail()
	}
	cNode := g.GetNode(c)
	if cNode == nil {
		t.Fail()
	}
	if len(cNode.IncomingConnections) != 2 {
		t.Fail()
	}
	err = g.RemoveNode(b)
	if err != nil {
		t.Error(err)
	}
	if len(aNode.OutgoingConnections) != 1 {
		t.Fail()
	}
	if len(cNode.IncomingConnections) != 1 {
		t.Fail()
	}
	t.Log(g.String())
	err = g.RemoveEdge(a, c)
	if err != nil {
		t.Error(err)
	}
	t.Log("removed edge")
	t.Log(g.String())
	d := g.AddNode(Foo{"this is node d"})

	_, err = g.AddEdge(d, a, Bar{"from d...", "...to a"})
	if err != nil {
		t.
			Error(err)
	}
	_, err = g.AddEdge(d, c, Bar{"from d...", "...to c"})
	if err != nil {
		t.Error(err)
	}

	t.Log(g.String())
}

func TestCyclicSort(t *testing.T) {
	g := graph.NewGraph[Foo, Bar]()
	a := g.AddNode(Foo{"This is node a"})
	b := g.AddNode(Foo{"This is node b"})
	c := g.AddNode(Foo{"This is node c"})
	g.AddEdge(a, b, Bar{"from a...", "...to b"})
	g.AddEdge(b, c, Bar{"from b...", "...to c"})
	g.AddEdge(a, c, Bar{"from a...", "...to c"})

	sort, err := g.TopologicalSort()
	if err != nil {
		t.Error(err)
	}
	t.Log(sort)

	// cyclic sort should fail if the graph has ANY cycles
	g.AddEdge(c, a, Bar{"cyclic from c...", "...back to a"})
	_, err = g.TopologicalSort()
	if err == nil {
		t.Errorf("Topological sort should fail on graph with cycles.")
	}
}
