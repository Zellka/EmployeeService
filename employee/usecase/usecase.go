package usecase

import (
	"MainGoTask/employee"
	"MainGoTask/models"
	"strings"
)

type EmployeekUseCase struct {
	employeeRepo employee.LogRepository
	webRepo      employee.WebRepository
}

func NewEmployeeUseCase(employeeRepo employee.LogRepository, webRepo employee.WebRepository) *EmployeekUseCase {
	return &EmployeekUseCase{
		employeeRepo: employeeRepo,
		webRepo:      webRepo,
	}
}

func (emp EmployeekUseCase) SetEmployees() []byte {
	return emp.webRepo.SetEmployees()
}

func (emp EmployeekUseCase) SaveEmployees(employees []models.Employee, num int64, city string) error {
	return emp.employeeRepo.SaveEmployees(employees)
}

func (emp EmployeekUseCase) GetEmployees() ([]models.Employee, error) {
	return emp.employeeRepo.GetEmployees()
}

func (emp EmployeekUseCase) DataProcess(employees []models.Employee, num int64, city string) []models.Employee {
	list := make([]models.Employee, 0, len(employees))
	for i := 0; i < len(employees); i++ {
		if (employees[i].NumYearWork > num) && (strings.HasPrefix(employees[i].Address, city)) {
			list = append(list, employees[i])
		}
	}
	return list
}
