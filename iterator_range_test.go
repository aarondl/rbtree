//go:build iterator

package rbtree

import (
	"cmp"
	"slices"
	"testing"
)

// TestIteratorsRange tests the iterators with the range built-in
// and the GOEXPERIMENT=rangefunc iterator support.
//
// Use the build tag `iterator` to enable these tests.
func TestIteratorsRange(t *testing.T) {
	tree := New(cmp.Compare[int])
	inserts := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, i := range inserts {
		tree.Insert(i)
	}

	t.Run("InOrder", func(t *testing.T) {
		var out []int
		want := inserts
		for i := range tree.Iterate(InOrder) {
			out = append(out, i)
		}
		if !slices.Equal(out, want) {
			t.Errorf("slices differ:\n%#v\n%#v", out, want)
		}
	})
	t.Run("PreOrder", func(t *testing.T) {
		var out []int
		want := []int{4, 2, 1, 3, 6, 5, 8, 7, 9, 10}
		for i := range tree.Iterate(PreOrder) {
			out = append(out, i)
		}
		if !slices.Equal(out, want) {
			t.Errorf("slices differ:\n%#v\n%#v", out, want)
		}
	})
	t.Run("PostOrder", func(t *testing.T) {
		var out []int
		want := []int{1, 3, 2, 5, 7, 10, 9, 8, 6, 4}
		for i := range tree.Iterate(PostOrder) {
			out = append(out, i)
		}
		if !slices.Equal(out, want) {
			t.Errorf("slices differ:\n%#v\n%#v", out, want)
		}
	})
	t.Run("LevelOrder", func(t *testing.T) {
		var out []int
		want := []int{4, 2, 6, 1, 3, 5, 8, 7, 9, 10}
		for i := range tree.Iterate(LevelOrder) {
			out = append(out, i)
		}
		if !slices.Equal(out, want) {
			t.Errorf("slices differ:\n%#v\n%#v", out, want)
		}
	})
}
