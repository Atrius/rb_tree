package rb_tree

import (
	"math/rand"
	"testing"
)

func TestRbTreeBasicRight(t *testing.T) {
	r := RbTree{nil}
	for i := 0; i < 10; i++ {
		r.Add(i)
		validateTree(&r, t)
	}

	for i := 0; i < 10; i++ {
		r.Remove(i)
		validateTree(&r, t)
	}
}

func TestRbTreeBasicLeft(t *testing.T) {
	r := RbTree{nil}
	for i := 0; i < 10; i++ {
		r.Add(9 - i)
		validateTree(&r, t)
	}

	for i := 0; i < 10; i++ {
		r.Remove(9 - i)
		validateTree(&r, t)
	}
}

func TestRbTreeInverseRight(t *testing.T) {
	r := RbTree{nil}
	for i := 0; i < 10; i++ {
		r.Add(i)
		validateTree(&r, t)
	}

	for i := 0; i < 10; i++ {
		r.Remove(9 - i)
		validateTree(&r, t)
	}
}

func TestRbTreeInverseLeft(t *testing.T) {
	r := RbTree{nil}
	for i := 0; i < 10; i++ {
		r.Add(9 - i)
		validateTree(&r, t)
	}

	for i := 0; i < 10; i++ {
		r.Remove(i)
		validateTree(&r, t)
	}
}

func TestRbTreeInner(t *testing.T) {
	r := RbTree{nil}
	// Adding nodes to the "middle" of the tree covers add case 4:  n's parent
	// is red but uncle is black, and n is the "inner" child of p.
	for i := 0; i < 5; i++ {
		r.Add(i)
		validateTree(&r, t)
		r.Add(9 - i)
		validateTree(&r, t)
	}

	// Removing nodes the same way we added them covers delete case 5's first
	// branch:  n is the left child of p, s's left child is red and right child
	// is black.
	for i := 0; i < 5; i++ {
		r.Remove(i)
		validateTree(&r, t)
		r.Remove(9 - i)
		validateTree(&r, t)
	}
}

func TestRbTreeInnerRemoveReverse(t *testing.T) {
	r := RbTree{nil}
	// Adding nodes to the "middle" of the tree covers add case 4:  n's parent
	// is red but uncle is black, and n is the "inner" child of p.
	for i := 0; i < 5; i++ {
		r.Add(i)
		validateTree(&r, t)
		r.Add(9 - i)
		validateTree(&r, t)
	}

	// Removing nodes the same way they were added, but starting with a higher
	// node covers delete case 5's second branch:  n is the right child of p,
	// s's right child is red and left child is black.
	for i := 0; i < 5; i++ {
		r.Remove(9 - i)
		validateTree(&r, t)
		r.Remove(i)
		validateTree(&r, t)
	}
}

func TestRbTreeOuter(t *testing.T) {
	r := RbTree{nil}
	for i := 0; i < 5; i++ {
		r.Add(5 - i)
		validateTree(&r, t)
		r.Add(5 + i)
		validateTree(&r, t)
	}

	// Removing nodes from the outside first covers the case where a black node
	// is deleted and replaced by a red child.
	// This circumvents the complex rebalancing on delete.
	for i := 0; i < 5; i++ {
		r.Remove(5 - i)
		validateTree(&r, t)
		r.Remove(5 + i)
		validateTree(&r, t)
	}
}

func TestRbTreeRemoveMissing(t *testing.T) {
	r := RbTree{nil}
	r.Add(0)
	for i := 0; i < 10; i++ {
		r.Remove(i)
		validateTree(&r, t)
	}
}

func BenchmarkAdd(b *testing.B) {
	r := RbTree{nil}
	for i := 0; i < b.N; i++ {
		r.Add(rand.Int())
	}
}

func BenchmarkRemove(b *testing.B) {
	r := RbTree{nil}
	vals := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		val := rand.Int()
		r.Add(val)
		vals[i] = val
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Remove(vals[i])
	}
}

func validateTree(tree *RbTree, t *testing.T) {
	if !tree.root.black() {
		t.Errorf("Root %p is not black", tree.root)
	}
	validateNode(tree.root, t)
}

func validateNode(n *node, t *testing.T) int {
	if n.parent != nil {
		validateHasChild(n.parent, n, t)
	}

	if n.leaf {
		if n.red() {
			t.Errorf("Leaf node %p is not black", n)
		}
		return 1 // No children to check.
	}

	// Validate children, value ordering, count of black nodes along paths.
	leftBlackCount := 0
	rightBlackCount := 0
	if n.left != nil {
		validateHasParent(n.left, n, t)
		if !n.left.leaf && n.less(n.val, n.left.val) {
			t.Errorf("Node %p (%v) has value < left child %p (%v)",
				n, n.val, n.left, n.left.val)
		}
		if n.red() && n.left.red() {
			t.Errorf("Red node %p has red left child %p", n, n.left)
		}
		leftBlackCount = validateNode(n.left, t)
	}
	if n.right != nil {
		validateHasParent(n.right, n, t)
		if !n.right.leaf && n.less(n.right.val, n.val) {
			t.Errorf("Node %p (%v) has value > right child %p (%v)",
				n, n.val, n.right, n.right.val)
		}
		if n.red() && n.right.red() {
			t.Errorf("Red node %p has red right child %p", n, n.right)
		}
		rightBlackCount = validateNode(n.right, t)
	}
	if leftBlackCount != rightBlackCount {
		t.Errorf(
			"Not all paths from %p have equal numbers of black nodes.  "+
				"Left (%p) = %v, right (%p) = %v",
			n, n.left, leftBlackCount, n.right, rightBlackCount)
	}

	if n.black() {
		return leftBlackCount + 1
	} else {
		return leftBlackCount
	}
}

func validateHasChild(n, child *node, t *testing.T) {
	if n.left != child && n.right != child {
		t.Errorf("Node %p does not have child %p", n, child)
	}
}

func validateHasParent(n, parent *node, t *testing.T) {
	if n.parent != parent {
		t.Errorf("Node %p does not have parent %p", n, parent)
	}
}
