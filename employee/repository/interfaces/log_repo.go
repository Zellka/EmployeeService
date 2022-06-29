package interfaces

import (
	model "MainGoTask/model"
)

type LogRepository interface {
	SaveEmployees(employees []model.Employee) error
	GetEmployees() ([]model.Employee, error)
}
