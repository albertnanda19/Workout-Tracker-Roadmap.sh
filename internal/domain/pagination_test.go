package domain

import "testing"

func TestNewPagination_DefaultsAndBounds(t *testing.T) {
	p := NewPagination(0, 0)
	if p.Page != 1 {
		t.Fatalf("expected page=1, got %d", p.Page)
	}
	if p.Limit != 10 {
		t.Fatalf("expected limit=10, got %d", p.Limit)
	}

	p = NewPagination(-10, -5)
	if p.Page != 1 {
		t.Fatalf("expected page=1, got %d", p.Page)
	}
	if p.Limit != 10 {
		t.Fatalf("expected limit=10, got %d", p.Limit)
	}

	p = NewPagination(2, 999)
	if p.Page != 2 {
		t.Fatalf("expected page=2, got %d", p.Page)
	}
	if p.Limit != 100 {
		t.Fatalf("expected limit=100, got %d", p.Limit)
	}
}
