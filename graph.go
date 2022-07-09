package graph

import (
	"fmt"
	"sort"
	"sync"
)

type NodeId int

type Node[N any, E any] struct {
	id                  NodeId
	Data                N
	IncomingConnections map[NodeId]*Edge[N, E]
	OutgoingConnections map[NodeId]*Edge[N, E]
}

func (n Node[N, E]) Id() NodeId {
	return n.id
}

type EdgeId struct {
	From NodeId
	To   NodeId
}

type Edge[N any, E any] struct {
	id   EdgeId
	Data E
	From *Node[N, E]
	To   *Node[N, E]
}

func (e Edge[N, E]) Id() EdgeId {
	return e.id
}

type Graph[N any, E any] struct {
	nodes      map[NodeId]*Node[N, E]
	edges      map[EdgeId]*Edge[N, E]
	nodeIdIncr int
	rw         sync.RWMutex
}

func NewGraph[N any, E any]() Graph[N, E] {
	return Graph[N, E]{
		nodes:      map[NodeId]*Node[N, E]{},
		edges:      map[EdgeId]*Edge[N, E]{},
		nodeIdIncr: 0,
		rw:         sync.RWMutex{},
	}
}

func NewGraphFrom[N any, E any](nodes []N, edges map[EdgeId]E) (Graph[N, E], error) {
	d := NewGraph[N, E]()
	for _, n := range nodes {
		d.addNode(n)
	}
	for id, v := range edges {
		_, err := d.addEdge(id.From, id.To, v)
		if err != nil {
			return d, err
		}
	}
	return d, nil
}

func (d *Graph[N, E]) Clone() (Graph[N, E], error) {
	d.rw.RLock()
	defer d.rw.RUnlock()
	return d.clone()
}

func (d *Graph[N, E]) clone() (Graph[N, E], error) {
	nodeIds := []int{}
	for _, n := range d.nodes {
		nodeIds = append(nodeIds, int(n.id))
	}
	sort.Sort(sort.IntSlice(nodeIds))
	nodes, edges := []N{}, map[EdgeId]E{}
	for _, id := range nodeIds {
		nodes = append(nodes, d.getNode(NodeId(id)).Data)
	}
	for _, e := range d.edges {
		edges[e.id] = d.getEdge(e.id.From, e.id.To).Data
	}
	return NewGraphFrom(nodes, edges)
}

func (d *Graph[N, E]) AddNode(n N) NodeId {
	d.rw.Lock()
	defer d.rw.Unlock()
	return d.addNode(n)
}

func (d *Graph[N, E]) addNode(n N) NodeId {

	id := NodeId(d.nodeIdIncr)
	d.nodes[id] = &Node[N, E]{
		id:                  id,
		Data:                n,
		IncomingConnections: map[NodeId]*Edge[N, E]{},
		OutgoingConnections: map[NodeId]*Edge[N, E]{},
	}
	d.nodeIdIncr += 1
	return id
}

func (d *Graph[N, E]) GetNode(id NodeId) *Node[N, E] {
	d.rw.RLock()
	defer d.rw.RUnlock()
	return d.getNode(id)
}

func (d *Graph[N, E]) getNode(id NodeId) *Node[N, E] {
	if v, ok := d.nodes[id]; ok {
		return v
	}
	return nil
}

func (d *Graph[N, E]) RemoveNode(id NodeId) error {
	d.rw.Lock()
	defer d.rw.Unlock()
	return d.removeNode(id)
}

func (d *Graph[N, E]) removeNode(id NodeId) error {
	node := d.getNode(id)
	if node == nil {
		return fmt.Errorf("Node %d not found to remove", id)
	}
	for edge, _ := range d.edges {
		if edge.From == id || edge.To == id {
			err := d.removeEdge(edge.From, edge.To)
			if err != nil {
				return err
			}
		}
	}
	delete(d.nodes, id)
	return nil
}

func (d *Graph[N, E]) AddEdge(from, to NodeId, value E) (EdgeId, error) {
	d.rw.Lock()
	defer d.rw.Unlock()
	return d.addEdge(from, to, value)
}

func (d *Graph[N, E]) addEdge(from, to NodeId, value E) (EdgeId, error) {
	fromNode, toNode := d.getNode(from), d.getNode(to)

	if fromNode == nil {
		return EdgeId{}, fmt.Errorf("Node %d not found", from)
	}
	if toNode == nil {
		return EdgeId{}, fmt.Errorf("Node %d not found", to)
	}
	id := EdgeId{from, to}
	if _, ok := d.edges[id]; ok {
		return EdgeId{}, fmt.Errorf("Edge %d -> %d already exists", from, to)
	}
	edge := Edge[N, E]{id, value, fromNode, toNode}

	// do it
	fromNode.OutgoingConnections[to] = &edge
	toNode.IncomingConnections[from] = &edge
	d.edges[id] = &edge

	return id, nil
}

func (d *Graph[N, E]) GetEdge(from, to NodeId) *Edge[N, E] {
	d.rw.RLock()
	defer d.rw.RUnlock()
	return d.getEdge(from, to)
}

func (d *Graph[N, E]) getEdge(from, to NodeId) *Edge[N, E] {
	if v, ok := d.edges[EdgeId{from, to}]; ok {
		return v
	}
	return nil
}

func (d *Graph[N, E]) RemoveEdge(from, to NodeId) error {
	d.rw.Lock()
	defer d.rw.Unlock()
	return d.removeEdge(from, to)
}

func (d *Graph[N, E]) removeEdge(from, to NodeId) error {
	fromNode, toNode := d.getNode(from), d.getNode(to)
	if fromNode == nil {
		return fmt.Errorf("Node %d not found", from)
	}
	if toNode == nil {
		return fmt.Errorf("Node %d not found", to)
	}

	// do it
	delete(d.edges, EdgeId{from, to})
	delete(fromNode.OutgoingConnections, to)
	delete(toNode.IncomingConnections, from)

	return nil
}

// https://stackoverflow.com/a/4577
func (original *Graph[N, E]) TopologicalSort() ([]NodeId, error) {
	d, err := original.Clone()
	if err != nil {
		return []NodeId{}, err
	}
	noIncoming := []NodeId{}
	for _, node := range d.nodes {
		if len(node.IncomingConnections) == 0 {
			noIncoming = append(noIncoming, node.id)
		}
	}
	l := []NodeId{}
	for len(noIncoming) > 0 {
		if len(l) > 40 {
			panic(">40")
		}
		a := noIncoming[0]
		noIncoming = noIncoming[1:]
		l = append(l, a)
		for _, e := range d.getNode(a).OutgoingConnections {
			m := e.To
			d.removeEdge(e.id.From, e.id.To)
			if len(m.IncomingConnections) == 0 {
				noIncoming = append(noIncoming, m.id)
			}
		}
	}
	if len(d.edges) > 0 {
		return []NodeId{}, fmt.Errorf("Graph is cyclic")
	}
	return l, nil
}

func (d *Graph[N, E]) String() string {
	s := "Nodes:\n"
	for _, node := range d.nodes {
		s += fmt.Sprintln(" ", node.id, "=", node.Data)
	}
	s += "Edges:\n"
	for _, edge := range d.edges {
		s += fmt.Sprintln(" ", edge.id.From, "->", edge.id.To, "=", edge.Data)
	}
	return s
}
