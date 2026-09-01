package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type employee struct {
	ID          	int    `json:"id"`
	EmployeeName 	string `json:"name"`
	Tel          	string `json:"tel"`
	Email        	string `json:"email"`
}

func main() {
	e := employee{}
	err := json.Unmarshal([]byte(`{"id":101,"name":"Fahsai","tel":"0900000000","email":"[EMAIL_ADDRESS]"}`), &e)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(e)
	fmt.Println(e.EmployeeName)
	fmt.Println(e.Tel)
	fmt.Println(e.Email)
}