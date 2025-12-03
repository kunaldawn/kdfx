package fxnode

import (
	"fmt"
	"kdfx/pkg/fxcontext"
)

// fxGraph represents a collection of nodes and their connections.
type fxGraph struct {
	nodes map[string]FXNode
}

// NewFXGraph creates a new empty fxGraph.
func NewFXGraph() FXGraph {
	return &fxGraph{
		nodes: make(map[string]FXNode),
	}
}

// AddNode adds a node to the fxGraph with a unique name.
func (g *fxGraph) AddNode(name string, node FXNode) {
	g.nodes[name] = node
}

// Connect connects the output of sourceNode to the input slot of targetNode.
func (g *fxGraph) Connect(sourceNodeName, targetNodeName, inputSlot string) error {
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

func (g *fxGraph) GetNode(name string) FXNode {
	return g.nodes[name]
}

func (g *fxGraph) Release() {
	for _, node := range g.nodes {
		node.Release()
	}
}

// fxPipeline manages the execution of a fxGraph.
type fxPipeline struct {
	fxGraph   FXGraph
	context fxcontext.FXContext
}

// NewFXPipeline creates a new fxPipeline for a given fxGraph and fxcontext.
func NewFXPipeline(ctx fxcontext.FXContext, fxGraph FXGraph) FXPipeline {
	return &fxPipeline{
		fxGraph:   fxGraph,
		context: ctx,
	}
}

// Execute runs the fxPipeline for a specific output fxnode.
// It recursively processes dependencies.
func (p *fxPipeline) Execute(outputNodeName string) error {
	node := p.fxGraph.GetNode(outputNodeName)
	if node == nil {
		return fmt.Errorf("output node %s not found", outputNodeName)
	}

	return node.Process(p.context)
}

// Release releases all resources in the fxGraph.
// Release releases all resources in the fxGraph.
func (p *fxPipeline) Release() {
	if p.fxGraph != nil {
		p.fxGraph.Release()
	}
}
