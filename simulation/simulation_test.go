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
	Simulation_RealData_Yes_No(buffer)
	fileText.Write(buffer.Bytes())
}

/* func TestSimulationCandQV(t *testing.T) {
	var fileText, err = os.Create("outputSimulation")
	require.Equal(t, err, nil, "Cannot create output file for graph viz")
	var fileText_QV, err_QV = os.Create("outputSimulation_QV")
	require.Equal(t, err_QV, nil, "Cannot create output file for graph viz")
	buffer := new(bytes.Buffer)
	buffer_QV := new(bytes.Buffer)
	Simulation_candidats_QV(buffer, buffer_QV)
	fileText.Write(buffer.Bytes())
	fileText_QV.Write(buffer_QV.Bytes())
} */
