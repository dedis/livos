package simulation

import (
	"bytes"
	"testing"
)

// GenerateItemsGraphviz creates a graphviz representation of the items. One can
// generate a graphical representation with `dot -Tpdf graph.dot -o graph.pdf`
func TestSimulation(t *testing.T) {
	buffer := new(bytes.Buffer)
	Simulation(buffer)
}
