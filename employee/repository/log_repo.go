package repository

import (
	"MainGoTask/models"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type EmployeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{
		db: db,
	}
}

func (r EmployeeRepository) SaveEmployees(employees []models.Employee) error {
	ctx := clickhouse.Context(context.Background(), clickhouse.WithSettings(clickhouse.Settings{
		"max_block_size": 10,
	}))
	if err := r.db.PingContext(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Catch exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return err
	}
	_, err := r.db.ExecContext(ctx, `
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
	scope, err := r.db.Begin()
	models := toModels(employees)
	if err != nil {
		return err
	}
	{
		batch, err := scope.PrepareContext(ctx, "INSERT INTO employee (Id, Name, Phone, Address, NumYearWork)")
		if err != nil {
			return err
		}
		for _, emp := range models {
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

func (r EmployeeRepository) GetEmployees() ([]models.Employee, error) {
	ctx := clickhouse.Context(context.Background(), clickhouse.WithSettings(clickhouse.Settings{
		"max_block_size": 10,
	}), clickhouse.WithProgress(func(p *clickhouse.Progress) {
		fmt.Println("progress: ", p)
	}))
	if err := r.db.PingContext(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			log.Fatalf("Catch exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM employee", 0, "xxx", time.Now())
	if err != nil {
		return nil, err
	}
	employees := []models.Employee{}
	for rows.Next() {
		var (
			id          int64
			name        string
			phone       string
			address     string
			numYearWork int64
		)
		if err := rows.Scan(&id, &name, &phone, &address, &numYearWork); err != nil {
			return nil, err
		}
		employees = append(employees, models.Employee{id, name, phone, address, numYearWork})
	}
	rows.Close()
	return employees, rows.Err()
}

type EmployeeDTO struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	NumYearWork int64  `json:"numYearWork"`
}

func toModel(emp models.Employee) *EmployeeDTO {
	return &EmployeeDTO{
		Id:          emp.Id,
		Name:        emp.Name,
		Phone:       emp.Phone,
		Address:     emp.Address,
		NumYearWork: emp.NumYearWork,
	}
}

func toModels(es []models.Employee) []*EmployeeDTO {
	out := make([]*EmployeeDTO, len(es))

	for i, b := range es {
		out[i] = toModel(b)
	}

	return out
}
