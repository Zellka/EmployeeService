package main

import (
	broker "MainGoTask/employee/delivery"
	rep "MainGoTask/employee/repository"
	usecase "MainGoTask/employee/usecase"
	"database/sql"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

//http://localhost:80/employees?city=Донецк&num=4

func main() {
	db := initClickHouse("clickhouse:9000")
	url := "https://square-meter.herokuapp.com/api/employees"

	logRepo := rep.NewLogRepository(db)
	webRepo := rep.NewWebRepository(url)

	usecase := usecase.NewGetEmployeesUseCase(logRepo, webRepo)
	br := broker.NewBroker(usecase)
	br.Start()
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
