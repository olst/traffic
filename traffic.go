package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	circle := make(chan Car, 8)
	exit := make(chan Car, 8)

	//go FirstInputRoad(circle)
	go SecondInputRoad(circle)
	//go ThirdInputRoad(circle)
	//go FourthOutputRoad(circle)
	//go FirstOutputRoad(exit)
	go SecondOutputRoad(exit)
	//go ThirdOutputRoad(exit)
	//go FourthOutputRoad(exit)

	go trafficChecker(circle, exit)

	ex := make(chan string)
	for {
		select {
		case <-ex:
			{
				os.Exit(0)
			}
		}
	}
}

// Counter ...
type Counter struct {
	sync.Mutex
	n int
}

var counter = Counter{n: 0}

// Increment a counter
func (c *Counter) Increment() {
	c.Lock()
	c.n++
	c.Unlock()
}

// Car ...
type Car struct {
	Name    string
	entered time.Time
}

// NewCar - creates a car
func NewCar() Car {
	car := Car{}
	counter.Increment()
	car.Name = "car-" + strconv.Itoa(counter.n)
	return car
}

// FirstInputRoad generates random number (0 to 5) of cars per second
func FirstInputRoad(cirlce chan Car) {
	n := getRandomInt(5)
	for {
		for i := 0; i < n; i++ {
			car := NewCar()
			car.entered = time.Now()
			fmt.Println(car.Name + " comes from 1st Input Road")
			cirlce <- car
		}
		time.Sleep(1 * time.Second)
	}
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
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	return r.Intn(n)
}

func trafficChecker(circle chan Car, toExit chan Car) {
	for {
		for car := range circle {
			delta := time.Now().Sub(car.entered)
			if delta > 3*time.Second {
				fmt.Printf("car: %s, time: %s\n", car.Name, delta)
				toExit <- car
			} else {
				circle <- car
			}
		}
	}
}
