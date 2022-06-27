package employee

import (
	"MainGoTask/models"
)

type UseCase interface {
	SetEmployees() []byte
	SaveEmployees(employees []models.Employee, num int64, city string) error
	DataProcess(employees []models.Employee, num int64, city string) []models.Employee
	GetEmployees() ([]models.Employee, error)
}
