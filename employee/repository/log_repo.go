package repository

import (
	model "MainGoTask/model"
	"database/sql"
	"log"
	"time"
)

type EmployeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS employee (
			Id Int64,
			Name String,
			Phone String,
			Address String,
			NumYearWork Int64
		) engine=Memory
	`)
	if err != nil {
		log.Fatal(err)
	}
	return &EmployeeRepository{
		db: db,
	}
}

func (r EmployeeRepository) SaveEmployees(employees []model.Employee) error {
	scope, err := r.db.Begin()
	models := toDTOModels(employees)
	if err != nil {
		return err
	}
	batch, err := scope.Prepare("INSERT INTO employee (Id, Name, Phone, Address, NumYearWork)")
	if err != nil {
		return err
	}
	for _, emp := range models {
		if _, err := batch.Exec(emp.Id, emp.Name, emp.Phone, emp.Address, emp.NumYearWork); err != nil {
			return err
		}
	}
	if err := scope.Commit(); err != nil {
		return err
	}
	return err
}

func (r EmployeeRepository) GetEmployees() ([]model.Employee, error) {
	rows, err := r.db.Query("SELECT * FROM employee", 0, "xxx", time.Now())
	if err != nil {
		return nil, err
	}
	employees := []model.Employee{}
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
		employees = append(employees, model.Employee{id, name, phone, address, numYearWork})
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

func toDTOModels(es []model.Employee) []*EmployeeDTO {
	out := make([]*EmployeeDTO, len(es))
	for i, b := range es {
		dto := EmployeeDTO(b)
		out[i] = &dto
	}
	return out
}
