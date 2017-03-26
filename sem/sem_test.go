package sem

import "testing"

func TestFatal(t *testing.T) {
	sem := NewSem(-1)

	if sem != nil {
		t.Error("sem should be nil")
	}
}

func TestSimple(t *testing.T) {
	sem := NewSem(0)
	cond := true

	go func() {
		cond = false
		sem.Up()
	}()

	sem.Down()

	if cond {
		t.Error("main goroutine awoken too early")
	}
}