package node

import (
	"fmt"
	"kdfx/pkg/context"
)

// graph represents a collection of nodes and their connections.
type graph struct {
	nodes map[string]Node
}

// NewGraph creates a new empty graph.
func NewGraph() Graph {
	return &graph{
		nodes: make(map[string]Node),
	}
}

// AddNode adds a node to the graph with a unique name.
func (g *graph) AddNode(name string, node Node) {
	g.nodes[name] = node
}

// Connect connects the output of sourceNode to the input slot of targetNode.
func (g *graph) Connect(sourceNodeName, targetNodeName, inputSlot string) error {
	source, ok := g.nodes[sourceNodeName]
	if !ok {
		return fmt.Errorf("source node %s not found", sourceNodeName)
	}
	target, ok := g.nodes[targetNodeName]
	if !ok {
		return fmt.Errorf("target node %s not found", targetNodeName)
	}

	target.SetInput(inputSlot, source)
	return nil
}

func (g *graph) GetNode(name string) Node {
	return g.nodes[name]
}

func (g *graph) Release() {
	for _, node := range g.nodes {
		node.Release()
	}
}

// pipeline manages the execution of a graph.
type pipeline struct {
	graph   Graph
	context context.Context
}

// NewPipeline creates a new pipeline for a given graph and context.
func NewPipeline(ctx context.Context, graph Graph) Pipeline {
	return &pipeline{
		graph:   graph,
		context: ctx,
	}
}

// Execute runs the pipeline for a specific output node.
// It recursively processes dependencies.
func (p *pipeline) Execute(outputNodeName string) error {
	node := p.graph.GetNode(outputNodeName)
	if node == nil {
		return fmt.Errorf("output node %s not found", outputNodeName)
	}

	return node.Process(p.context)
}

// Release releases all resources in the graph.
// Release releases all resources in the graph.
func (p *pipeline) Release() {
	if p.graph != nil {
		p.graph.Release()
	}
}
