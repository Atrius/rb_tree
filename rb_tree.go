package rb_tree

type node struct {
	val         int
	left, right *node
	parent      *node
	black, leaf bool
}

func (n *node) add(val int) {
	if n.leaf {
		n.leaf = false
		n.black = false
		n.val = val
		n.left = newLeaf(n)
		n.right = newLeaf(n)
		rebalanceAdd(n)
		return
	}
	if val < n.val {
		n.left.add(val)
	} else {
		n.right.add(val)
	}
}

func rebalanceAdd(n *node) {
	p := n.parent

	// Case 1:  n is the root.
	if p == nil {
		n.black = true
		return
	}

	// Case 2:  n's parent is black and n is red.  Nothing to do.
	if p.black {
		return
	}

	// If p was root, it would have been black, so g must exist.
	g := p.parent
	u := g.left // u may be a leaf, we don't care.
	if g.left == p {
		u = g.right
	}

	// Case 3:  Parent and uncle are red.
	if !u.black {
		// Turn p and u black, g red.
		p.black = true
		u.black = true
		g.black = false
		// Rebalance at g.
		rebalanceAdd(g)
		return // No further operations on n, done rebalancing this portion.
	}

	// Case 4:  Parent is red but uncle is black, inner case.
	if (p == g.left && n == p.right) || (p == g.right && n == p.left) {
		rotate(p, n)
		// Rotation has caused n, p to swap roles.
		n = p
		p = n.parent
	}

	// Case 5:  Parent is red but uncle is black, outer case.
	g.black = false
	p.black = true
	rotate(g, p)
}

func newLeaf(parent *node) *node {
	return &node{0, nil, nil, parent, true, true}
}

func (n *node) find(val int) *node {
	if n.leaf {
		return nil
	}
	if n.val == val {
		return n
	} else if val < n.val {
		return n.left.find(val)
	} else {
		return n.right.find(val)
	}
}

func (n *node) remove() {
	candidate := n
	if !n.left.leaf {
		candidate = n.left
		for !candidate.right.leaf {
			candidate = candidate.right
		}
	} else if !n.right.leaf {
		candidate = n.right
		for !candidate.left.leaf {
			// Unreachable due to rebalancing.
			// We will always have a left child before adding multiple nodes to
			// the right subtree, and select the left child as candidate above
			// instead.
			candidate = candidate.left
		}
	}

	n.val = candidate.val
	candidate.delete()
}

func (n *node) delete() {
	p := n.parent
	// Handle deleting the root.  Resets the root to a leaf node.
	if p == nil {
		// The root being deleted means no other candidate for removal
		// was found.  Hence, the tree is now empty.
		n.leaf = true
		n.left = nil
		n.right = nil
		// No rebalancing necessary when deleting the root.
		return
	}

	replacement := n.left
	if !n.right.leaf {
		// Unreachable due to rebalancing.
		// We always prefer going left first, then as far right as possible to
		// find the deletion candidate.  Because of rebalancing, any node
		// will get a non-leaf left child before multiple nodes to the right.
		// Then we will select the left child's rightmost descendant as the
		// deletion candidate.  That node is as far right as possible, so it's
		// right child is always a leaf.
		// If we only have one right child, it is selected as the deletion
		// candidate and must have two leaves as children.
		replacement = n.right
	}
	p.replaceChild(n, replacement)
	if n.black {
		// If we're removing a black node and replacing it with a red one, just
		// turn the replacement black and there's no rebalancing necessary.
		if !replacement.black {
			replacement.black = true
		} else {
			rebalanceDelete(replacement)
		}
	}
	// Removing a red node does not require rebalancing.
}

func (n *node) replaceChild(child, replacement *node) {
	if n.left == child {
		n.left = replacement
	} else if n.right == child {
		n.right = replacement
	} else {
		panic("invariant violation:  attempted to replaceChild non-child")
	}
	replacement.parent = n
}

func rebalanceDelete(n *node) {
	// Case 1:  n is root.  Turn it black.
	if n.parent == nil {
		n.black = true
		return
	}

	// Initialize context.
	p := n.parent
	s := p.left
	if n == p.left {
		s = p.right
	}

	// Case 2:  s is red.
	if !s.black {
		s.black = true
		p.black = false
		rotate(p, s)

		// Recompute s.
		s = p.left
		if n == p.left {
			s = p.right
		}
	}

	// Case 3:  p, s, and s's children are all black.
	if p.black && s.black && s.right.black && s.left.black {
		s.black = false
		rebalanceDelete(p)
		return // Nothing more to rebalance in this subtree.
	}

	// Case 4:  p is red, s and s's children are black.
	if !p.black && s.black && s.right.black && s.left.black {
		p.black = true
		s.black = false
		return // Nothing more to rebalance.
	}

	// Case 5:  s has one red child, inner case.
	if n == p.left && s.right.black && !s.left.black {
		rotate(s, s.left)
		s.black = false // Turn the old s red.
		s = p.right     // Relabel s.
		s.black = true  // Turn the new s black.
	} else if n == p.right && s.left.black && !s.right.black {
		rotate(s, s.right)
		s.black = false // Turn the old s red.
		s = p.left      // Relabel s.
		s.black = true  // Turn the new s black.
	}

	// Case 6:  s has one red child, outer, or two red children.
	// Turn s's outer child black.
	if n == p.left {
		s.right.black = true
	} else {
		s.left.black = true
	}

	// Rotate p to s.
	rotate(p, s)

	// Swap colors of s and p.  S is known to be black due to case 2.
	s.black = p.black
	p.black = true
}

func rotate(oldParent, newParent *node) {
	// Switch grandparent link to new child.
	grandparent := oldParent.parent
	if grandparent != nil {
		grandparent.replaceChild(oldParent, newParent)
	} else {
		newParent.parent = nil
	}

	// Detect which way to rotate.
	newChild := newParent.left
	if oldParent.left == newParent {
		newChild = newParent.right
	}

	// Rotate old parent and new parent's link.
	oldParent.replaceChild(newParent, newChild)
	newParent.replaceChild(newChild, oldParent)
}

type RbTree struct {
	root *node
}

func (t *RbTree) Add(val int) {
	if t.root == nil {
		t.root = newLeaf(nil)
	}

	t.root.add(val)
	t.root = t.root.findRoot()
}

func (t *RbTree) Remove(val int) {
	n := t.root.find(val)
	if n == nil {
		return
	}

	n.remove()
	t.root = t.root.findRoot()
}

func (n *node) findRoot() *node {
	if n.parent == nil {
		return n
	}
	return n.parent.findRoot()
}
