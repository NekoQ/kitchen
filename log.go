package main

import (
	"fmt"
	"os"
)

func LogOrders() {
	file, err := os.OpenFile("logs/orders.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer file.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	for logMsg := range OrdersLogChan {
		fmt.Fprintln(file, logMsg)
	}
}

func LogCooks() {
	file, err := os.OpenFile("logs/cooks.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer file.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	for logMsg := range CooksLogChan {
		fmt.Fprintln(file, logMsg)
	}
}
