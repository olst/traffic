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

	//ex := make(chan struct{})
	//for {
	//	select {
	//	case <-ex:
	//		os.Exit(0)
	//	}
	//}

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

// Increment a counter
func (c *Counter) Increment() int {
	c.Lock()
	defer c.Unlock()
	c.n++
	return c.n
}
func (c *Counter) GetN() int {
	c.RLock()
	defer c.RUnlock()
	return c.n
}

var counter2 int32

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

func executeEvery(fn func(), every time.Duration) {
	for {
		fn()
		time.Sleep(every * time.Second)
	}
}

// FirstInputRoad generates random number (0 to 5) of cars per second
func FirstInputRoad(cirlce chan Car) {
	n := getRandomInt(5)
	executeEvery(func() {
		for i := 0; i < n; i++ {
			car := NewCar()
			car.entered = time.Now()
			fmt.Println(car.Name + " comes from 1st Input Road")
			cirlce <- car
		}
	}, 1)
}

// SecondInputRoad generates 1 car each second
func SecondInputRoad(cirlce chan Car) {
	for {
		car := NewCar()
		car.entered = time.Now()
		fmt.Println(car.Name + " comes from 2nd Input Road")
		cirlce <- car
		time.Sleep(1 * time.Second)
	}
}

// ThirdInputRoad generates a car each two seconds
func ThirdInputRoad(cirlce chan Car) {
	for {
		car := NewCar()
		car.entered = time.Now()
		fmt.Println(car.Name + " comes from 3rd Input Road")
		cirlce <- car
		time.Sleep(2 * time.Second)
	}
}

// FourthInputRoad generates 10 cars per second
func FourthInputRoad(cirlce chan Car) {
	for {
		for i := 0; i < 10; i++ {
			car := NewCar()
			car.entered = time.Now()
			fmt.Println(car.Name + " comes from 4th Input Road")
			cirlce <- car
		}
		time.Sleep(1 * time.Second)
	}
}

// FirstOutputRoad receives random number (0 to 5) of cars per second
func FirstOutputRoad(exit chan Car) {
	n := getRandomInt(5)
	for {
		for i := 0; i < n; i++ {
			car := <-exit
			fmt.Println(car.Name + " leaves into 1st Output Road")
		}
		time.Sleep(1 * time.Second)
	}
}

// SecondOutputRoad receives one car each second
func SecondOutputRoad(exit chan Car) {
	for {
		car := <-exit
		fmt.Println(car.Name + " leaves into 2nd Output Road")
		time.Sleep(1 * time.Second)
	}
}

// ThirdOutputRoad receives a car each hour
func ThirdOutputRoad(exit chan Car) {
	for {
		car := <-exit
		fmt.Println(car.Name + " leaves into 3rd Output Road")
		time.Sleep(1 * time.Hour)
	}
}

// FourthOutputRoad receives 10 cars per second
func FourthOutputRoad(exit chan Car) {
	for {
		for i := 0; i < 10; i++ {
			car := <-exit
			fmt.Println(car.Name + " leaves into 1st Output Road")
		}
		time.Sleep(1 * time.Second)
	}
}

func getRandomInt(n int) int {
	randSeed := rand.New(rand.NewSource(time.Now().UnixNano()))
	return randSeed.Intn(n)
}

func trafficChecker(circle chan Car, toExit chan Car) {
	for car := range circle {
		delta := time.Now().Sub(car.entered)
		// TODO: Time depends
		if delta < 3*time.Second {
			// AAAAAAAAAAAAAA!!!!!!!!!!!!!!!!!!!!!!!!!
			circle <- car
			time.Sleep(time.Millisecond * 100)
			continue
		}
		fmt.Fprintf(os.Stdout, "car: %s, time: %s\n. %d Cars left", car.Name, delta, len(circle))
		toExit <- car
	}
}

//input1-        - Exit1 -output1 4s
//input2- CIRCLE - Exit2 -output2 1s
//input3-        - Exit3 -output3 2s
//input4-        - Exit4 -output4 3s
