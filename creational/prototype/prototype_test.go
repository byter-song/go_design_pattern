package prototype

import "testing"

func sampleDocument() *Document {
	return &Document{
		Title: "Go Design Patterns",
		Tags:  []string{"go", "design-pattern"},
		Metadata: map[string]string{
			"author": "song",
		},
		Sections: []Section{
			{
				Heading: "Introduction",
				Notes:   []string{"why", "how"},
			},
		},
	}
}

func TestShallowClone(t *testing.T) {
	original := sampleDocument()
	cloned := original.ShallowClone()

	if cloned == original {
		t.Fatal("expected different struct pointer")
	}

	cloned.Tags[0] = "golang"
	cloned.Metadata["author"] = "alice"
	cloned.Sections[0].Notes[0] = "updated"

	if original.Tags[0] != "golang" {
		t.Fatal("expected tags slice to be shared in shallow clone")
	}
	if original.Metadata["author"] != "alice" {
		t.Fatal("expected metadata map to be shared in shallow clone")
	}
	if original.Sections[0].Notes[0] != "updated" {
		t.Fatal("expected nested slice to be shared in shallow clone")
	}
}

func TestDeepClone(t *testing.T) {
	original := sampleDocument()
	cloned := original.DeepClone()

	if cloned == original {
		t.Fatal("expected different struct pointer")
	}

	cloned.Tags[0] = "golang"
	cloned.Metadata["author"] = "alice"
	cloned.Sections[0].Notes[0] = "updated"
	cloned.Sections[0].Heading = "Intro 2"

	if original.Tags[0] != "go" {
		t.Fatal("expected tags slice to be independent in deep clone")
	}
	if original.Metadata["author"] != "song" {
		t.Fatal("expected metadata map to be independent in deep clone")
	}
	if original.Sections[0].Notes[0] != "why" {
		t.Fatal("expected nested slice to be independent in deep clone")
	}
	if original.Sections[0].Heading != "Introduction" {
		t.Fatal("expected section struct to be independent in deep clone")
	}
}

func TestNilClone(t *testing.T) {
	var doc *Document
	if doc.ShallowClone() != nil {
		t.Fatal("expected nil shallow clone")
	}
	if doc.DeepClone() != nil {
		t.Fatal("expected nil deep clone")
	}
}
