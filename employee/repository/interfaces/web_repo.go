package interfaces

import model "MainGoTask/model"

type WebRepository interface {
	GetEmployeesFromWeb() []model.Employee
}
