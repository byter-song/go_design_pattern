package flyweight

import "testing"

func TestStyleFactorySharing(t *testing.T) {
	factory := NewStyleFactory()

	laser1 := factory.GetStyle("laser", "green", ">>>")
	laser2 := factory.GetStyle("laser", "green", ">>>")
	plasma := factory.GetStyle("plasma", "blue", "***")

	if laser1 != laser2 {
		t.Fatal("expected identical flyweight instance for same style")
	}
	if laser1 == plasma {
		t.Fatal("expected different flyweight instance for different style")
	}
	if factory.Count() != 2 {
		t.Fatalf("expected 2 shared styles, got %d", factory.Count())
	}
}

func TestBulletPoolAcquireRelease(t *testing.T) {
	factory := NewStyleFactory()
	pool := NewBulletPool()
	style := factory.GetStyle("laser", "green", ">>>")

	bullet := pool.Acquire(style, "player-1", 10, 5, 20)
	if bullet.Style != style {
		t.Fatal("expected shared style to be assigned")
	}
	if bullet.Owner != "player-1" || bullet.X != 10 || bullet.Y != 5 || bullet.Speed != 20 || !bullet.Active {
		t.Fatalf("unexpected bullet state: %#v", bullet)
	}

	bullet.Move(0.5)
	if bullet.X != 20 {
		t.Fatalf("expected x position 20, got %.1f", bullet.X)
	}

	pool.Release(bullet)
	if bullet.Style != nil || bullet.Owner != "" || bullet.X != 0 || bullet.Y != 0 || bullet.Speed != 0 || bullet.Active {
		t.Fatalf("expected released bullet to be reset, got %#v", bullet)
	}

	reused := pool.Acquire(style, "player-2", 1, 2, 3)
	if reused.Style != style || reused.Owner != "player-2" || reused.X != 1 || reused.Y != 2 || reused.Speed != 3 || !reused.Active {
		t.Fatalf("unexpected reused bullet state: %#v", reused)
	}
}

func TestBulletPoolCreatedCount(t *testing.T) {
	factory := NewStyleFactory()
	pool := NewBulletPool()
	style := factory.GetStyle("plasma", "blue", "***")

	first := pool.Acquire(style, "a", 0, 0, 1)
	second := pool.Acquire(style, "b", 0, 0, 1)

	if pool.CreatedCount() != 2 {
		t.Fatalf("expected 2 created bullets, got %d", pool.CreatedCount())
	}

	pool.Release(first)
	pool.Release(second)

	third := pool.Acquire(style, "c", 0, 0, 1)
	if third == nil {
		t.Fatal("expected acquired bullet")
	}
	if pool.CreatedCount() < 2 {
		t.Fatalf("expected created count to remain at least 2, got %d", pool.CreatedCount())
	}
}

func BenchmarkBulletPool(b *testing.B) {
	factory := NewStyleFactory()
	pool := NewBulletPool()
	style := factory.GetStyle("laser", "green", ">>>")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bullet := pool.Acquire(style, "bench", 0, 0, 10)
		bullet.Move(1)
		pool.Release(bullet)
	}
}
