package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func main() {
	url := "https://square-meter.herokuapp.com/api/employees"
	employees := []Employee{}
	getDataFromWeb(url, &employees)

	conn := initClickHouse("clickhouse:9000")
	fmt.Println("Ожидание...")
	http.HandleFunc("/employees", saveRequest(conn, employees))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func getDataFromWeb(url string, empList *[]Employee) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
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
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	jsonErr := json.Unmarshal(body, empList)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}

func initClickHouse(host string) *sql.DB {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{host},
		Auth: clickhouse.Auth{
			Database: "default",
			Password: "20ilona01",
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: 5 * time.Second,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	})
	return conn
}

func saveRequest(conn *sql.DB, employees []Employee) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		city := req.URL.Query().Get("city")
		numYearWork, err := strconv.ParseInt(req.URL.Query().Get("num"), 10, 64)
		if err != nil {
			fmt.Println(err.Error())
		}
		if err := saveDataToDB(conn, dataProcess(employees, numYearWork, city)); err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "<h1>Save to ClickHouse</h>")
	}
}

func dataProcess(employees []Employee, num int64, city string) []Employee {
	list := make([]Employee, 0, len(employees))
	for i := 0; i < len(employees); i++ {
		if (employees[i].NumYearWork > num) && (strings.Contains(employees[i].Address, city)) {
			list = append(list, employees[i])
		}
	}
	return list
}

func saveDataToDB(conn *sql.DB, employees []Employee) error {
	ctx := clickhouse.Context(context.Background(), clickhouse.WithSettings(clickhouse.Settings{
		"max_block_size": 10,
	}), clickhouse.WithProgress(func(p *clickhouse.Progress) {
		fmt.Println("progress: ", p)
	}))
	if err := conn.PingContext(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Catch exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return err
	}
	_, err := conn.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS employee (
			Id Int64,
			Name String,
			Phone String,
			Address String,
			NumYearWork Int64
		) engine=Memory
	`)
	if err != nil {
		return err
	}
	scope, err := conn.Begin()
	if err != nil {
		return err
	}
	{
		batch, err := scope.PrepareContext(ctx, "INSERT INTO employee (Id, Name, Phone, Address, NumYearWork)")
		if err != nil {
			return err
		}
		for _, emp := range employees {
			if _, err := batch.Exec(emp.Id, emp.Name, emp.Phone, emp.Address, emp.NumYearWork); err != nil {
				return err
			}
		}
	}
	if err := scope.Commit(); err != nil {
		return err
	}
	return err
}

type Employee struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	NumYearWork int64  `json:"numYearWork"`
}
