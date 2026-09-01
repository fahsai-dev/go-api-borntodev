package main

import (
	"encoding/json"
	"fmt"
)

type employee struct {
	ID          	int    `json:"id"`
	EmployeeName 	string `json:"name"`
	Tel          	string `json:"tel"`
	Email        	string `json:"email"`
}

func main() {
	data, _ := json.Marshal(&employee{101, "Fahsai", "0900000000", "[EMAIL_ADDRESS]"})

	fmt.Println(string(data))
}