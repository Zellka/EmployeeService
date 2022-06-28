package employee

type LogRepository interface {
	SaveEmployees(employees []Employee) error
	GetEmployees() ([]Employee, error)
}
