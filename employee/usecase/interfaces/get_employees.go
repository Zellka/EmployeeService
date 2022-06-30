package interfaces

import (
	model "MainGoTask/model"
)

type GetEmployeesUseCase interface {
	SaveEmployees(employees []model.Employee) error
	DataProcess(employees []model.Employee, num int64, city string) []model.Employee
	GetEmployees(num int64, city string) ([]model.Employee, error)
}
