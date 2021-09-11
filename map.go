// Copyright 2021 冯立强 mr.fengliqiang@gmail.com.  All rights reserved.

package goinline

// MapIterator
type MapIterator struct {
	node *RBTnode
}

// IsEnd
//go:nosplit
func (it MapIterator) IsEnd() bool {
	return it.node == nil
}

// Next : next Iterator
//go:nosplit
func (it MapIterator) Next() MapIterator {
	if it.node != nil {
		return MapIterator{it.node.Next()}
	}
	return MapIterator{nil}
}

// Pre : pre Iterator
//go:nosplit
func (it MapIterator) Pre() MapIterator {
	if it.node != nil {
		return MapIterator{it.node.Pre()}
	}
	return MapIterator{nil}
}

// Value return node data
func (it MapIterator) Value() *RBTpaire {
	if it.node == nil {
		return nil
	}
	return it.node.Get()
}

// Map by rbtree
type Map struct {
	tree RBtree
	size uint64
}

// Size : items count
//go:nosplit
func (m *Map) Size() uint64 {
	return m.size
}

// Init is the map constructor
// compaire:
//         if a < b :=>    ret < 0
//         if a > b :=>    ret > 0
//         if a == b :=>   ret = 0
//go:nosplit
func (m *Map) Init(compaire func(a, b interface{}) int) *Map {
	m.tree.Init(compaire)
	m.size = 0
	return m
}

// Clear
func (m *Map) Clear() {
	for m.size > 0 {
		m.Erase(m.Begin())
	}
}

// Begin
func (m *Map) Begin() MapIterator {
	if m.size > 0 {
		return MapIterator{m.tree.Begin()}
	}
	return MapIterator{nil}
}

// Rbegin : right begin
func (m *Map) Rbegin() MapIterator {
	if m.size > 0 {
		return MapIterator{m.tree.Rbegin()}
	}
	return MapIterator{nil}
}

// End
//go:nosplit
func (m *Map) End() MapIterator {
	return MapIterator{nil}
}

// Erase
func (m *Map) Erase(it MapIterator) {
	if it.node != nil && it.node.valid {
		m.tree.Remove(it.node)
		m.size--
	}
}

// Remove
func (m *Map) Remove(key interface{}) {
	m.Erase(m.Find(key))
}

// Set
func (m *Map) Set(key, value interface{}) MapIterator {
	isparent, node := m.tree.Find(key)
	if !isparent && node != nil {
		node.Get().Value = value
		return MapIterator{node}
	}
	newnode := (&RBTnode{}).init(colorRed)
	newnode.Value.first = key
	newnode.Value.Value = value
	m.tree.Insert(node, newnode)
	m.size++
	return MapIterator{newnode}
}

// Find
func (m *Map) Find(key interface{}) MapIterator {
	if isparent, node := m.tree.Find(key); (!isparent) && (node != nil) {
		return MapIterator{node}
	}
	return MapIterator{nil}
}
