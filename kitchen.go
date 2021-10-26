package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	TIME_UNIT = 250

	ReceivedOrdersChan = make(chan Order, 100)
	OrdersChan         = make(map[int]chan Food)

	Foods    = make([]Food, 100)
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
	for order := range ReceivedOrdersChan {
		// Log that the order is received
		timeNow := time.Now().Format(time.Stamp)
		logStr := fmt.Sprintf("%v Order %v received", timeNow, order.ID)
		fmt.Println(logStr)
		OrdersLogChan <- logStr

		// Create a waiting for food channel
		OrdersChan[order.ID] = make(chan Food, 100)
		go waitForOrder(order)

		// Send the food from the order to the FoodChan
		for _, id := range order.Items {
			food := Foods[id-1]
			food.orderID = order.ID
			FoodChan <- food
		}
	}
}

// Start picking cooks for food
func distributeFood() {

	for {
		select {
		case food := <-FoodChan:
			if pickCook(food) == false {
				FoodChan <- food
			}
		default:
		}
	}
}

// Pick a cook
// based on proficiency and availability
func pickCook(food Food) bool {
	for i := 0; i < 3; i++ {
		if food.Complexity <= i+1 {
			select {
			case cook := <-RankChans[i]:
				go cooking(cook, food)
				return true
			default:
			}
		}
	}
	return false
}

// Cook the food and return the cook
// to the cooks channel
func cooking(cook Cook, food Food) {
	timeNow := time.Now().Format(time.Stamp)
	CooksLogChan <- fmt.Sprintf("%v %v starts cooking %v order %v", timeNow, cook.Name, food.ID, food.orderID)
	food.cookID = cook.ID

	time.Sleep(time.Millisecond * time.Duration(food.PreparationTime*TIME_UNIT))
	OrdersChan[food.orderID] <- food
	RankChans[cook.Rank-1] <- cook

	timeNow = time.Now().Format(time.Stamp)
	CooksLogChan <- fmt.Sprintf("%v %v finished cooking %v order %v", timeNow, cook.Name, food.ID, food.orderID)
}

func waitForOrder(order Order) {
	for food := range OrdersChan[order.ID] {
		var detail CookingDetail
		detail.CookID = food.cookID
		detail.FoodID = food.ID
		order.CookingDetails = append(order.CookingDetails, detail)
		if len(order.CookingDetails) == len(order.Items) {
			order.CookingTime = (int(time.Now().UnixMilli()) - int(order.PickUpTime)) / TIME_UNIT
			sendOrder(order)
		}
	}
}

func sendOrder(order Order) {
	url := "http://172.17.0.3:81/distribution"
	jsonValue, _ := json.Marshal(order)
	http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	timeNow := time.Now().Format(time.Stamp)
	logStr := fmt.Sprintf("%v Order %v finished", timeNow, order.ID)
	fmt.Println(logStr)
	OrdersLogChan <- logStr
}
