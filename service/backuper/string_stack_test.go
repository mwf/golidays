package backuper

import (
	"reflect"
	"testing"
)

func TestStringStackHead_empty(t *testing.T) {
	s := newStringStack(1)

	if s.Head() != "" {
		t.Fail()
	}
}

func TestStringStackPutHead(t *testing.T) {
	s := newStringStack(10)

	s.Put("foo")
	s.Put("bar")

	if s.Head() != "bar" {
		t.Fail()
	}
}

func TestStringStackPut_overflow(t *testing.T) {
	s := newStringStack(2)

	_, purged := s.Put("foo")
	if purged {
		t.Fail()
	}
	_, purged = s.Put("bar")
	if purged {
		t.Fail()
	}

	str, purged := s.Put("buz")
	if !purged {
		t.Fail()
	}
	if str != "foo" {
		t.Fail()
	}
}

func TestStringStackList_overflow(t *testing.T) {
	strings := []string{"foo", "bar", "buz"}
	s := newStringStack(2)

	for _, str := range strings {
		s.Put(str)
	}

	expected := []string{"buz", "bar"}
	if !reflect.DeepEqual(s.List(), expected) {
		t.Fatalf("stack list %s != expected %s", s.List(), expected)
	}
}

func TestStringStackList(t *testing.T) {
	strings := []string{"foo", "bar", "buz"}
	s := newStringStack(10)

	for _, str := range strings {
		s.Put(str)
	}

	reversed := []string{"buz", "bar", "foo"}
	if !reflect.DeepEqual(s.List(), reversed) {
		t.Fatalf("stack list %s != expected %s", s.List(), reversed)
	}
}

func TestStringStackList_empty(t *testing.T) {
	s := newStringStack(10)

	if len(s.List()) > 0 {
		t.Fatalf("list %s should be empty", s.List())
	}
}

func TestStringStackList_single(t *testing.T) {
	s := newStringStack(1)
	s.Put("foo")

	expected := []string{"foo"}
	if !reflect.DeepEqual(s.List(), expected) {
		t.Fatalf("stack list %s != expected %s", s.List(), expected)
	}
}
