package stack

type NodeType int

const (
	_              NodeType = iota
	NodeTypeObject          = 1
	NodeTypeArray           = 2
)

func (t NodeType) String() string {
	switch t {
	case NodeTypeObject:
		return "object"
	case NodeTypeArray:
		return "array"
	}
	return ""
}

type Node struct {
	Type             NodeType
	ArrLen           int
	KeyStart, KeyEnd int
}

// Stack is a LIFO Stack
type Stack struct {
	items []*Node
	top   int
}

// New creates a new stack instance.
// New will preallocate at least 16.
func New(prealloc int) *Stack {
	if prealloc < 16 {
		prealloc = 16
	}
	i := make([]*Node, prealloc)
	for x := range i {
		i[x] = &Node{}
	}
	return &Stack{
		items: i,
		top:   -1,
	}
}

// Reset resets the stack
func (s *Stack) Reset() {
	s.top = -1
}

// Len returns the current stack length
func (s *Stack) Len() int {
	return s.top + 1
}

// Push pushes a new node onto the stack
func (s *Stack) Push(t NodeType, ai, keyS, keyE int) {
	newTop := s.top + 1
	if newTop < len(s.items) {
		s.items[newTop].Type = t
		s.items[newTop].ArrLen = ai
		s.items[newTop].KeyStart = keyS
		s.items[newTop].KeyEnd = keyE
	} else {
		s.items = append(s.items, &Node{Type: t, ArrLen: ai})
	}
	s.top = newTop
}

// Pop drops the last pushed node.
// Pop does nothing if Len is 0
func (s *Stack) Pop() {
	if s.top < 0 {
		return
	}
	s.top--
}

// Top returns the last pushed node
func (s *Stack) Top() *Node {
	if s.top < 0 {
		return nil
	}
	return s.items[s.top]
}

// TopOffset returns the n'th last pushed node
func (s *Stack) TopOffset(offset int) *Node {
	i := s.top - offset
	if i < 0 {
		return nil
	}
	return s.items[i]
}
