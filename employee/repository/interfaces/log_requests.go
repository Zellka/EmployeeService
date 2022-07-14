package interfaces

import (
	model "MainGoTask/model"
)

type LogRequests interface {
	SaveEmployees(employees []model.Employee) error
	GetEmployees() ([]model.Employee, error)
}
