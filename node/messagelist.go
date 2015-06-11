//
// Copy from GO library
//
package main

// MessageElement is an element of a linked list.
type MessageElement struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *MessageElement

	// The list to which this element belongs.
	list *MessageList

	// The value stored with this element.
	Value *Message
}

// Next returns the next list element or nil.
func (e *MessageElement) Next() *MessageElement {
	if p := e.next; p != &e.list.root {
		return p
	}
	return nil
}

// Prev returns the previous list element or nil.
func (e *MessageElement) Prev() *MessageElement {
	if p := e.prev; p != &e.list.root {
		return p
	}
	return nil
}

// MessageList represents a doubly linked list.
// The zero value for MessageList is an empty list ready to use.
type MessageList struct {
	root MessageElement // sentinel list element, only &root, root.prev, and root.next are used
	len  int            // current list length excluding (this) sentinel element
}

// Init initializes or clears list l.
func (l *MessageList) Init() *MessageList {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// New returns an initialized list.
func NewMessageList() *MessageList { return new(MessageList).Init() }

// Len returns the number of elements of list l.
func (l *MessageList) Len() int { return l.len }

// Front returns the first element of list l or nil
func (l *MessageList) Front() *MessageElement {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back returns the last element of list l or nil.
func (l *MessageList) Back() *MessageElement {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// lazyInit lazily initializes a zero MessageList value.
func (l *MessageList) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// insert inserts e after at, increments l.len, and returns e.
func (l *MessageList) insert(e, at *MessageElement) *MessageElement {
	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e
	e.list = l
	l.len++
	return e
}

// insertValue is a convenience wrapper for insert(&MessageElement{Value: v}, at).
func (l *MessageList) insertValue(v *Message, at *MessageElement) *MessageElement {
	return l.insert(&MessageElement{Value: v}, at)
}

// remove removes e from its list, decrements l.len, and returns e.
func (l *MessageList) remove(e *MessageElement) *MessageElement {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--
	return e
}

// Remove removes e from l if e is an element of list l.
// It returns the element value e.Value.
func (l *MessageList) Remove(e *MessageElement) *Message {
	if e.list == l {
		// if e.list == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero MessageElement) and l.remove will crash
		l.remove(e)
	}
	return e.Value
}

// Pushfront inserts a new element e with value v at the front of list l and returns e.
func (l *MessageList) PushFront(v *Message) *MessageElement {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

// PushBack inserts a new element e with value v at the back of list l and returns e.
func (l *MessageList) PushBack(v *Message) *MessageElement {
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}
