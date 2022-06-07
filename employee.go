package main

type Employee struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Phone       string  `json:"phone"`
	Address     string  `json:"address"`
	NumYearWork int     `json:"numYearWork"`
	Country     Country `json:"countryId"`
}
