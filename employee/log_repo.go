package employee

import (
	"MainGoTask/models"
)

type LogRepository interface {
	SaveEmployees(employees []models.Employee) error
	GetEmployees() ([]models.Employee, error)
}
