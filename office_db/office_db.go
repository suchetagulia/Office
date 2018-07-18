package officedb

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Company struct {
	ID        string
	Name      string
	Offices   []Office
	Employees []Person
	Director  Person
}

type Office struct {
	ID        string
	City      string
	CompanyID string
	Employees []Person
}

type Person struct {
	ID          string `json:"id"`
	Name        string `json: "name"`
	Designation string `json: "designation"`
	OfficeID    string `json: "officeId"`
}

var db *sql.DB
var err error

func init() {
	db, err = sql.Open("mysql", "root:asdfg@tcp(127.0.0.1:3306)/")
	if err != nil {
		fmt.Println("error occurred while connecting to db: ", err)
	}

	_, err := db.Exec("USE officedb")
	if err != nil {
		fmt.Println("error occurred while using db: ", err)
	}
}

func GetPerson(id string) *Person {
	stmt := "SELECT id, name, designation, officeid FROM person WHERE id = '" + id + "'"
	rows, err := db.Query(stmt)
	if err != nil {
		fmt.Println("error occurred while executing stmt: ", err)
	}
	var p Person
	for rows.Next() {
		err = rows.Scan(&p.ID, &p.Name, &p.Designation, &p.OfficeID)
		if err != nil {
			fmt.Println("error occurred while executing rows.Scan of person: ", err)
		}
	}
	return &p
}

func GetOffice(id string) *Office {
	stmt := "select * from office where id = '" + id + "'"
	rows, err := db.Query(stmt)
	if err != nil {
		fmt.Println("error occurred while executing office stmt: ", err)
	}
	var off Office
	for rows.Next() {
		err = rows.Scan(&off.ID, &off.City, &off.CompanyID)
		if err != nil {
			fmt.Println("error occurred while executing rows.Scan of office: ", err)
		}
	}
	off.Employees = GetOfficeEmp(id)
	return &off
}

func GetOfficeEmp(id string) []Person {
	var employees []Person
	stmt := "select * from person where officeid = '" + id + "'"
	rows, err := db.Query(stmt)
	if err != nil {
		fmt.Println("error occurred while executing GetOfficeEmp stmt: ", stmt, " err: ", err)
	}
	for rows.Next() {
		var emp Person
		err = rows.Scan(&emp.ID, &emp.Name, &emp.Designation, &emp.OfficeID)
		if err != nil {
			fmt.Println("error occurred while executing rows.Scan of GetOfficeEmp: ", err)
		}
		employees = append(employees, emp)
	}
	return employees
}

func GetCompany(id string) Company {
	stmt := "select * from company where id = '" + id + "'"
	rows, err := db.Query(stmt)
	if err != nil {
		fmt.Println("GetCompany: ", err)
	}
	var com Company
	for rows.Next() {
		err = rows.Scan(&com.ID, &com.Name)
		if err != nil {
			fmt.Println("error in rows.Scan of GetCompany: ", err)
		}
	}
	com.Offices, com.Employees = GetCompanyOff(id)
	return com
}

func GetCompanyOff(id string) ([]Office, []Person) {
	var offs []Office
	var offices []Office
	var employees []Person
	stmt := "select * from office where cid = '" + id + "'"
	rows, err := db.Query(stmt)
	if err != nil {
		fmt.Println("GetCompanyOff: ", err)
	}
	for rows.Next() {
		var off Office
		err = rows.Scan(&off.ID, &off.City, &off.CompanyID)
		if err != nil {
			fmt.Println("rows.Scan of GetCompanyOff: ", err)
		}
		offs = append(offs, off)
	}
	for _, off := range offs {
		off.Employees = GetOfficeEmp(off.ID)
		employees = append(employees, off.Employees...)
		offices = append(offices, off)
	}
	return offices, employees
}

func AddPerson(person Person) {
	stmt := "INSERT INTO person VALUES ('" + person.ID + "', '" + person.Name + "', '" + person.Designation + "', '" + person.OfficeID + "')"
	_, err := db.Exec(stmt)
	if err != nil {
		fmt.Println("AddPerson: ", stmt, " err: ", err)
	}
}

func AddOffice(office Office) {
	stmt := "INSERT INTO office VALUES ('" + office.ID + "', '" + office.City + "', '" + office.CompanyID + "')"
	_, err := db.Exec(stmt)
	if err != nil {
		fmt.Println("AddOffice: ", stmt, " err: ", err)
	}
	for _, emp := range office.Employees {
		AddPerson(emp)
	}
}

func AddCompany(company Company) {
	stmt := "INSERT INTO company VALUES ('" + company.ID + "', '" + company.Name + "')"
	_, err := db.Exec(stmt)
	if err != nil {
		fmt.Println("AddCompany: ", stmt, " err: ", err)
	}
	for _, office := range company.Offices {
		AddOffice(office)
	}
	for _, emp := range company.Employees {
		AddPerson(emp)
	}
	AddPerson(company.Director)
}

func CloseDB() {
	db.Close()
}
