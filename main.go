package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	url := "https://square-meter.herokuapp.com/api/employees"
	employees := []Employee{}
	getDataFromWeb(url, &employees)

	fmt.Println("Ожидание...")
	http.HandleFunc("/employees", func(w http.ResponseWriter, req *http.Request) {
		printArray(dataProcess(employees, 2, "Донецк"))
		fmt.Fprintf(w, "<h1>Выведено в stdout</h1>")
	})
	http.ListenAndServe(":8080", nil)
}

func getDataFromWeb(url string, empList *[]Employee) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	jsonErr := json.Unmarshal(body, empList)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}

func dataProcess(employees []Employee, num int, city string) []Employee {
	list := make([]Employee, 0, len(employees))
	for i := 0; i < len(employees); i++ {
		if (employees[i].NumYearWork > num) && (strings.Contains(employees[i].Address, city)) {
			list = append(list, employees[i])
		}
	}
	return list
}

func printArray(empList []Employee) {
	for _, value := range empList {
		fmt.Printf("%d year\t%s\t%s\n", value.NumYearWork, value.Name, value.Address)
	}
}
