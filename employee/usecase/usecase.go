package usecase

import (
	rep "MainGoTask/employee/repository/interfaces"
	model "MainGoTask/model"
	"log"
	"strings"
)

type EmployeekUseCase struct {
	employeeRepo rep.LogRepository
	webRepo      rep.WebRepository
}

func NewEmployeeUseCase(employeeRepo rep.LogRepository, webRepo rep.WebRepository) *EmployeekUseCase {
	return &EmployeekUseCase{
		employeeRepo: employeeRepo,
		webRepo:      webRepo,
	}
}

func (emp EmployeekUseCase) GetEmployeesFromWeb() []model.Employee {
	return emp.webRepo.GetEmployeesFromWeb()
}

func (emp EmployeekUseCase) SaveEmployees(employees []model.Employee) error {
	return emp.employeeRepo.SaveEmployees(employees)
}

func (emp EmployeekUseCase) GetEmployees(num int64, city string) ([]model.Employee, error) {
	employees := emp.GetEmployeesFromWeb()
	if err := emp.SaveEmployees(emp.DataProcess(employees, num, city)); err != nil {
		log.Fatal(err)
	}
	return emp.employeeRepo.GetEmployees()
}

func (emp EmployeekUseCase) DataProcess(employees []model.Employee, num int64, city string) []model.Employee {
	list := make([]model.Employee, 0, len(employees))
	for i := 0; i < len(employees); i++ {
		if (employees[i].NumYearWork > num) && (strings.HasPrefix(employees[i].Address, city)) {
			list = append(list, employees[i])
		}
	}
	return list
}
