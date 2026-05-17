package scheduler

import (
	"sync"
	"time"
)

type Scheduler struct {
	mu       sync.Mutex
	interval time.Duration
	stop     chan struct{}
	reset    chan time.Duration
	fire     func()
	running  bool
}

func New(interval time.Duration, fire func()) *Scheduler {
	return &Scheduler{
		interval: interval,
		fire:     fire,
		stop:     make(chan struct{}),
		reset:    make(chan time.Duration, 1),
	}
}

func (s *Scheduler) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	go s.loop()
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	close(s.stop)
	s.stop = make(chan struct{})
	s.mu.Unlock()
}

func (s *Scheduler) SetInterval(d time.Duration) {
	if d <= 0 {
		return
	}
	s.mu.Lock()
	s.interval = d
	s.mu.Unlock()
	select {
	case s.reset <- d:
	default:
	}
}

func (s *Scheduler) Interval() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.interval
}

func (s *Scheduler) FireNow() {
	if s.fire != nil {
		go s.fire()
	}
}

func (s *Scheduler) loop() {
	s.mu.Lock()
	t := time.NewTimer(s.interval)
	s.mu.Unlock()
	defer t.Stop()

	for {
		select {
		case <-s.stop:
			return
		case d := <-s.reset:
			if !t.Stop() {
				select {
				case <-t.C:
				default:
				}
			}
			t.Reset(d)
		case <-t.C:
			if s.fire != nil {
				s.fire()
			}
			s.mu.Lock()
			next := s.interval
			s.mu.Unlock()
			t.Reset(next)
		}
	}
}
