package backuper

import (
	"container/list"
)

type stringStack struct {
	list   *list.List
	maxLen int
}

func newStringStack(maxLen int) *stringStack {
	return &stringStack{
		list:   list.New(),
		maxLen: maxLen,
	}
}

// Put adds string to stack. Returns popped item if maxLen is exceeded
func (s *stringStack) Put(str string) (string, bool) {
	s.list.PushBack(str)
	if s.list.Len() > s.maxLen {
		popped := s.list.Remove(s.list.Front())
		return popped.(string), true
	}

	return "", false
}

// Head returns last-added string
func (s *stringStack) Head() string {
	last := s.list.Back()
	if last == nil {
		return ""
	}
	return last.Value.(string)
}

// List returns all elements in stack
func (s *stringStack) List() []string {
	strings := make([]string, 0, s.list.Len())
	for e := s.list.Back(); e != nil; e = e.Prev() {
		strings = append(strings, e.Value.(string))
	}
	return strings
}
