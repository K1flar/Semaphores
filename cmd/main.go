package main

import (
	"critical/internal/semaphore"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	numFloors = 5

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
	CurrentFloor   int
	Direction      ElevatorDirection
	floorReq       []bool
	targetFloorReq []bool
	targetFloors   []int
	Mutex          sync.Mutex
	upSemaphore    semaphore.Semaphore
	downSemaphore  semaphore.Semaphore
}

func NewElevator() *Elevator {
	return &Elevator{
		CurrentFloor:   1,
		Direction:      Up,
		floorReq:       make([]bool, numFloors),
		targetFloorReq: make([]bool, numFloors),
		targetFloors:   make([]int, numFloors),
		upSemaphore:    *semaphore.New(numFloors),
		downSemaphore:  *semaphore.New(numFloors),
	}
}

func (e *Elevator) Move() {
	for {
		time.Sleep(time.Second / 2)
		e.Mutex.Lock()
		if e.CurrentFloor < 1 || e.CurrentFloor > numFloors {
			if e.CurrentFloor < 1 {
				e.CurrentFloor = 1
				e.Direction = Up
			} else {
				e.CurrentFloor = numFloors
				e.Direction = Down
			}
			// fmt.Println("Лифт закончил свою работу")
			e.Mutex.Unlock()
			return
		}

		if e.floorReq[e.CurrentFloor-1] {
			fmt.Printf("Забрали пассажира на %d этаже, надо высадить на: %d\n", e.CurrentFloor, e.targetFloors[e.CurrentFloor-1])
			e.targetFloorReq[e.targetFloors[e.CurrentFloor-1]-1] = true
			e.floorReq[e.CurrentFloor-1] = false
		}

		if e.targetFloorReq[e.CurrentFloor-1] {
			fmt.Printf("Высадили пассажира на %d этаже\n", e.CurrentFloor)
			e.targetFloorReq[e.CurrentFloor-1] = false
		}

		// fmt.Println(e.floorReq)

		e.CurrentFloor += int(e.Direction)
		e.Mutex.Unlock()
	}
}

func (e *Elevator) Request(req Request) {

	for {
		e.Mutex.Lock()
		delta := 0
		direction := e.Direction
		if e.CurrentFloor != req.Floor {
			delta = req.Floor - e.CurrentFloor
			direction = ElevatorDirection(delta / abs(delta))
		}

		// if e.Direction != Stop && (e.Direction != req.Direction) {
		if (e.Direction != direction) && (e.Direction != req.Direction) {
			e.Mutex.Unlock()
			continue
		}
		switch req.Direction {
		case Up:
			if e.downSemaphore.Len() == 0 {
				e.upSemaphore.Down()
				e.targetFloors[req.Floor-1] = req.TargetFloor
				e.floorReq[req.Floor-1] = true
				// e.Direction = direction
				// fmt.Println(e.Direction)
				e.Mutex.Unlock()
				fmt.Printf("Вызов с %d этажа %s на %d этаж\n", req.Floor, req.Direction, req.TargetFloor)
				e.Move()
				e.upSemaphore.Up()
				return
			}
		case Down:
			if e.upSemaphore.Len() == 0 {
				e.downSemaphore.Down()
				e.targetFloors[req.Floor-1] = req.TargetFloor
				e.floorReq[req.Floor-1] = true
				// e.Direction = direction
				// fmt.Println(e.Direction)

				e.Mutex.Unlock()
				fmt.Printf("Вызов с %d этажа %s на %d этаж\n", req.Floor, req.Direction, req.TargetFloor)
				e.Move()
				e.downSemaphore.Up()
				return
			}
		}

		e.Mutex.Unlock()
	}

}

func main() {
	elevator := NewElevator()

	wg := sync.WaitGroup{}

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
			elevator.Request(request)
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
