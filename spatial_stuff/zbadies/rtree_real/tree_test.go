package rtree

import (
	"testing"
)

func TestSearch(t *testing.T) {
	min := 15 // After a split each node will contain at least min entries.
	max := 30 // A node will be split if it contains more than max entries.
	tree := New(min, max)
}
