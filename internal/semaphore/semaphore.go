package semaphore

type Semaphore struct {
	s chan struct{}
}

func New(n int) *Semaphore {
	return &Semaphore{
		s: make(chan struct{}, n),
	}
}

func (s *Semaphore) Down() {
	s.s <- struct{}{}
}

func (s *Semaphore) Up() {
	<-s.s
}

func (s *Semaphore) Len() int {
	return len(s.s)
}
