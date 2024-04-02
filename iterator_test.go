package rbtree

import (
	"cmp"
	"slices"
	"testing"
)

// runIterator runs the GOEXPERIMENT=rangefunc style iterators
// without any dependency on the new package
func runIterator(fn func(func(int) bool)) []int {
	var out []int
	fn(func(val int) bool {
		out = append(out, val)
		return true
	})

	return out
}

func TestIterators(t *testing.T) {
	tree := New(cmp.Compare[int])
	inserts := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, i := range inserts {
		tree.Insert(i)
	}

	t.Run("InOrder", func(t *testing.T) {
		out := runIterator(tree.Iterate(InOrder))
		want := inserts
		if !slices.Equal(out, want) {
			t.Errorf("slices differ:\n%#v\n%#v", out, want)
		}
	})
	t.Run("PreOrder", func(t *testing.T) {
		out := runIterator(tree.Iterate(PreOrder))
		want := []int{4, 2, 1, 3, 6, 5, 8, 7, 9, 10}
		if !slices.Equal(out, want) {
			t.Errorf("slices differ:\n%#v\n%#v", out, want)
		}
	})
	t.Run("PostOrder", func(t *testing.T) {
		out := runIterator(tree.Iterate(PostOrder))
		want := []int{1, 3, 2, 5, 7, 10, 9, 8, 6, 4}
		if !slices.Equal(out, want) {
			t.Errorf("slices differ:\n%#v\n%#v", out, want)
		}
	})
	t.Run("LevelOrder", func(t *testing.T) {
		out := runIterator(tree.Iterate(LevelOrder))
		want := []int{4, 2, 6, 1, 3, 5, 8, 7, 9, 10}
		if !slices.Equal(out, want) {
			t.Errorf("slices differ:\n%#v\n%#v", out, want)
		}
	})
}
