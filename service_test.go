package main

import (
	broker "MainGoTask/employee/delivery"
	domain "MainGoTask/employee/domain"
	rep "MainGoTask/employee/repository"
	usecase "MainGoTask/employee/usecase"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
)

func TestRequest(t *testing.T) {
	composeFilePaths := []string{"docker-compose.yml"}
	compose := tc.NewLocalDockerCompose(composeFilePaths, "my_app")
	upCompose(compose, composeFilePaths)

	db := initClickHouse("localhost:9001")
	url := "https://square-meter.herokuapp.com/api/employees"
	logRepo := rep.NewEmployeeRepository(db)
	webRepo := rep.NewWebRepository(url)
	usecase := usecase.NewEmployeeUseCase(logRepo, webRepo)
	br := broker.NewBroker(usecase)

	checkHTTPResponse(t, br)
	checkSaveToDB(t, br)

	defer destroyCompose(t, compose)
}

func upCompose(compose *tc.LocalDockerCompose, composeFilePaths []string) {
	execError := compose.
		WithCommand([]string{"up", "-d"}).
		WithEnv(map[string]string{
			"CLICKHOUSE_PASSWORD": "20ilona01",
		}).
		Invoke()
	err := execError.Error
	if err != nil {
		log.Fatalf("Could not run compose file: %v - %v\n", composeFilePaths, err)
	}
}

func destroyCompose(t *testing.T, compose *tc.LocalDockerCompose) {
	err := compose.Down()
	if err.Error != nil {
		t.Fatal(err.Error)
	}
}

func checkHTTPResponse(t *testing.T, b *broker.Broker) {
	employees := broker.ParseEmployee(b.SetEmployees())

	req, err := http.NewRequest("GET", "/employees", nil)
	if err != nil {
		t.Fatal(err)
	}
	cities := []string{"Донецк", "Макеевка"}
	for i := 0; i < 6; i++ {
		for _, city := range cities {
			q := req.URL.Query()
			q.Add("city", city)
			q.Add("num", strconv.Itoa(i))
			req.URL.RawQuery = q.Encode()

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(CheckHandler)
			handler.ServeHTTP(rr, req)
			require.Equal(t, http.StatusOK, rr.Code, "HTTP status code")

			if status := b.SaveEmployees(employees, int64(i), city); status != http.StatusOK {
				t.Fatal(status)
			}
		}
	}
}

func checkSaveToDB(t *testing.T, b *broker.Broker) {
	employees, status := b.GetEmployees()
	if status != http.StatusOK {
		t.Fatal(status)
	}
	checkTrueParseData(t, employees)
}

func checkTrueParseData(t *testing.T, employees []domain.Employee) {
	employee := domain.Employee{Id: 2, Name: "Федосеев Владислав", Phone: "0713504125", Address: "Донецк, ул.Кирова, 255", NumYearWork: 2}
	if employee != employees[0] {
		t.Fatal("Error parse data: ", employees[0])
	}
}

func CheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
