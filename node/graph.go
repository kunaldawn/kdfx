package node

import (
	"fmt"
	"kimg/context"
)

// Graph represents a collection of nodes and their connections.
type Graph struct {
	Nodes map[string]Node
}

// NewGraph creates a new empty graph.
func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]Node),
	}
}

// AddNode adds a node to the graph with a unique name.
func (g *Graph) AddNode(name string, node Node) {
	g.Nodes[name] = node
}

// Connect connects the output of sourceNode to the input slot of targetNode.
func (g *Graph) Connect(sourceNodeName, targetNodeName, inputSlot string) error {
	source, ok := g.Nodes[sourceNodeName]
	if !ok {
		return fmt.Errorf("source node %s not found", sourceNodeName)
	}
	target, ok := g.Nodes[targetNodeName]
	if !ok {
		return fmt.Errorf("target node %s not found", targetNodeName)
	}

	target.SetInput(inputSlot, source)
	return nil
}

// Pipeline manages the execution of a graph.
type Pipeline struct {
	Graph   *Graph
	Context context.Context
}

// NewPipeline creates a new pipeline for a given graph and context.
func NewPipeline(ctx context.Context, graph *Graph) *Pipeline {
	return &Pipeline{
		Graph:   graph,
		Context: ctx,
	}
}

// Execute runs the pipeline for a specific output node.
// It recursively processes dependencies.
func (p *Pipeline) Execute(outputNodeName string) error {
	node, ok := p.Graph.Nodes[outputNodeName]
	if !ok {
		return fmt.Errorf("output node %s not found", outputNodeName)
	}

	return node.Process(p.Context)
}

// Release releases all resources in the graph.
func (p *Pipeline) Release() {
	for _, node := range p.Graph.Nodes {
		node.Release()
	}
}
