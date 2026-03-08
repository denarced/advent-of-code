package aoc2308

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestCountSteps(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	lines, err := inr.ReadPath("testdata/in.txt", inr.IncludeEmpty())
	req.NoError(err, "failed to read test data")
	req.Equal(2, CountSteps(lines))
}

func TestParseLines(t *testing.T) {
	shared.InitTestLogging(t)
	path, nodes := parseLines([]string{
		"LRL",
		"",
		"AAA = (BBB, CCC)",
	})
	req := require.New(t)
	req.Equal("LRL", path)
	req.NotNil(nodes, "nodes is nil")
	req.Equal(3, len(nodes), "there should be 3 nodes")
	for _, each := range []string{"AAA", "BBB", "CCC"} {
		nod := nodes[each]
		req.NotNilf(nod, "nodes should contain %s", each)
		req.Equalf(each, nod.name, "node name should match map key %s", each)
	}
	aaaNode := nodes["AAA"]
	bbbNode := nodes["BBB"]
	cccNode := nodes["CCC"]
	req.Equal(bbbNode, aaaNode.left, "AAA -> BBB")
	req.Equal(cccNode, aaaNode.right, "AAA -> CCC")

	req.Nil(bbbNode.left, "BBB left is nil")
	req.Nil(bbbNode.right, "BBB right is nil")

	req.Nil(cccNode.left, "CCC left is nil")
	req.Nil(cccNode.right, "CCC right is nil")
}
