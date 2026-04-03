package iterator

import (
	"slices"
	"testing"
)

func TestClosureIterator(t *testing.T) {
	collection := NewIntCollection(1, 2, 3)
	next := collection.Iterator()

	values := make([]int, 0)
	for {
		v, ok := next()
		if !ok {
			break
		}
		values = append(values, v)
	}

	if !slices.Equal(values, []int{1, 2, 3}) {
		t.Fatalf("unexpected values: %v", values)
	}
}

func TestSeqIterator(t *testing.T) {
	collection := NewIntCollection(2, 4, 6)
	values := make([]int, 0)

	for item := range collection.Seq() {
		values = append(values, item)
	}

	if !slices.Equal(values, []int{2, 4, 6}) {
		t.Fatalf("unexpected values: %v", values)
	}
}

func TestFilterSeq(t *testing.T) {
	collection := NewIntCollection(1, 2, 3, 4, 5)
	values := make([]int, 0)

	for item := range FilterSeq(collection.Seq(), func(v int) bool {
		return v%2 == 0
	}) {
		values = append(values, item)
	}

	if !slices.Equal(values, []int{2, 4}) {
		t.Fatalf("unexpected values: %v", values)
	}
}
