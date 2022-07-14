package repository

import (
	model "MainGoTask/model"
	"database/sql"
	"log"
	"time"
)

type LogRepository struct {
	db *sql.DB
}

func NewLogRepository(db *sql.DB) *LogRepository {
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
	return &LogRepository{
		db: db,
	}
}

func (r LogRepository) SaveEmployees(employees []model.Employee) error {
	scope, err := r.db.Begin()
	emploeesDTO := toDTOModels(employees)
	if err != nil {
		return err
	}
	batch, err := scope.Prepare("INSERT INTO employee (Id, Name, Phone, Address, NumYearWork)")
	if err != nil {
		return err
	}
	for _, emp := range emploeesDTO {
		if _, err := batch.Exec(emp.Id, emp.Name, emp.Phone, emp.Address, emp.NumYearWork); err != nil {
			return err
		}
	}
	if err := scope.Commit(); err != nil {
		return err
	}
	return err
}

func (r LogRepository) GetEmployees() ([]model.Employee, error) {
	rows, err := r.db.Query("SELECT * FROM employee", 0, "xxx", time.Now())
	if err != nil {
		return nil, err
	}
	emp := model.Employee{}
	employees := []model.Employee{}
	for rows.Next() {
		if err := rows.Scan(&emp.Id, &emp.Name, &emp.Phone, &emp.Address, &emp.NumYearWork); err != nil {
			return nil, err
		}
		employees = append(employees, emp)
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
