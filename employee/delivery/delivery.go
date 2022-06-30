package delivery

import (
	usecase "MainGoTask/employee/usecase/interfaces"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Broker struct {
	useCase usecase.GetEmployeesUseCase
	router  *mux.Router
}

func NewBroker(useCase usecase.GetEmployeesUseCase) *Broker {
	return &Broker{
		useCase: useCase,
		router:  mux.NewRouter(),
	}
}

func (b *Broker) Start() {
	fmt.Println("Ожидание...")

	b.router.HandleFunc("/employees", b.HandleGetEmployee)
	http.Handle("/", b.router)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (b *Broker) HandleGetEmployee(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	city := req.URL.Query().Get("city")
	numYearWork, err := strconv.ParseInt(req.URL.Query().Get("num"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	employees, err := b.useCase.GetEmployees(numYearWork, city)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	jsonResp, err := json.Marshal(employees)
	if err != nil {
		log.Fatalf("Error JSON marshal. Err: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(jsonResp)
}
