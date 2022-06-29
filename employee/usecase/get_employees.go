package usecase

import (
	rep "MainGoTask/employee/repository/interfaces"
	model "MainGoTask/model"
	"log"
	"strings"
)

type GetEmployeesUseCase struct {
	employeeRepo rep.LogRequests
	webRepo      rep.WebRepository
}

func NewGetEmployeesUseCase(employeeRepo rep.LogRequests, webRepo rep.WebRepository) *GetEmployeesUseCase {
	return &GetEmployeesUseCase{
		employeeRepo: employeeRepo,
		webRepo:      webRepo,
	}
}

func (emp GetEmployeesUseCase) SaveEmployees(employees []model.Employee) error {
	return emp.employeeRepo.SaveEmployees(employees)
}

func (emp GetEmployeesUseCase) GetEmployees(num int64, city string) ([]model.Employee, error) {
	employees := emp.webRepo.GetEmployeesFromWeb()
	if err := emp.SaveEmployees(emp.DataProcess(employees, num, city)); err != nil {
		log.Fatal(err)
	}
	return emp.employeeRepo.GetEmployees()
}

func (emp GetEmployeesUseCase) DataProcess(employees []model.Employee, num int64, city string) []model.Employee {
	empListReady := make([]model.Employee, 0, len(employees))
	for i := 0; i < len(employees); i++ {
		if (employees[i].NumYearWork > num) && (strings.HasPrefix(employees[i].Address, city)) {
			empListReady = append(empListReady, employees[i])
		}
	}
	return empListReady
}
