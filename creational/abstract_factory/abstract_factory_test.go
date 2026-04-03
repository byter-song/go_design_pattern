package abstract_factory

import "testing"

func TestAbstractFactory(t *testing.T) {
	tests := []struct {
		name            string
		factory         UIFactory
		expectedTheme   string
		expectedButton  string
		expectedChecked string
	}{
		{
			name:            "MacFactory",
			factory:         &MacFactory{},
			expectedTheme:   "Mac",
			expectedButton:  "MacButton clicked",
			expectedChecked: "MacCheckbox ON",
		},
		{
			name:            "WinFactory",
			factory:         &WinFactory{},
			expectedTheme:   "Windows",
			expectedButton:  "WinButton clicked",
			expectedChecked: "WinCheckbox ON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.factory.Theme() != tt.expectedTheme {
				t.Fatalf("expected theme %s, got %s", tt.expectedTheme, tt.factory.Theme())
			}

			button := tt.factory.CreateButton()
			checkbox := tt.factory.CreateCheckbox()

			if button.Click() != tt.expectedButton {
				t.Fatalf("expected button click result %s, got %s", tt.expectedButton, button.Click())
			}
			if checkbox.Check(true) != tt.expectedChecked {
				t.Fatalf("expected checkbox state %s, got %s", tt.expectedChecked, checkbox.Check(true))
			}
		})
	}
}

func TestBuildUI(t *testing.T) {
	output := BuildUI(&MacFactory{})
	if output != "[Mac] MacButton | MacCheckbox" {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestInterfaceImplementation(t *testing.T) {
	var _ UIFactory = (*MacFactory)(nil)
	var _ UIFactory = (*WinFactory)(nil)
}
