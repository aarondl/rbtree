// Package rbtree implements a red-black tree capable of
// containing any type with generics and the compare function.
//
// The API allows for nodes to be held on to externally to
// encourage special indexing concerns, but care must be
// taken in those advanced sorts of used cases.
package rbtree

import (
	"fmt"
	"strings"
)

type color bool

const (
	black color = false
	red   color = true
)

// RBTree is a generic red black tree.
type RBTree[T any] struct {
	root    *Node[T]
	nil     *Node[T]
	compare func(a, b T) int
}

// New constructs a red black tree, note that compare can never return 0.
func New[T any](compare func(a, b T) int) *RBTree[T] {
	nil := &Node[T]{color: black}
	return &RBTree[T]{compare: compare, root: nil, nil: nil}
}

// Node for the red black tree, only exposes it's value publicly
// to lessen the possibility of accidental manipulation, granted
// the value is enough to make a mess.
type Node[T any] struct {
	color color

	parent *Node[T]
	left   *Node[T]
	right  *Node[T]

	Value T
}

func (n *Node[T]) getColor() color {
	if n == nil {
		return black
	}
	return n.color
}

// Insert val and return a pointer to the Node that was inserted
// for indexing purposes.
func (r *RBTree[T]) Insert(val T) *Node[T] {
	if r.root == r.nil {
		// recolor from red to black to avoid fixup call
		r.root = &Node[T]{
			Value: val,
			color: black,
			left:  r.nil,
			right: r.nil,
		}
		return r.root
	}

	current := r.root
	insert := &Node[T]{
		Value: val,
		color: red,
		left:  r.nil,
		right: r.nil,
	}
	for current != r.nil {
		test := r.compare(val, current.Value)
		if test < 0 {
			if current.left == r.nil {
				insert.parent = current
				current.left = insert
				break
			}
			current = current.left
		} else if test > 0 {
			if current.right == r.nil {
				insert.parent = current
				current.right = insert
				break
			}
			current = current.right
		} else {
			panic("duplicate value")
		}
	}

	r.insertFixup(insert)
	return insert
}

func (r *RBTree[T]) insertFixup(check *Node[T]) {
	for check.parent.getColor() == red {
		grandParent := check.parent.parent

		// Check direction of our parent (left or right)
		if check.parent == grandParent.left {
			uncle := grandParent.right   // uncle will be on the right
			if uncle.getColor() == red { // right uncle is red
				check.parent.color = black
				uncle.color = black
				grandParent.color = red
				check = grandParent
			} else {
				if check == check.parent.right { // right uncle black, triangle case
					check = check.parent
					r.rotateLeft(check)
				}
				// right uncle black, line case
				check.parent.color = black
				check.parent.parent.color = red
				r.rotateRight(check.parent.parent)
			}
		} else {
			uncle := grandParent.left    // uncle will be on the left
			if uncle.getColor() == red { // left uncle is red
				check.parent.color = black
				uncle.color = black
				grandParent.color = red
				check = grandParent
			} else {
				if check == check.parent.left { // left uncle black, triangle case
					check = check.parent
					r.rotateRight(check)
				}
				// left uncle black, line case
				check.parent.color = black
				check.parent.parent.color = red
				r.rotateLeft(check.parent.parent)
			}
		}
	}

	r.root.color = black
}

// Delete a value. This is the equivalent of DeleteNode(Search(val))
func (r *RBTree[T]) Delete(val T) bool {
	return r.DeleteNode(r.Search(val))
}

// DeleteNode deletes the provided node, this provides an easy
// way to delete a node that's been indexed outside of this
// data structure.
func (r *RBTree[T]) DeleteNode(n *Node[T]) bool {
	if n == nil {
		return false
	}

	var odd *Node[T]
	originalColor := n.color

	if n.left == r.nil {
		// case 1: left child nil
		odd = n.right
		r.transplant(n, odd)
	} else if n.right == r.nil {
		// case 2: right child nil
		odd = n.left
		r.transplant(n, odd)
	} else {
		// case 3: neither nil
		minimum := n.right
		for minimum.left != r.nil {
			minimum = minimum.left
		}

		originalColor = minimum.color
		odd = minimum.right

		if minimum.parent == n {
			if odd != nil {
				odd.parent = minimum
			}
		} else {
			r.transplant(minimum, minimum.right)
			minimum.right = n.right
			minimum.right.parent = minimum
		}

		r.transplant(n, minimum)
		minimum.left = n.left
		minimum.left.parent = minimum
		minimum.color = n.color
	}

	if originalColor == black {
		r.deleteFixup(odd)
	}

	return true
}

func (r *RBTree[T]) transplant(u, v *Node[T]) {
	if u.parent == nil {
		r.root = v
	} else if u == u.parent.left {
		u.parent.left = v
	} else {
		u.parent.right = v
	}

	v.parent = u.parent
}

func (r *RBTree[T]) deleteFixup(n *Node[T]) {
	for n != r.root && n.getColor() == black {
		if n == n.parent.left {
			sibling := n.parent.right

			// case 1: sibling is red
			if sibling.getColor() == red {
				sibling.color = black
				n.parent.color = red
				r.rotateLeft(n.parent)
				sibling = n.parent.right
			}

			// case 2: sibling has two black descendants
			if sibling.left.getColor() == black && sibling.right.getColor() == black {
				sibling.color = red
				n = n.parent
			} else {
				// case 3
				if sibling.right.getColor() == black {
					sibling.left.color = black
					sibling.color = red
					r.rotateRight(sibling)
					sibling = n.parent.right
				}

				// case 4
				sibling.color = n.parent.color
				n.parent.color = black
				sibling.right.color = black
				r.rotateLeft(n.parent)
				n = r.root
			}
		} else {
			sibling := n.parent.left

			// case 1: sibling is red
			if sibling.getColor() == red {
				sibling.color = black
				n.parent.color = red
				r.rotateRight(n.parent)
				sibling = n.parent.left
			}

			// case 2: sibling has two black descendants
			if sibling.right.getColor() == black && sibling.left.getColor() == black {
				sibling.color = red
				n = n.parent
			} else {
				// case 3
				if sibling.left.getColor() == black {
					sibling.right.color = black
					sibling.color = red
					r.rotateLeft(sibling)
					sibling = n.parent.left
				}

				// case 4
				sibling.color = n.parent.color
				n.parent.color = black
				sibling.left.color = black
				r.rotateRight(n.parent)
				n = r.root
			}
		}
	}
	n.color = black
}

// Search for a node in the tree, returns nil if not found
func (r *RBTree[T]) Search(val T) *Node[T] {
	current := r.root
	for current != r.nil {
		test := r.compare(val, current.Value)
		if test < 0 {
			current = current.left
		} else if test > 0 {
			current = current.right
		} else {
			return current
		}
	}

	return nil
}

// Has is a convenience method that is equivalent to Search(val) != nil
func (r *RBTree[T]) Has(val T) bool {
	return r.Search(val) != nil
}

func (r *RBTree[T]) rotateLeft(n *Node[T]) {
	if n.right == r.nil {
		panic("is this possible?")
	}

	// set all the descendants
	newRoot := n.right
	n.right = newRoot.left
	newRoot.left = n
	if n.right != r.nil {
		// fix the parent of the newly adopted descendant
		n.right.parent = n
	}

	// fix the parent of the old root
	if n.parent == nil {
		r.root = newRoot
	} else if n.parent.left == n {
		n.parent.left = newRoot
	} else {
		n.parent.right = newRoot
	}

	// fix new root parent
	newRoot.parent = n.parent

	// fix the old root parent
	n.parent = newRoot
}

func (r *RBTree[T]) rotateRight(n *Node[T]) {
	if n.left == nil {
		panic("is this possible?")
	}

	// set all the descendants
	newRoot := n.left
	n.left = newRoot.right
	newRoot.right = n
	if n.left != r.nil {
		// fix the parent of the newly adopted descendant
		n.left.parent = n
	}

	// fix the parent of the old root
	if n.parent == nil {
		r.root = newRoot
	} else if n.parent.right == n {
		n.parent.right = newRoot
	} else {
		n.parent.left = newRoot
	}

	// fix new root parent
	newRoot.parent = n.parent

	// fix the old root parent
	n.parent = newRoot
}

// Successor looks up the successor to the given node.
// Can be helpful in certain odd iteration scenarios.
// Returns nil if there is none.
func (t *RBTree[T]) Successor(node *Node[T]) *Node[T] {
	if node == nil {
		return nil
	}

	if node.right != t.nil {
		node = node.right
		for node.left != t.nil {
			node = node.left
		}
		return node
	}

	succ := node.parent
	for succ != nil && node == succ.right {
		node = succ
		succ = succ.parent
	}
	return succ
}

// Predecessor looks up the predecessor to the given node.
// Can be helpful in certain odd iteration scenarios.
// Returns nil if there is none.
func (t *RBTree[T]) Predecessor(node *Node[T]) *Node[T] {
	if node == nil {
		return nil
	}

	if node.left != t.nil {
		node = node.left
		for node.right != t.nil {
			node = node.right
		}
		return node
	}

	pred := node.parent
	for pred != nil && node == pred.left {
		node = pred
		pred = pred.parent
	}
	return pred
}

func (r *RBTree[T]) String() string {
	var builder strings.Builder
	builder.WriteString("digraph RBTree {\n")
	builder.WriteString("  node [shape = circle];\n")
	r.nodeToDot(r.root, &builder)
	builder.WriteString("}\n")
	return builder.String()
}

func (r *RBTree[T]) nodeToDot(n *Node[T], builder *strings.Builder) {
	if n == nil || n == r.nil {
		return
	}

	color := "black"
	if n.color == red {
		color = "red"
	}

	builder.WriteString(fmt.Sprintf("  \"%p\" [label = \"%v\", color=%s];\n", n, n.Value, color))

	if n.left != r.nil {
		builder.WriteString(fmt.Sprintf("  \"%p\" -> \"%p\";\n", n, n.left))
		r.nodeToDot(n.left, builder)
	}
	if n.right != r.nil {
		builder.WriteString(fmt.Sprintf("  \"%p\" -> \"%p\";\n", n, n.right))
		r.nodeToDot(n.right, builder)
	}
}
