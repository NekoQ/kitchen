package main

import (
	"fmt"
	"time"
)

var OrderChan = make(chan Order, 100)

var Foods = make([]Food, 10)
var FoodChan = make(chan Food, 1000)
var TmpFoodChan = make(chan Food, 1000)
var tmpNotEmpty bool

var Cooks = make([]Cook, 100)
var RankChans [3]chan Cook

func main() {

	// Prepare the data
	UnmarshalFood()
	UnmarshalCooks()
	FillRankChans()

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
		for _, id := range order.Items {
			FoodChan <- Foods[id-1]
		}
	}
}

func distributeFood() {
	for food := range FoodChan {
		go pickCook(food)
	}
}

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

func cooking(cook Cook, rank int, id int) {
	fmt.Printf("Start cooking %v %v\n", cook.Name, id)
	time.Sleep(time.Second * time.Duration((rank+1)*3))
	fmt.Printf("%v cooked %v\n", cook.Name, id)
	RankChans[rank] <- cook
}
