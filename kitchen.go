package main

import (
	"fmt"
	"time"
)

var (
	OrderChan = make(chan Order, 100)

	Foods    = make([]Food, 10)
	FoodChan = make(chan Food, 1000)

	Cooks     = make([]Cook, 100)
	RankChans [3]chan Cook

	OrdersLogChan = make(chan string, 100)
	CooksLogChan  = make(chan string, 100)
)

func main() {

	// Prepare the data
	UnmarshalFood()
	UnmarshalCooks()
	FillRankChans()

	// Start logging
	go LogOrders()
	go LogCooks()

	// Start the goroutines
	go pickOrders()
	go distributeFood()

	// Start the server
	a := App{}
	a.Initialize()
	a.Run(":80")
}

// When there is a new order in the order channel
// add the foods to the food channel
func pickOrders() {
	for order := range OrderChan {
		// Log that the order is received
		timeNow := time.Now().Format(time.Stamp)
		OrdersLogChan <- fmt.Sprintf("%v Order %v received", timeNow, order.ID)

		// Send the food from the order to the FoodChan
		for _, id := range order.Items {
			FoodChan <- Foods[id-1]
		}
	}
}

// Start picking cooks for food
func distributeFood() {
	for food := range FoodChan {
		go pickCook(food)
	}
}

// Pick a cook
// based on proficiency and avilability
func pickCook(food Food) {
	for {
		for i := 0; i < 3; i++ {
			if food.Complexity <= i+1 {
				select {
				case cook := <-RankChans[i]:
					go cooking(cook, i, food.ID)
					return
				default:
				}
			}
		}
	}
}

// Cook the food and return the cook
// to the cooks channel
func cooking(cook Cook, rank int, id int) {
	timeNow := time.Now().Format(time.Stamp)
	CooksLogChan <- fmt.Sprintf("%v %v starts cooking %v", timeNow, cook.Name, id)

	time.Sleep(time.Second * time.Duration((rank+1)*3))

	timeNow = time.Now().Format(time.Stamp)
	CooksLogChan <- fmt.Sprintf("%v %v finished cooking %v", timeNow, cook.Name, id)

	RankChans[rank] <- cook
}
