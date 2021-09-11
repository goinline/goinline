// Copyright 2021 冯立强 mr.fengliqiang@gmail.com.  All rights reserved.

package goinline

type rbColor int8

const (
	colorRed   rbColor = 0
	colorBlack rbColor = 1
)

// RBTpaire : k-v paire
type RBTpaire struct {
	first interface{}
	Value interface{}
}

// Key
func (p *RBTpaire) Key() interface{} {
	return p.first
}

// RBTnode
type RBTnode struct {
	Value  RBTpaire
	left   *RBTnode
	right  *RBTnode
	parent *RBTnode
	color  rbColor
	valid  bool
}

// Get k-v Paire
//go:nosplit
func (t *RBTnode) Get() *RBTpaire {
	if t == nil {
		return nil
	}
	return &t.Value
}

//go:nosplit
func (t *RBTnode) isblack() bool {
	return t.color != colorRed
}

//go:nosplit
func (t *RBTnode) isred() bool {
	return t.color == colorRed
}

//go:nosplit
func (t *RBTnode) setblack() {
	t.color = colorBlack
}

//go:nosplit
func (t *RBTnode) setred() {
	t.color = colorRed
}

//go:nosplit
func (t *RBTnode) init(color rbColor) *RBTnode {
	t.left = nil
	t.right = nil
	t.parent = nil
	t.color = color
	t.valid = true
	return t
}

//go:nosplit
func (t *RBTnode) copycolor(other *RBTnode) {
	t.color = other.color
}

//go:nosplit
func (t *RBTnode) swapcolor(other *RBTnode) {
	_color := other.color
	other.color = t.color
	t.color = _color
}

// Pre node
//go:nosplit
func (t *RBTnode) Pre() *RBTnode {
	ret := t.left
	if ret != nil {
		for ret.right != nil {
			ret = ret.right
		}
		return ret
	}
	top := t
	for ret = t.parent; ret != nil; {
		if ret.right == top {
			return ret
		}
		top = ret
		ret = top.parent
	}
	return nil
}

// Next node
//go:nosplit
func (t *RBTnode) Next() *RBTnode {
	top := t
	ret := t.right
	if ret != nil {
		for ret.left != nil {
			ret = ret.left
		}
		return ret
	}

	for ret = top.parent; ret != nil; {
		if ret.left == top {
			return ret
		}
		top = ret
		ret = top.parent
	}
	return nil
}

// RBtree : Red/Black tree
type RBtree struct {
	compaire func(a, b interface{}) int
	root     *RBTnode
}

// Init struct
//go:nosplit
func (rbt *RBtree) Init(compaire func(a, b interface{}) int) *RBtree {
	rbt.root = nil
	// 他山之石
	rbt.compaire = compaire
	return rbt
}

//go:nosplit
func (rbt *RBtree) rotate(current *RBTnode) {
	for {
		parent := current.parent
		if parent == nil {
			current.setblack()
			break
		}
		if parent.isblack() {
			break
		}
		grandparent := parent.parent
		uncle := grandparent.left
		if uncle == parent {
			uncle = grandparent.right
		}
		if uncle != nil && uncle.isred() {
			parent.setblack()
			uncle.setblack()
			current = grandparent
			current.setred()
			continue
		}
		if grandparent.left == parent && parent.right == current {
			grandparent.left = current
			current.parent = grandparent
			parent.parent = current
			parent.right = current.left
			if parent.right != nil {
				parent.right.parent = parent
			}
			current.left = parent
			current = parent
			continue
		}
		if grandparent.right == parent && parent.left == current {
			grandparent.right = current
			current.parent = grandparent
			parent.parent = current
			parent.left = current.right
			if parent.left != nil {
				parent.left.parent = parent
			}
			current.right = parent
			current = parent
			continue
		}
		gg := grandparent.parent
		if gg == nil {
			rbt.root = parent
		} else if gg.left == grandparent {
			gg.left = parent
		} else {
			gg.right = parent
		}
		parent.parent = gg
		parent.setblack()
		grandparent.setred()
		if current == parent.left {
			grandparent.left = parent.right
			if parent.right != nil {
				parent.right.parent = grandparent
			}
			grandparent.parent = parent
			parent.right = grandparent
		} else {
			grandparent.right = parent.left
			if parent.left != nil {
				parent.left.parent = grandparent
			}
			grandparent.parent = parent
			parent.left = grandparent
		}
		break
	}
}

//go:nosplit
func replaceparent(x, a, b *RBTnode) {
	if x != nil {
		if x.left == a {
			x.left = b
		} else {
			x.right = b
		}
	}
}

//go:nosplit
func swapnode(a, b *RBTnode) {
	x := a.parent
	a.parent = b.parent
	replaceparent(x, a, b)
	b.parent = x
	replaceparent(a.parent, b, a)
	x = a.left
	a.left = b.left
	b.left = x
	if x != nil {
		x.parent = b
	}
	if a.left != nil {
		a.left.parent = a
	}
	x = a.right
	a.right = b.right
	b.right = x
	if x != nil {
		x.parent = b
	}
	if a.right != nil {
		a.right.parent = a
	}
	a.swapcolor(b)
}

func (rbt *RBtree) removeone(node *RBTnode) {
	child := node.left
	if child == nil {
		child = node.right
	}
	parent := node.parent
	if parent == nil {
		if child == nil {
			if node.isblack() {
				rbt.root = nil
			}
		} else {
			child.setblack()
			child.parent = nil
			rbt.root = child
		}
		return
	}
	var _tempnode *RBTnode = nil

	if node.isblack() && child == nil {
		_tempnode = (&RBTnode{}).init(colorBlack)
		child = _tempnode // 借鸡生蛋
	}
	if parent.left == node {
		parent.left = child
	} else {
		parent.right = child
	}
	if node.isred() {
		return
	}
	child.parent = parent
	if node.isblack() {
		if child.isred() {
			child.setblack()
		} else {
			rbt.removeCaseN(child)
		}
	}
	if child == _tempnode { // 好借好还
		parent = child.parent
		if parent.left == child {
			parent.left = nil
		} else {
			parent.right = nil
		}
		child.parent = nil
	}
}

//go:nosplit
func (rbt *RBtree) rotateleft(node *RBTnode) {
	parent := node.parent
	s := node.right
	s.parent = parent
	node.right = s.left
	if node.right != nil {
		node.right.parent = node
	}
	node.parent = s
	s.left = node
	if parent != nil {
		if parent.left == node {
			parent.left = s
		} else {
			parent.right = s
		}
	} else {
		rbt.root = s
	}
}

//go:nosplit
func (rbt *RBtree) rotateright(node *RBTnode) {
	parent := node.parent
	s := node.left
	s.parent = parent
	node.left = s.right
	if node.left != nil {
		node.left.parent = node
	}
	node.parent = s
	s.right = node
	if parent != nil {
		if parent.left == node {
			parent.left = s
		} else {
			parent.right = s
		}
	} else {
		rbt.root = s
	}
}

//go:nosplit
func isblack(p *RBTnode) bool {
	if p == nil {
		return true
	}
	return p.isblack()
}

func (rbt *RBtree) removeCaseN(n *RBTnode) {
	parent := n.parent
	var s *RBTnode = nil
	for {
		if parent == nil {
			break
		}
		if parent.left == n {
			s = parent.right
		} else {
			s = parent.left
		}
		if s.isred() {
			parent.setred()
			s.setblack()
			if parent.left == n {
				rbt.rotateleft(parent)
			} else {
				rbt.rotateright(parent)
			}
			if parent.left == n {
				s = parent.right
			} else {
				s = parent.left
			}
		}
		if parent.isblack() && s.isblack() && isblack(s.left) && isblack(s.right) {
			s.setred()
			n = parent
			parent = n.parent
			continue
		}
		if parent.isred() && s.isblack() && isblack(s.left) && isblack(s.right) {
			s.setred()
			parent.setblack()
			break
		}
		if s.isblack() {
			if n == parent.left && isblack(s.right) && !isblack(s.left) {
				s.setred()
				s.left.setblack()
				rbt.rotateright(s)
			} else if n == parent.right && isblack(s.left) && !isblack(s.right) {
				s.setred()
				s.right.setblack()
				rbt.rotateleft(s)
			}
			if parent.left == n {
				s = parent.right
			} else {
				s = parent.left
			}
		}
		s.copycolor(parent)
		parent.setblack()
		if n == parent.left {
			s.right.setblack()
			rbt.rotateleft(parent)
		} else {
			s.left.setblack()
			rbt.rotateright(parent)
		}
		break
	}
}

// Find : Return target node, Or parent of new node by this Key
func (rbt *RBtree) Find(key interface{}) (isParent bool, r *RBTnode) {
	node := rbt.root
	if node != nil {
		for {
			// 可以攻玉
			cmp := rbt.compaire(key, node.Value.first)
			if cmp == 0 {
				return false, node
			}
			if cmp < 0 {
				if node.left == nil {
					break
				}
				node = node.left
			} else {
				if node.right == nil {
					break
				}
				node = node.right
			}
		}
	}
	return true, node
}

// Insert : After Find, put new node on parent sub
func (rbt *RBtree) Insert(parent, node *RBTnode) {
	node.parent = parent
	node.left = nil
	node.right = nil
	node.valid = true
	if parent != nil {
		node.setred()
	} else {
		node.setblack()
	}
	if parent == nil {
		rbt.root = node
	} else {
		if rbt.compaire(node.Value.first, parent.Value.first) < 0 {
			parent.left = node
		} else {
			parent.right = node
		}
		rbt.rotate(node)
	}
}

// Remove node
func (rbt *RBtree) Remove(node *RBTnode) {
	if !node.valid {
		return
	}
	if node.right != nil {
		next := node.right
		for next.left != nil {
			next = next.left
		}
		swapnode(node, next)
		if next.parent == nil {
			rbt.root = next
		}
	}
	rbt.removeone(node)
	node.valid = false
}

// Begin : Node on tree left
//go:nosplit
func (rbt *RBtree) Begin() *RBTnode {
	ret := rbt.root
	if ret != nil {
		for ret.left != nil {
			ret = ret.left
		}
	}
	return ret
}

// Rbegin : Node on tree right
//go:nosplit
func (rbt *RBtree) Rbegin() *RBTnode {
	ret := rbt.root
	if ret != nil {
		for ret.right != nil {
			ret = ret.right
		}
	}
	return ret
}
