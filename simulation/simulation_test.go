package simulation

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// GenerateItemsGraphviz creates a graphviz representation of the items. One can
// generate a graphical representation with `dot -Tpdf graph.dot -o graph.pdf`
func TestSimulation(t *testing.T) {
	var fileText, err = os.Create("outputSimulation")
	require.Equal(t, err, nil, "Cannot create output file for graph viz")
	buffer := new(bytes.Buffer)
	Simulation(buffer)
	fileText.Write(buffer.Bytes())
}
