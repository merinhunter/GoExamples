package rendez

import (
	"fmt"
	"testing"
	"time"
)

const (
	Sleept = 10 * time.Millisecond
	Nthreads = 500
)

func f(n int, tag int, s int) {
	fmt.Printf("%d: Entra %d con tag %d\n", n, s, tag)
	str := Rendezvous(tag, s)
	fmt.Printf("%d: He cambiado %d por %d\n", n, s, str)
}

func TestConcurrency(t *testing.T) {
	for i := 1; i < Nthreads; i++ {
		go f(i, 1, i)
		time.Sleep(Sleept)
	}

	f(0, 1, 0)
}

func TestMultipleThreads(t *testing.T) {
	for i := 0; i < 15; i++ {
		go f(i, i, i)
	}

	time.Sleep(Sleept)

	for i := 14; i >= 0; i-- {
		go f(i, i, i)
	}

	time.Sleep(200 * Sleept)
}