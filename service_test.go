package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
)

func TestRequest(t *testing.T) {
	composeFilePaths := []string{"docker-compose.yml"}
	compose := tc.NewLocalDockerCompose(composeFilePaths, "my_app")
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
	conn := initClickHouse("localhost:9001")
	checkHTTPResponse(t)
	checkSaveToDB(conn, t)

	destroyCompose := func() {
		err := compose.Down()
		if err.Error != nil {
			log.Fatal(err.Error)
		}
	}
	defer destroyCompose()
}

func checkHTTPResponse(t *testing.T) {
	url := "https://square-meter.herokuapp.com/api/employees"
	employees := []Employee{}
	getDataFromWeb(url, &employees)

	req, err := http.NewRequest("GET", "/employees", nil)
	if err != nil {
		t.Fatal(err)
	}
	conn := initClickHouse("localhost:9001")
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
			if err := saveDataToDB(conn, dataProcess(employees, int64(i), city)); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func checkSaveToDB(conn *sql.DB, t *testing.T) {
	ctx := clickhouse.Context(context.Background(), clickhouse.WithSettings(clickhouse.Settings{
		"max_block_size": 10,
	}), clickhouse.WithProgress(func(p *clickhouse.Progress) {
		fmt.Println("progress: ", p)
	}))
	if err := conn.PingContext(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			t.Fatalf("Catch exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		t.Fatal(err)
	}
	rows, err := conn.QueryContext(ctx, "SELECT * FROM employee", 0, "xxx", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if !rows.Next() {
		t.Fatal("data base is empty")
	}
	checkTrueParseData(t, rows)
	rows.Close()
}

func checkTrueParseData(t *testing.T, rows *sql.Rows) {
	empList := []Employee{}
	for rows.Next() {
		var (
			id          int64
			name        string
			phone       string
			address     string
			numYearWork int64
		)
		if err := rows.Scan(&id, &name, &phone, &address, &numYearWork); err != nil {
			t.Fatal(err)
		}
		empList = append(empList, Employee{id, name, phone, address, numYearWork})
	}
	employee := Employee{Id: 17, Name: "Никифорова Анастасия", Phone: "0716404125", Address: "Донецк, ул.Бирюзова, 25/47", NumYearWork: 5}
	if employee != empList[0] {
		t.Fatal("Error parse data")
	}
}

func CheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
