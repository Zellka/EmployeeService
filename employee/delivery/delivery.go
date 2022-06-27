package delivery

import (
	"MainGoTask/employee"
	"MainGoTask/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Broker struct {
	useCase employee.UseCase
	router  *mux.Router
}

func NewBroker(useCase employee.UseCase) *Broker {
	return &Broker{
		useCase: useCase,
		router:  mux.NewRouter(),
	}
}

func (b *Broker) Start() {
	employees := ParseEmployee(b.SetEmployees())
	fmt.Println("Ожидание...")

	b.router.HandleFunc("/employees", b.HandleGetEmployee(employees))
	http.Handle("/", b.router)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(http.StatusInternalServerError)
	}
}

func (b *Broker) HandleGetEmployee(employees []models.Employee) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		city := req.URL.Query().Get("city")
		numYearWork, err := strconv.ParseInt(req.URL.Query().Get("num"), 10, 64)
		if err != nil {
			log.Fatal(http.StatusInternalServerError)
		}
		if status := b.SaveEmployees(employees, numYearWork, city); status != http.StatusOK {
			log.Fatal(status)
		}
		employees, status := b.GetEmployees()
		if status != http.StatusOK {
			log.Fatal(status)
		}
		printEmployees(w, employees)
	}
}

func printEmployees(w http.ResponseWriter, employees []models.Employee) {
	fmt.Fprintf(w, "<h1>Список сотрудников</h>")
	for _, emp := range employees {
		fmt.Fprintf(w, "<h2>"+emp.Name+"</h>")
	}
}

func (b *Broker) SetEmployees() []byte {
	return b.useCase.SetEmployees()
}

func ParseEmployee(body []byte) []models.Employee {
	empList := []Employee{}
	jsonErr := json.Unmarshal(body, &empList)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return toEmployees(empList)
}

func (b *Broker) SaveEmployees(employees []models.Employee, num int64, city string) int {
	emp := b.useCase.DataProcess(employees, num, city)
	if err := b.useCase.SaveEmployees(emp, num, city); err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func (b *Broker) GetEmployees() ([]models.Employee, int) {
	employees, err := b.useCase.GetEmployees()
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return employees, http.StatusOK
}

type Employee struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	NumYearWork int64  `json:"numYearWork"`
}

func toEmployee(emp Employee) models.Employee {
	return models.Employee{
		Id:          emp.Id,
		Name:        emp.Name,
		Phone:       emp.Phone,
		Address:     emp.Address,
		NumYearWork: emp.NumYearWork,
	}
}

func toEmployees(employees []Employee) []models.Employee {
	out := make([]models.Employee, len(employees))
	for i, b := range employees {
		out[i] = toEmployee(b)
	}
	return out
}
