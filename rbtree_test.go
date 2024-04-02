package rbtree

import (
	"cmp"
	"math/rand"
	"testing"
)

func TestRedBlackTreeInserts(t *testing.T) {
	t.Run("InsertOrdered", func(t *testing.T) {
		tree := New(cmp.Compare[int])
		inserts := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		for _, v := range inserts {
			out := tree.Insert(v)
			if out.Value != v {
				t.Errorf("got: %d, want: %d", out.Value, v)
			}
			isRedBlackTree(t, tree, tree.root)
		}
	})

	t.Run("InsertReverseOrdered", func(t *testing.T) {
		tree := New(cmp.Compare[int])
		inserts := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
		for _, v := range inserts {
			out := tree.Insert(v)
			if out.Value != v {
				t.Errorf("got: %d, want: %d", out.Value, v)
			}
			isRedBlackTree(t, tree, tree.root)
		}
	})

	t.Run("InsertShuffled", func(t *testing.T) {
		tree := New(cmp.Compare[int])
		inserts := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		rand.Shuffle(len(inserts), func(i, j int) { inserts[i], inserts[j] = inserts[j], inserts[i] })
		for _, v := range inserts {
			out := tree.Insert(v)
			if out.Value != v {
				t.Errorf("got: %d, want: %d", out.Value, v)
			}
			isRedBlackTree(t, tree, tree.root)
		}
	})
}

func TestRedBlackTreeDeletes(t *testing.T) {
	t.Run("DeleteOrdered", func(t *testing.T) {
		tree := New(cmp.Compare[int])
		inserts := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		for _, v := range inserts {
			tree.Insert(v)
			isRedBlackTree(t, tree, tree.root)
		}
		for _, ins := range inserts {
			if !tree.Delete(ins) {
				t.Errorf("failed to delete: %d", ins)
			}
			isRedBlackTree(t, tree, tree.root)
		}
	})

	t.Run("DeleteReverseOrdered", func(t *testing.T) {
		tree := New(cmp.Compare[int])
		inserts := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
		for _, v := range inserts {
			tree.Insert(v)
			isRedBlackTree(t, tree, tree.root)
		}
		for _, v := range inserts {
			if !tree.Delete(v) {
				t.Errorf("failed to delete: %d", v)
			}
			isRedBlackTree(t, tree, tree.root)
		}
	})

	t.Run("DeleteShuffled", func(t *testing.T) {
		tree := New(cmp.Compare[int])
		inserts := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		rand.Shuffle(len(inserts), func(i, j int) { inserts[i], inserts[j] = inserts[j], inserts[i] })
		for _, v := range inserts {
			tree.Insert(v)
			isRedBlackTree(t, tree, tree.root)
		}
		for _, v := range inserts {
			if !tree.Delete(v) {
				t.Errorf("failed to delete: %d", v)
			}
			isRedBlackTree(t, tree, tree.root)
		}
	})

	t.Run("DeleteHalf", func(t *testing.T) {
		tree := New(cmp.Compare[int])
		inserts := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		rand.Shuffle(len(inserts), func(i, j int) { inserts[i], inserts[j] = inserts[j], inserts[i] })
		for _, ins := range inserts[:len(inserts)/2] {
			tree.Insert(ins)
			isRedBlackTree(t, tree, tree.root)
		}
		for _, v := range inserts[:len(inserts)/2] {
			if !tree.Delete(v) {
				t.Errorf("failed to delete: %d", v)
			}
			isRedBlackTree(t, tree, tree.root)
		}
		for _, v := range inserts[len(inserts)/2:] {
			tree.Insert(v)
			isRedBlackTree(t, tree, tree.root)
		}
		for _, v := range inserts[len(inserts)/2:] {
			if !tree.Delete(v) {
				t.Errorf("failed to delete: %d", v)
			}
			isRedBlackTree(t, tree, tree.root)
		}
	})
}

// isRedBlackTree checks if a tree satisfies the red-black properties.
func isRedBlackTree[T any](t *testing.T, tree *RBTree[T], n *Node[T]) bool {
	t.Helper()

	if n == nil {
		// Base case: NIL nodes are black
		return true
	}

	// Check Red property (Red nodes have only black children)
	if n.color == red {
		if (n.left != nil && n.left.color == red) || (n.right != nil && n.right.color == red) {
			t.Errorf("Red node with red child detected:\n%s", tree)
			return false
		}
	}

	// Check for consistent black height and recurse
	leftBlackHeight := blackHeight(t, n.left)
	rightBlackHeight := blackHeight(t, n.right)
	if leftBlackHeight != rightBlackHeight || leftBlackHeight == -1 {
		t.Errorf("Black height inconsistent:\n%s", tree)
		return false
	}
	return isRedBlackTree(t, tree, n.left) && isRedBlackTree(t, tree, n.right)
}

// blackHeight calculates the black height of a node.
func blackHeight[T any](t *testing.T, n *Node[T]) int {
	if n == nil {
		return 1
	}
	leftHeight := blackHeight(t, n.left)
	if leftHeight == -1 {
		return -1
	}
	rightHeight := blackHeight(t, n.right)
	if rightHeight == -1 {
		return -1
	}
	if leftHeight != rightHeight {
		return -1
	}
	if n.color == black {
		return leftHeight + 1
	}
	return leftHeight
}

func TestRedBlackTreeSteps(t *testing.T) {
	tree := New(cmp.Compare[int])
	inserts := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, i := range inserts {
		tree.Insert(i)
	}

	five := tree.Search(5)
	t.Run("Predecessor", func(t *testing.T) {
		current := five
		for i := 3; i >= 0; i-- {
			current = tree.Predecessor(current)
			if current.Value != inserts[i] {
				t.Errorf("want: %d got: %d", inserts[i], current.Value)
			}
		}
	})
	t.Run("Successor", func(t *testing.T) {
		current := five
		for i := 5; i < len(inserts); i++ {
			current = tree.Successor(current)
			if current.Value != inserts[i] {
				t.Errorf("want: %d got: %d", inserts[i], current.Value)
			}
		}
	})
	t.Run("Edges", func(t *testing.T) {
		one := tree.Search(1)
		if tree.Predecessor(one) != nil {
			t.Error("we don't want this to be t.nil")
		}
		ten := tree.Search(10)
		if tree.Successor(ten) != nil {
			t.Error("we don't want this to be t.nil")
		}
	})
	t.Run("Wiggle", func(t *testing.T) {
		current := five
		current = tree.Predecessor(current)
		if current.Value != 4 {
			t.Errorf("want: %d got: %d", 4, current.Value)
		}
		current = tree.Successor(current)
		if current.Value != five.Value {
			t.Errorf("want: %d got: %d", five.Value, current.Value)
		}

		current = five
		current = tree.Successor(current)
		if current.Value != 6 {
			t.Errorf("want: %d got: %d", 6, current.Value)
		}
		current = tree.Predecessor(current)
		if current.Value != five.Value {
			t.Errorf("want: %d got: %d", five.Value, current.Value)
		}
	})
}
