package main

import (
	rep "MainGoTask/employee/repository"
	usecase "MainGoTask/employee/usecase"
	"MainGoTask/models"
	"testing"
)

func TestDataProcess(t *testing.T) {
	usecase := usecase.NewEmployeeUseCase(rep.EmployeeRepository{}, rep.WebRepository{})

	employees := []models.Employee{
		{Id: 6, Name: "Ivan Ivanov", Phone: "0713508765", Address: "Makeevka, University street, 24", NumYearWork: 5},
		{Id: 7, Name: "Ivan1 Ivanov", Phone: "0713508765", Address: "Donetsk, University street, 24", NumYearWork: 2},
		{Id: 8, Name: "Ivan2 Ivanov", Phone: "0713508765", Address: "Gorlovka, Lenina street, 20", NumYearWork: 1},
		{Id: 9, Name: "Ivan3 Ivanov", Phone: "0713508765", Address: "Donetsk, University street, 24", NumYearWork: 6},
		{Id: 10, Name: "Ivan4 Ivanov", Phone: "0713508765", Address: "Gorlovka, University street, 88", NumYearWork: 1},
	}
	arr := usecase.DataProcess(employees, 2, "Donetsk")
	got := arr[0].Name
	want := "Ivan3 Ivanov"
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
