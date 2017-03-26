package rendez

import "sync"

type Value struct {
	value interface{}
	wg sync.WaitGroup
}

var (
	mutex sync.Mutex
	goroutines = make(map[int]*Value)
)

func Rendezvous(tag int, val interface{}) interface{} {
	mutex.Lock()

	s, exist := goroutines[tag]

	if exist {
		delete(goroutines, tag)
		mutex.Unlock()
		val, s.value = s.value, val
		s.wg.Done()
		return val
	}

	s = new(Value)
	s.value = val
	s.wg.Add(1)

	goroutines[tag] = s

	mutex.Unlock()

	s.wg.Wait()

	return s.value
}