package main

import "testing"

func TestDataProcess(t *testing.T) {
	country := Country{Id: 2, Name: "DPR"}
	employees := []Employee{
		{Id: 6, Name: "Ivan Ivanov", Phone: "0713508765", Address: "Makeevka, University street, 24", NumYearWork: 5, Country: country},
		{Id: 7, Name: "Ivan1 Ivanov", Phone: "0713508765", Address: "Donetsk, University street, 24", NumYearWork: 2, Country: country},
		{Id: 8, Name: "Ivan2 Ivanov", Phone: "0713508765", Address: "Gorlovka, Lenina street, 20", NumYearWork: 1, Country: country},
		{Id: 9, Name: "Ivan3 Ivanov", Phone: "0713508765", Address: "Donetsk, University street, 24", NumYearWork: 6, Country: country},
		{Id: 10, Name: "Ivan4 Ivanov", Phone: "0713508765", Address: "Gorlovka, University street, 88", NumYearWork: 1, Country: country},
	}
	arr := dataProcess(employees, 2, "Donetsk")
	got := arr[0].Name
	want := "Ivan3 Ivanov"
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
