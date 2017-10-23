package store

import "testing"

func TestSet(t *testing.T) {
	s := NewMemoryStore(100)
	if err := s.Set("foo", "value"); err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	if _, ok := s.s.Load("foo"); !ok {
		t.Errorf("key not stored")
	}
}

func TestGet(t *testing.T) {
	s := NewMemoryStore(100)
	s.s.Store("foo", "bar")
	v, ok := s.Get("foo")
	if !ok {
		t.Error("unexpected false return")
	}

	if v != "bar" {
		t.Errorf("got: %v, wants: %v", v, "bar")
	}

	_, ok = s.Get("baz")
	if ok {
		t.Error("unexpected key found")
	}
}

func TestDelete(t *testing.T) {
	s := NewMemoryStore(100)
	s.s.Store("foo", "bar")
	s.Delete("foo")
	if _, ok := s.s.Load("foo"); ok {
		t.Errorf("unexpected key found")
	}
}
