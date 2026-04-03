package bridge

import (
	"strings"
	"testing"
)

func TestBasicRemote(t *testing.T) {
	tv := NewTV("Living Room TV")
	remote := NewBasicRemote(tv)

	if remote.DeviceName() != "Living Room TV" {
		t.Fatalf("expected device name, got %s", remote.DeviceName())
	}

	remote.TogglePower()
	if !tv.IsEnabled() {
		t.Fatal("expected tv to be enabled")
	}

	remote.TogglePower()
	if tv.IsEnabled() {
		t.Fatal("expected tv to be disabled")
	}
}

func TestSmartRemoteWithTV(t *testing.T) {
	tv := NewTV("Bedroom TV")
	remote := NewSmartRemote(tv)

	t.Run("PowerAndControls", func(t *testing.T) {
		remote.TogglePower()
		remote.VolumeUp(15)
		remote.ChannelNext()

		if !tv.IsEnabled() {
			t.Fatal("expected tv to be enabled")
		}
		if tv.Volume() != 35 {
			t.Fatalf("expected volume 35, got %d", tv.Volume())
		}
		if tv.Channel() != 2 {
			t.Fatalf("expected channel 2, got %d", tv.Channel())
		}
	})

	t.Run("Status", func(t *testing.T) {
		status := remote.Status()
		expectedParts := []string{"Bedroom TV", "power=true", "volume=35", "channel=2"}
		for _, part := range expectedParts {
			if !strings.Contains(status, part) {
				t.Fatalf("expected status to contain %q, got %s", part, status)
			}
		}
	})
}

func TestSmartRemoteWithStreamingBox(t *testing.T) {
	box := NewStreamingBox("Apple TV")
	remote := NewSmartRemote(box)

	remote.TogglePower()
	remote.VolumeUp(100)
	remote.ChannelNext()

	if !box.IsEnabled() {
		t.Fatal("expected streaming box to be enabled")
	}
	if box.Volume() != 50 {
		t.Fatalf("expected capped volume 50, got %d", box.Volume())
	}
	if box.Channel() != 101 {
		t.Fatalf("expected channel 101, got %d", box.Channel())
	}
}

func TestBridgeInterfaceImplementation(t *testing.T) {
	var _ Device = (*TV)(nil)
	var _ EntertainmentDevice = (*TV)(nil)
	var _ EntertainmentDevice = (*StreamingBox)(nil)
}
