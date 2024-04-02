package rbtree

type IterationMethod int

const (
	InOrder IterationMethod = iota
	PreOrder
	PostOrder
	LevelOrder
)

// Iterate over the collection with the desired iteration method.
//
// Notably this function signature is compatible with the
// GOEXPERIMENT=rangefunc iter.Seq[V any] iterator proposal.
// This means it can be used with the range built-in if
// the environment variable is set.
func (r *RBTree[T]) Iterate(method IterationMethod) func(func(T) bool) {
	switch method {
	case InOrder:
		iterator := inOrderIter[T]{tree: r}
		return iterator.Iterate
	case PreOrder:
		iterator := preOrderIter[T]{tree: r}
		return iterator.Iterate
	case PostOrder:
		iterator := postOrderIter[T]{tree: r}
		return iterator.Iterate
	case LevelOrder:
		iterator := levelOrderIter[T]{tree: r}
		return iterator.Iterate
	default:
		panic("unknown iteration method")
	}
}

type inOrderIter[T any] struct {
	tree  *RBTree[T]
	stack []*Node[T]
}

type preOrderIter[T any] struct {
	tree  *RBTree[T]
	stack []*Node[T]
}

type postOrderIter[T any] struct {
	tree      *RBTree[T]
	stack     []*Node[T]
	lastVisit *Node[T]
}

type levelOrderIter[T any] struct {
	tree  *RBTree[T]
	queue []*Node[T]
}

func (i *inOrderIter[T]) Iterate(yield func(T) bool) {
	current := i.tree.root

	for current != i.tree.nil || len(i.stack) > 0 {
		for current != i.tree.nil {
			i.stack = append(i.stack, current)
			current = current.left
		}
		current = i.stack[len(i.stack)-1]
		i.stack = i.stack[:len(i.stack)-1]

		if !yield(current.Value) {
			return
		}

		current = current.right
	}
}

func (i *preOrderIter[T]) Iterate(yield func(T) bool) {
	if i.tree.root == i.tree.nil {
		return
	}
	i.stack = []*Node[T]{i.tree.root}

	for len(i.stack) > 0 {
		node := i.stack[len(i.stack)-1]
		i.stack = i.stack[:len(i.stack)-1]

		if !yield(node.Value) {
			return
		}

		if node.right != i.tree.nil {
			i.stack = append(i.stack, node.right)
		}
		if node.left != i.tree.nil {
			i.stack = append(i.stack, node.left)
		}
	}
}

func (i *postOrderIter[T]) Iterate(yield func(T) bool) {
	if i.tree.root == i.tree.nil {
		return
	}
	current := i.tree.root

	for len(i.stack) > 0 || current != i.tree.nil {
		if current != i.tree.nil {
			i.stack = append(i.stack, current)
			current = current.left
		} else {
			peekNode := i.stack[len(i.stack)-1]
			if peekNode.right != i.tree.nil && i.lastVisit != peekNode.right {
				current = peekNode.right
			} else {
				if !yield(peekNode.Value) {
					return
				}
				i.lastVisit = i.stack[len(i.stack)-1]
				i.stack = i.stack[:len(i.stack)-1]
			}
		}
	}
}

func (i *levelOrderIter[T]) Iterate(yield func(T) bool) {
	if i.tree.root == i.tree.nil {
		return
	}
	i.queue = []*Node[T]{i.tree.root}

	for len(i.queue) > 0 {
		node := i.queue[0]
		i.queue = i.queue[1:]

		if !yield(node.Value) {
			return
		}

		if node.left != i.tree.nil {
			i.queue = append(i.queue, node.left)
		}
		if node.right != i.tree.nil {
			i.queue = append(i.queue, node.right)
		}
	}
}
