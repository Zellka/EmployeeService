package repository

import (
	model "MainGoTask/model"
	"encoding/json"
	"log"
	"net/http"
)

type WebRepository struct {
	url string
}

func NewWebRepository(url string) *WebRepository {
	return &WebRepository{
		url: url,
	}
}

func (r WebRepository) GetEmployeesFromWeb() []model.Employee {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, r.url, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	empList := []Employee{}
	errDecode := json.NewDecoder(res.Body).Decode(&empList)
	if errDecode != nil {
		log.Fatal(errDecode)
	}
	return toEmployees(empList)
}

type Employee struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	NumYearWork int64  `json:"numYearWork"`
}

func toEmployees(employees []Employee) []model.Employee {
	out := make([]model.Employee, len(employees))
	for i, b := range employees {
		emp := model.Employee(b)
		out[i] = emp
	}
	return out
}
