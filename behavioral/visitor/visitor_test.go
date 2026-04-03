package visitor

import (
	"strings"
	"testing"
)

func TestVisitors(t *testing.T) {
	root := &Directory{
		Name: "root",
		Children: []Element{
			&File{Name: "a.txt", Size: 10},
			&Directory{
				Name: "docs",
				Children: []Element{
					&File{Name: "b.md", Size: 20},
				},
			},
		},
	}

	sizeV := &SizeVisitor{}
	nameV := &NameVisitor{}

	root.Accept(sizeV)
	root.Accept(nameV)

	if sizeV.Total != 30 {
		t.Fatalf("expected total size 30, got %d", sizeV.Total)
	}
	summary := strings.Join(nameV.Names, ",")
	expected := "dir:root,file:a.txt,dir:docs,file:b.md"
	if summary != expected {
		t.Fatalf("expected %s, got %s", expected, summary)
	}
}

func TestVisitorInterfaces(t *testing.T) {
	var _ Element = (*File)(nil)
	var _ Element = (*Directory)(nil)
}
