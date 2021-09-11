// Copyright 2021 冯立强 mr.fengliqiang@gmail.com.  All rights reserved.

package goinline

type listNode struct {
	pre   *listNode
	nxt   *listNode
	value interface{}
}

//go:nosplit
func (n *listNode) reset() {
	n.pre = nil
	n.nxt = nil
	n.value = nil
}

//Iterator : iterator of list
type ListIterator struct {
	root *List
	node *listNode
}

//Set : set the node value
func (i ListIterator) Set(value interface{}) bool {
	if i.node == nil {
		return false
	}
	i.node.value = value
	return true
}

//Valid : check iterator valid
//go:nosplit
func (i ListIterator) Valid() bool {
	if i.root == nil {
		return false
	}
	if i.node == nil {
		return true
	}
	if i.root.count == 0 {
		return false
	}
	if i.node.pre == nil {
		return i.root.first == i.node
	}
	if i.node.nxt == nil {
		return i.root.last == i.node
	}
	return i.node.pre.nxt == i.node && i.node.nxt.pre == i.node
}

//Destroy iterator
//go:nosplit
func (i ListIterator) Destroy() {
	i.node = nil
}

//Remove node from list
func (i ListIterator) Remove() (interface{}, bool) {
	if !i.Valid() || i.root.count == 0 {
		return nil, false
	}
	defer i.Destroy()
	if i.node.pre == nil {
		return i.root.PopFront()
	}
	if i.node.nxt == nil {
		return i.root.PopBack()
	}
	i.node.pre.nxt = i.node.nxt
	i.node.nxt.pre = i.node.pre
	i.root.count--
	ret := i.node.value
	i.node.reset()
	return ret, true
}

//InsertFront : insert a node front of current
func (i ListIterator) InsertFront(value interface{}) (ListIterator, bool) {
	if !i.Valid() || i.node == nil {
		return ListIterator{i.root, nil}, false
	}
	n := &listNode{i.node.pre, i.node, value}
	i.node.pre = n
	if n.pre == nil {
		i.root.first = n
	} else {
		n.pre.nxt = n
	}
	i.root.count++
	return ListIterator{i.root, n}, true
}

//InsertBack : insert a node back of current
func (i ListIterator) InsertBack(value interface{}) (ListIterator, bool) {
	if !i.Valid() || i.node == nil {
		return ListIterator{i.root, nil}, false
	}
	n := &listNode{i.node, i.node.nxt, value}
	i.node.nxt = n
	if n.nxt == nil {
		i.root.last = n
	} else {
		n.nxt.pre = n
	}
	i.root.count++
	return ListIterator{i.root, n}, true
}

//Value
//go:nosplit
func (i ListIterator) Value() (interface{}, bool) {

	if i.Valid() && i.node != nil {
		return i.node.value, true
	}
	return nil, false
}

//Back : next node
//go:nosplit
func (i ListIterator) Back() ListIterator {

	if i.Valid() && i.node != nil {
		return ListIterator{root: i.root, node: i.node.nxt}
	}
	return ListIterator{i.root, nil}
}

//Front
//go:nosplit
func (i ListIterator) Front() ListIterator {

	if i.Valid() && i.node != nil {
		return ListIterator{root: i.root, node: i.node.pre}
	}
	return ListIterator{i.root, nil}
}

//List : general list
type List struct {
	first *listNode
	last  *listNode
	count int
}

//Size
//go:nosplit
func (l *List) Size() int {
	return l.count
}

//Clear
func (l *List) Clear() {
	for l.Size() > 0 {
		l.PopFront()
	}
}
func (l *List) End() ListIterator {
	return ListIterator{l, nil}
}

//Front
//go:nosplit
func (l *List) Front() ListIterator {
	return ListIterator{l, l.first}
}

//Back
//go:nosplit
func (l *List) Back() ListIterator {
	return ListIterator{l, l.last}
}

//PushBack
func (l *List) PushBack(value interface{}) ListIterator {
	n := &listNode{l.last, nil, value}
	if l.count == 0 {
		l.first = n
	} else {
		l.last.nxt = n
	}
	l.last = n
	l.count++
	return ListIterator{l, n}
}

//PushFront
func (l *List) PushFront(value interface{}) ListIterator {
	n := &listNode{nil, l.first, value}
	if l.count == 0 {
		l.last = n
	} else {
		l.first.pre = n
	}
	l.first = n
	l.count++
	return ListIterator{l, n}
}

//PopFront
func (l *List) PopFront() (interface{}, bool) {
	if l.count == 0 {
		return nil, false
	}
	rnode := l.first
	ret := rnode.value
	l.first = rnode.nxt
	rnode.value = nil
	l.count--
	if l.count == 0 {
		l.last = nil
	} else {
		l.first.pre = nil
		rnode.nxt = nil
	}
	return ret, true
}

//PopBack
func (l *List) PopBack() (interface{}, bool) {
	if l.count == 0 {
		return nil, false
	}
	rnode := l.last
	ret := rnode.value
	l.last = rnode.pre
	rnode.value = nil
	l.count--
	if l.count == 0 {
		l.first = nil
	} else {
		l.last.nxt = nil
		rnode.pre = nil
	}
	return ret, true
}
