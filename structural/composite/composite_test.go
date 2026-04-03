package composite

import (
	"strings"
	"testing"
)

func TestDirectoryComposite(t *testing.T) {
	root := NewDirectory("root")
	docs := NewDirectory("docs")
	images := NewDirectory("images")
	readme := NewFile("README.md", 12)
	spec := NewFile("spec.md", 20)
	logo := NewFile("logo.png", 256)

	docs.Add(readme, spec)
	images.Add(logo)
	root.Add(docs, images, NewFile("go.mod", 1))

	t.Run("AggregateSize", func(t *testing.T) {
		if root.Size() != 289 {
			t.Fatalf("expected total size 289, got %d", root.Size())
		}
	})

	t.Run("Count", func(t *testing.T) {
		if root.Count() != 7 {
			t.Fatalf("expected 7 nodes, got %d", root.Count())
		}
	})

	t.Run("Find", func(t *testing.T) {
		node := root.Find("logo.png")
		if node == nil {
			t.Fatal("expected to find logo.png")
		}
		if node.Size() != 256 {
			t.Fatalf("expected logo.png size 256, got %d", node.Size())
		}
	})

	t.Run("Render", func(t *testing.T) {
		output := root.Render(0)
		expectedFragments := []string{
			"+ root/",
			"  + docs/",
			"    - README.md (12KB)",
			"    - spec.md (20KB)",
			"  + images/",
			"    - logo.png (256KB)",
			"  - go.mod (1KB)",
		}
		for _, fragment := range expectedFragments {
			if !strings.Contains(output, fragment) {
				t.Fatalf("expected render output to contain %q, got:\n%s", fragment, output)
			}
		}
	})
}

func TestDirectoryMutation(t *testing.T) {
	root := NewDirectory("root")
	config := NewFile("config.yaml", 3)
	root.Add(config)

	t.Run("RemoveExisting", func(t *testing.T) {
		if !root.Remove("config.yaml") {
			t.Fatal("expected remove to succeed")
		}
		if root.Find("config.yaml") != nil {
			t.Fatal("expected config.yaml to be removed")
		}
	})

	t.Run("RemoveNonExisting", func(t *testing.T) {
		if root.Remove("missing") {
			t.Fatal("expected remove to fail")
		}
	})
}

func TestDirectoryWalk(t *testing.T) {
	root := NewDirectory("root")
	src := NewDirectory("src")
	root.Add(src, NewFile("README.md", 1))
	src.Add(NewFile("main.go", 10), NewFile("util.go", 8))

	visited := make([]string, 0)
	root.Walk(func(node Node) {
		visited = append(visited, node.Name())
	})

	expected := []string{"root", "src", "main.go", "util.go", "README.md"}
	if strings.Join(visited, ",") != strings.Join(expected, ",") {
		t.Fatalf("expected %v, got %v", expected, visited)
	}
}

func TestCompositeInterfaceImplementation(t *testing.T) {
	var _ Node = (*File)(nil)
	var _ Node = (*Directory)(nil)
}
