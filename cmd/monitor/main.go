package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	numFloors = 10

	Up   ElevatorDirection = 1
	Down ElevatorDirection = -1
	Stop ElevatorDirection = 0
)

type ElevatorDirection int

func (d ElevatorDirection) String() string {
	return [...]string{"Вниз", "Стоп", "Вверх"}[d+1]
}

type Request struct {
	Floor       int
	TargetFloor int
	Direction   ElevatorDirection
}

type Elevator struct {
	CurrentFloor int
	Direction    ElevatorDirection
	targetFloors []bool
	requests     []*Request
}

type Monitor struct {
	*Elevator
	Mutex sync.Mutex
}

func NewMonitor(e *Elevator) *Monitor {
	return &Monitor{
		Elevator: e,
		Mutex:    sync.Mutex{},
	}
}

func (m *Monitor) Move() {
	for {
		time.Sleep(time.Second / 2)
		m.Mutex.Lock()
		m.Elevator.Move()
		m.Mutex.Unlock()
	}
}

func (m *Monitor) Request(req Request) {
	for {
		m.Mutex.Lock()
		res := m.Elevator.Request(req)
		if res {
			m.Mutex.Unlock()
			return
		}
		m.Mutex.Unlock()
	}
}

func NewElevator() *Elevator {
	return &Elevator{
		CurrentFloor: 1,
		Direction:    Stop,
		targetFloors: make([]bool, numFloors),
		requests:     make([]*Request, numFloors),
	}
}

func (e *Elevator) Move() {
	if e.CurrentFloor < 1 || e.CurrentFloor > numFloors {
		if e.CurrentFloor < 1 {
			e.CurrentFloor = 1
			e.Direction = Up
		} else {
			e.CurrentFloor = numFloors
			e.Direction = Down
		}
	}

	for i, needExit := range e.targetFloors {
		direction := e.Direction
		if e.CurrentFloor != i+1 {
			delta := i + 1 - e.CurrentFloor
			direction = ElevatorDirection(delta / abs(delta))
		}
		if i+1 == e.CurrentFloor && needExit && e.Direction == direction {
			fmt.Printf("Высадили на %d этаже\n", e.CurrentFloor)
			e.targetFloors[i] = false
			e.requests[i] = nil
		}
	}

	e.CurrentFloor += int(e.Direction)
}

func (e *Elevator) Request(req Request) bool {
	delta := 0
	direction := e.Direction
	if e.CurrentFloor != req.Floor {
		delta = req.Floor - e.CurrentFloor
		direction = ElevatorDirection(delta / abs(delta))
	}

	if ((e.Direction != direction) || (e.Direction != req.Direction)) && (e.Direction != Stop) {
		return false
	}

	if e.Direction == Stop {
		e.Direction = Up
	}

	fmt.Printf("Надо забрать с %d этажа\n", req.Floor)
	e.targetFloors[req.TargetFloor-1] = true
	return true
}

func main() {
	elevator := NewElevator()
	monitor := NewMonitor(elevator)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		monitor.Move()
	}()

	for i := 0; i < numFloors; i++ {
		dir := Up
		if i%2 != 0 || i == numFloors-1 {
			dir = Down
		}

		var target int
		if dir == Up {
			target = i + (1 + rand.Intn(numFloors-i-1))
		} else {
			target = i - (1 + rand.Intn(i))
		}
		req := Request{Floor: i + 1, TargetFloor: target + 1, Direction: dir}

		wg.Add(1)
		go func(request Request) {
			defer wg.Done()
			fmt.Printf("Вызов с %d этажа %s на %d этаж\n", req.Floor, req.Direction, req.TargetFloor)
			monitor.Request(req)
		}(req)
	}
	wg.Wait()
	time.Sleep(time.Minute)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
