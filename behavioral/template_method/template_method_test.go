package template_method

import "testing"

func TestReportGenerator(t *testing.T) {
	generator := NewReportGenerator(SalesReport{})
	got := generator.Generate()
	want := "[SALES] SALES-TOTAL=600 FROM SALES:100,200,300"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestPipeline(t *testing.T) {
	pipeline := DefaultCSVImport()
	got := pipeline.Execute()
	want := "saved:id|name\n1|go"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestTemplateInterfaceImplementation(t *testing.T) {
	var _ ReportSteps = SalesReport{}
}
