package tmproofs

import (
	"fmt"
	"math/bits"

	ics23 "github.com/confio/ics23/go"
	"github.com/tendermint/tendermint/crypto/merkle"
)

// ConvertExistenceProof will convert the given proof into a valid
// existence proof, if that's what it is.
//
// This is the simplest case of the range proof and we will focus on
// demoing compatibility here
func ConvertExistenceProof(p *merkle.Proof, key, value []byte) (*ics23.ExistenceProof, error) {
	path, err := convertInnerOps(p)
	if err != nil {
		return nil, err
	}

	proof := &ics23.ExistenceProof{
		Key:   key,
		Value: value,
		Leaf:  convertLeafOp(),
		Path:  path,
	}
	return proof, nil
}

// this is adapted from merkle/hash.go:leafHash()
// and merkle/simple_map.go:KVPair.Bytes()
func convertLeafOp() *ics23.LeafOp {
	prefix := []byte{0}

	return &ics23.LeafOp{
		Hash:         ics23.HashOp_SHA256,
		PrehashKey:   ics23.HashOp_NO_HASH,
		PrehashValue: ics23.HashOp_SHA256,
		Length:       ics23.LengthOp_VAR_PROTO,
		Prefix:       prefix,
	}
}

func convertInnerOps(p *merkle.Proof) ([]*ics23.InnerOp, error) {
	var inners []*ics23.InnerOp
	path := buildPath(p.Index, p.Total)

	if len(p.Aunts) != len(path) {
		return nil, fmt.Errorf("Calculated a path different length (%d) than provided by SimpleProof (%d)", len(path), len(p.Aunts))
	}

	for i, aunt := range p.Aunts {
		auntRight := path[i]

		// combine with: 0x01 || lefthash || righthash
		inner := &ics23.InnerOp{Hash: ics23.HashOp_SHA256}
		if auntRight {
			inner.Prefix = []byte{1}
			inner.Suffix = aunt
		} else {
			inner.Prefix = append([]byte{1}, aunt...)
		}
		inners = append(inners, inner)
	}
	return inners, nil
}

// buildPath returns a list of steps from leaf to root
// in each step, true means index is left side, false index is right side
// code adapted from merkle/simple_proof.go:computeHashFromAunts
func buildPath(idx int, total int) []bool {
	if total < 2 {
		return nil
	}
	numLeft := getSplitPoint(total)
	goLeft := idx < numLeft

	// we put goLeft at the end of the array, as we recurse from top to bottom,
	// and want the leaf to be first in array, root last
	if goLeft {
		return append(buildPath(idx, numLeft), goLeft)
	}
	return append(buildPath(idx-numLeft, total-numLeft), goLeft)
}

func getSplitPoint(length int) int {
	if length < 1 {
		panic("Trying to split a tree with size < 1")
	}
	uLength := uint(length)
	bitlen := bits.Len(uLength)
	k := 1 << uint(bitlen-1)
	if k == length {
		k >>= 1
	}
	return k
}
