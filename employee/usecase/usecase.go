package usecase

import (
	domain "MainGoTask/employee/domain"
	"strings"
)

type EmployeekUseCase struct {
	employeeRepo domain.LogRepository
	webRepo      domain.WebRepository
}

func NewEmployeeUseCase(employeeRepo domain.LogRepository, webRepo domain.WebRepository) *EmployeekUseCase {
	return &EmployeekUseCase{
		employeeRepo: employeeRepo,
		webRepo:      webRepo,
	}
}

func (emp EmployeekUseCase) SetEmployees() []byte {
	return emp.webRepo.SetEmployees()
}

func (emp EmployeekUseCase) SaveEmployees(employees []domain.Employee, num int64, city string) error {
	return emp.employeeRepo.SaveEmployees(employees)
}

func (emp EmployeekUseCase) GetEmployees() ([]domain.Employee, error) {
	return emp.employeeRepo.GetEmployees()
}

func (emp EmployeekUseCase) DataProcess(employees []domain.Employee, num int64, city string) []domain.Employee {
	list := make([]domain.Employee, 0, len(employees))
	for i := 0; i < len(employees); i++ {
		if (employees[i].NumYearWork > num) && (strings.HasPrefix(employees[i].Address, city)) {
			list = append(list, employees[i])
		}
	}
	return list
}
