package main

type Response struct {
	Success bool
	Result  float64
	Query   struct {
		From   string
		To     string
		Amount float32
	}
	Date string
}
