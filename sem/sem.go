package sem

import "sync"

type UpDowner interface {
	Up()
	Down()
}

type Sem struct {
	ntok int
	cond *sync.Cond
}

var mutex sync.Mutex

func condition(s *Sem) bool {
	if s.ntok != 0 {
		return true
	}

	return false
}

func NewSem(ntok int) *Sem {
	if ntok < 0 {
		return nil
	}

	sem := new(Sem)
	sem.ntok = ntok
	sem.cond = sync.NewCond(&mutex)

	return sem
}

func (s *Sem) Up() {
	s.cond.L.Lock()

	s.ntok++
	s.cond.Signal()

	s.cond.L.Unlock()
}

func (s *Sem) Down() {
	s.cond.L.Lock()
	for !condition(s) {
		s.cond.Wait()
	}

	s.ntok--
	s.cond.L.Unlock()
}