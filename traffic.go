package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
)

func main() {
	ctx := context.Background()
	circle := make(chan Car, 8)
	exit := make(chan Car, 8)

	go FirstInputRoad(circle)
	go SecondInputRoad(circle)
	go ThirdInputRoad(circle)
	go FourthInputRoad(circle)
	go FirstOutputRoad(exit)
	go SecondOutputRoad(exit)
	go ThirdOutputRoad(exit)
	go FourthOutputRoad(exit)

	go trafficChecker(circle, exit)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	// Block until a signal is received.
	select {
	case s := <-c:
		fmt.Println("Got signal:", s)
		os.Exit(0)
	case <-ctx.Done():
		fmt.Println("Got ctx.Done:", ctx.Err())
		os.Exit(0)
	}

}

// Counter ...
type Counter struct {
	sync.RWMutex
	n int
}

var counter = Counter{}
var counter2 int32

// Increment a counter
func (c *Counter) Increment() int {
	c.Lock()
	defer c.Unlock()
	c.n++
	return c.n
}

// Car ...
type Car struct {
	Name    string
	entered time.Time
}

// NewCar - creates a car
func NewCar() Car {
	car := Car{Name: "car-" + strconv.Itoa(counter.Increment())}
	//car.Name = "car-" + strconv.Itoa(int(atomic.AddInt32(&counter2, 1)))
	return car
}

func getRandomInt(n int) int {
	randSeed := rand.New(rand.NewSource(time.Now().UnixNano()))
	return randSeed.Intn(n)
}

func trafficChecker(circle chan Car, toExit chan Car) {
	for car := range circle {
		go func(car Car) {
			time.Sleep(1 * time.Second)
			toExit <- car
		}(car)
	}
}

func inputCar(circle chan Car, roadName string, every time.Duration) {
	car := NewCar()
	car.entered = time.Now()
	fmt.Printf("%s comes from %s Input Road\n", car.Name, roadName)
	circle <- car
	time.Sleep(every)
}

// FirstInputRoad generates random number (0 to 5) of cars per second
func FirstInputRoad(cirlce chan Car) {
	n := getRandomInt(5)
	for i := 0; i < n; i++ {
		inputCar(cirlce, "1st", 0)
	}
}

// SecondInputRoad generates 1 car each second
func SecondInputRoad(circle chan Car) {
	for {
		inputCar(circle, "2nd", 1*time.Second)
	}
}

// ThirdInputRoad generates a car each two seconds
func ThirdInputRoad(cirlce chan Car) {
	for {
		inputCar(cirlce, "3rd", 2*time.Second)
	}
}

// FourthInputRoad generates 10 cars per second
func FourthInputRoad(cirlce chan Car) {
	for {
		for i := 0; i < 10; i++ {
			inputCar(cirlce, "4th", 0)
		}
		time.Sleep(1 * time.Second)
	}
}

func outputCar(exit chan Car, roadName string, every time.Duration) {
	car := <-exit
	fmt.Printf("%s is exiting into %s Output Road, time: %s\n",
		car.Name, roadName, time.Now().Sub(car.entered))
	time.Sleep(every)
}

// FirstOutputRoad receives random number (0 to 5) of cars per second
func FirstOutputRoad(exit chan Car) {
	n := getRandomInt(5)
	for {
		for i := 0; i < n; i++ {
			outputCar(exit, "1st", 0)
		}
		time.Sleep(1 * time.Second)
	}
}

// SecondOutputRoad receives one car each second
func SecondOutputRoad(exit chan Car) {
	for {
		outputCar(exit, "2nd", 1*time.Second)
	}
}

// ThirdOutputRoad receives a car each hour
func ThirdOutputRoad(exit chan Car) {
	for {
		outputCar(exit, "3rd", 1*time.Hour)
	}
}

// FourthOutputRoad receives 10 cars per second
func FourthOutputRoad(exit chan Car) {
	for {
		for i := 0; i < 10; i++ {
			outputCar(exit, "4th", 0)
		}
		time.Sleep(1 * time.Second)
	}
}
