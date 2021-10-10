package main

import (
	"fmt"
	"time"
)

var OrderChan = make(chan Order, 100)

var Foods = make([]Food, 10)
var FoodChan = make(chan Food, 1000)

var Cooks = make([]Cook, 100)
var RankChans [3]chan Cook

func main() {

	// Prepare the data
	UnmarshalFood()
	UnmarshalCooks()
	FillRankChans()

	// Start the goroutines
	go pickOrders()
	go pickCooks()

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

func pickCooks() {
	for food := range FoodChan {
	L:
		for {
			for i := 0; i < 3; i++ {
				if food.Complexity <= i+1 {
					select {
					case cook := <-RankChans[i]:
						go cooking(cook, i, food.ID)
						break L
					default:
					}
				}
			}
		}
	}
}

func cooking(cook Cook, rank int, id int) {
	fmt.Printf("Start cooking %v\n", cook.Name)
	time.Sleep(time.Second * time.Duration((rank+1)*3))
	fmt.Printf("%v cooked %v\n", cook.Name, id)
	RankChans[rank] <- cook

}
