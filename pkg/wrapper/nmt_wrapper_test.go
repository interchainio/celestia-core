package wrapper

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"sort"
	"testing"

	"github.com/celestiaorg/nmt"
	"github.com/celestiaorg/rsmt2d"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/pkg/consts"
)

func TestPushErasuredNamespacedMerkleTree(t *testing.T) {
	testCases := []struct {
		name       string
		squareSize int
	}{
		{"extendedSquareSize = 16", 8},
		{"extendedSquareSize = 256", 128},
	}
	for _, tc := range testCases {
		tc := tc
		n := NewErasuredNamespacedMerkleTree(uint64(tc.squareSize))
		tree := n.Constructor()

		// push test data to the tree
		for i, d := range generateErasuredData(t, tc.squareSize, consts.DefaultCodec()) {
			// push will panic if there's an error
			tree.Push(d, rsmt2d.SquareIndex{Axis: uint(0), Cell: uint(i)})
		}
	}
}

func TestRootErasuredNamespacedMerkleTree(t *testing.T) {
	// check that the root is different from a standard nmt tree this should be
	// the case, because the ErasuredNamespacedMerkleTree should add namespaces
	// to the second half of the tree
	size := 8
	data := generateRandNamespacedRawData(size, consts.NamespaceSize, consts.MsgShareSize)
	n := NewErasuredNamespacedMerkleTree(uint64(size))
	tree := n.Constructor()
	nmtTree := nmt.New(sha256.New())

	for i, d := range data {
		tree.Push(d, rsmt2d.SquareIndex{Axis: uint(0), Cell: uint(i)})
		err := nmtTree.Push(d)
		if err != nil {
			t.Error(err)
		}
	}

	assert.NotEqual(t, nmtTree.Root(), tree.Root())
}

func TestErasureNamespacedMerkleTreePanics(t *testing.T) {
	testCases := []struct {
		name  string
		pFunc assert.PanicTestFunc
	}{
		{
			"push over square size",
			assert.PanicTestFunc(
				func() {
					data := generateErasuredData(t, 16, consts.DefaultCodec())
					n := NewErasuredNamespacedMerkleTree(uint64(15))
					tree := n.Constructor()
					for i, d := range data {
						tree.Push(d, rsmt2d.SquareIndex{Axis: uint(0), Cell: uint(i)})
					}
				}),
		},
		{
			"push in incorrect lexigraphic order",
			assert.PanicTestFunc(
				func() {
					data := generateErasuredData(t, 16, consts.DefaultCodec())
					n := NewErasuredNamespacedMerkleTree(uint64(16))
					tree := n.Constructor()
					for i := len(data) - 1; i > 0; i-- {
						tree.Push(data[i], rsmt2d.SquareIndex{Axis: uint(0), Cell: uint(i)})
					}
				},
			),
		},
	}
	for _, tc := range testCases {
		tc := tc
		assert.Panics(t, tc.pFunc, tc.name)

	}
}

func TestExtendedDataSquare(t *testing.T) {
	squareSize := 4
	// data for a 4X4 square
	raw := generateRandNamespacedRawData(
		squareSize*squareSize,
		consts.NamespaceSize,
		consts.MsgShareSize,
	)

	tree := NewErasuredNamespacedMerkleTree(uint64(squareSize))

	_, err := rsmt2d.ComputeExtendedDataSquare(raw, consts.DefaultCodec(), tree.Constructor)
	assert.NoError(t, err)
}

func TestErasuredNamespacedMerkleTree(t *testing.T) {
	// check that the Tree() returns exact underlying nmt tree
	size := 8
	data := generateRandNamespacedRawData(size, consts.NamespaceSize, consts.MsgShareSize)
	n := NewErasuredNamespacedMerkleTree(uint64(size))
	tree := n.Constructor()

	for i, d := range data {
		tree.Push(d, rsmt2d.SquareIndex{Axis: uint(0), Cell: uint(i)})
	}

	assert.Equal(t, n.Tree(), n.tree)
	assert.Equal(t, n.Tree().Root(), n.tree.Root())
}

// generateErasuredData produces a slice that is twice as long as it erasures
// the data
func generateErasuredData(t *testing.T, numLeaves int, codec rsmt2d.Codec) [][]byte {
	raw := generateRandNamespacedRawData(
		numLeaves,
		consts.NamespaceSize,
		consts.MsgShareSize,
	)
	erasuredData, err := codec.Encode(raw)
	if err != nil {
		t.Error(err)
	}
	return append(raw, erasuredData...)
}

// this code is copy pasted from the plugin, and should likely be exported in the plugin instead
func generateRandNamespacedRawData(total int, nidSize int, leafSize int) [][]byte {
	data := make([][]byte, total)
	for i := 0; i < total; i++ {
		nid := make([]byte, nidSize)
		_, err := rand.Read(nid)
		if err != nil {
			panic(err)
		}
		data[i] = nid
	}

	sortByteArrays(data)
	for i := 0; i < total; i++ {
		d := make([]byte, leafSize)
		_, err := rand.Read(d)
		if err != nil {
			panic(err)
		}
		data[i] = append(data[i], d...)
	}

	return data
}

func sortByteArrays(src [][]byte) {
	sort.Slice(src, func(i, j int) bool { return bytes.Compare(src[i], src[j]) < 0 })
}
