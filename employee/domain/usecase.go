package employee

type UseCase interface {
	SetEmployees() []byte
	SaveEmployees(employees []Employee, num int64, city string) error
	DataProcess(employees []Employee, num int64, city string) []Employee
	GetEmployees() ([]Employee, error)
}
