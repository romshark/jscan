package stack_test

import (
	"testing"

	"github.com/romshark/jscan/internal/stack"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	r := require.New(t)
	s := stack.New(0)

	r.Zero(s.Top())
	r.Zero(s.TopOffset(0))
	r.Zero(s.Len())

	s.Push(42, 0, 0, 0)
	r.Equal(&stack.Node{Type: 42}, s.TopOffset(0))
	r.Equal(1, s.Len())

	s.Push(43, 0, 0, 0)
	r.Equal(&stack.Node{Type: 43}, s.TopOffset(0))
	r.Equal(2, s.Len())

	s.Pop()
	r.Equal(&stack.Node{Type: 42}, s.TopOffset(0))
	r.Equal(1, s.Len())

	s.Pop()
	r.Nil(s.TopOffset(0))
	r.Equal(0, s.Len())
}

func TestPopNoop(t *testing.T) {
	r := require.New(t)
	s := stack.New(0)

	r.Zero(s.TopOffset(0))
	r.Zero(s.Len())

	s.Pop() // No-op

	r.Zero(s.Top())
	r.Zero(s.TopOffset(0))
	r.Zero(s.Len())
}

func TestTopOffset(t *testing.T) {
	r := require.New(t)
	s := stack.New(0)

	s.Push(1, 0, 0, 0)
	s.Push(2, 0, 0, 0)
	s.Push(3, 0, 0, 0)
	s.Push(4, 0, 0, 0)
	s.Push(5, 0, 0, 0)

	r.Equal(&stack.Node{Type: 5}, s.Top())
	r.Equal(&stack.Node{Type: 5}, s.TopOffset(0))
	r.Equal(&stack.Node{Type: 4}, s.TopOffset(1))
	r.Equal(&stack.Node{Type: 3}, s.TopOffset(2))
	r.Equal(&stack.Node{Type: 2}, s.TopOffset(3))
	r.Equal(&stack.Node{Type: 1}, s.TopOffset(4))
	r.Nil(s.TopOffset(5))
}

func TestAllocReset(t *testing.T) {
	s := stack.New(16)
	const ln = 128
	for i := 0; i < ln; i++ {
		s.Push(0, i, 0, 0)
	}
	require.Equal(t, ln, s.Len())

	s.Reset()
	require.Nil(t, s.Top())
	require.Zero(t, s.Len())
}

func TestNodeTypeString(t *testing.T) {
	for _, tt := range []struct {
		in  stack.NodeType
		exp string
	}{
		{stack.NodeTypeObject, "object"},
		{stack.NodeTypeArray, "array"},
		{stack.NodeType(0), ""},
	} {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tt.exp, tt.in.String())
		})
	}
}
