package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	NumFloors = 5
	NumCalls  = 10
)

type Elevator struct {
	currentFloor int
	direction    int // 1 for up, -1 for down
}

func NewElevator() *Elevator {
	return &Elevator{
		currentFloor: 1,
		direction:    1,
	}
}

func (e *Elevator) Move() {
	for {
		// Check if there are any requests in the current direction
		if e.direction == 1 {
			if e.currentFloor == NumFloors {
				e.direction = -1
			}
		} else {
			if e.currentFloor == 1 {
				e.direction = 1
			}
		}

		fmt.Printf("Elevator is on floor %d\n", e.currentFloor)

		// Sleep for simulation of moving to next floor
		time.Sleep(1 * time.Second)

		// Update current floor based on direction
		e.currentFloor += e.direction
	}
}

func main() {
	elevator := NewElevator()
	var mutex sync.Mutex

	go elevator.Move()

	// Simulating request buttons being pressed
	for i := 0; i < NumCalls; i++ {
		// Start a new goroutine for each call
		go func() {
			// Simulate a request being made from different floors
			floor := 1 + time.Now().Second()%NumFloors

			// Simulate critical section entry
			for {
				mutex.Lock()
				// fmt.Printf("Trying to enter critical section for floor %d\n", floor)

				// Process request
				if (floor == elevator.currentFloor) || (floor == NumFloors && elevator.direction == 1) || (floor == 1 && elevator.direction == -1) {
					fmt.Printf("Entering critical section for floor %d\n", floor)
					// Process request
					time.Sleep(2 * time.Second)
					fmt.Printf("Exiting critical section for floor %d\n", floor)
					mutex.Unlock()
					return
				}

				// Simulate critical section exit
				mutex.Unlock()
			}
		}()
	}

	// Keep the main goroutine running
	select {}
}
