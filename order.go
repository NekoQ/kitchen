package main

type Order struct {
	ID             int             `json:"order_id"`
	TableID        int             `json:"table_id"`
	WaiterID       int             `json:"waiter_id"`
	Items          []int           `json:"items"`
	Priority       int             `json:"priority"`
	MaxWait        int             `json:"max_wait"`
	PickUpTime     int64           `json:"pick_up_time"`
	CookingTime    int             `json:"cooking_time"`
	CookingDetails []CookingDetail `json:"cooking_details"`
}

type CookingDetail struct {
	FoodID int `json:"food_id"`
	CookID int `json:"cook_id"`
}
