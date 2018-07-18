package office_schema

import (
	"../office_db"
	"github.com/graphql-go/graphql"
)

var personType *graphql.Object
var officeType *graphql.Object
var companyType *graphql.Object

var personInput *graphql.InputObject
var officeInput *graphql.InputObject
var companyInput *graphql.InputObject

func init() {

	personType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Person",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"designation": &graphql.Field{
				Type: graphql.String,
			},
			"officeId": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	personInput = graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "PersonInputType",
		Fields: graphql.InputObjectConfigFieldMap{
			"id": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"name": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"designation": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"officeId": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
		},
	})

	officeType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Office",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"city": &graphql.Field{
				Type: graphql.String,
			},
			"companyId": &graphql.Field{
				Type: graphql.String,
			},
			"employees": &graphql.Field{
				Type: graphql.NewList(personType),
			},
		},
	})

	officeInput = graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "OfficeInputType",
		Fields: graphql.InputObjectConfigFieldMap{
			"id": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"city": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"companyId": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"employees": &graphql.InputObjectFieldConfig{
				Type: graphql.NewList(personInput),
			},
		},
	})

	companyType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Company",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"offices": &graphql.Field{
				Type: graphql.NewList(officeType),
			},
			"employees": &graphql.Field{
				Type: graphql.NewList(personType),
			},
			"director": &graphql.Field{
				Type: personType,
			},
		},
	})

	companyInput = graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CompanyInputType",
		Fields: graphql.InputObjectConfigFieldMap{
			"id": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"name": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"offices": &graphql.InputObjectFieldConfig{
				Type: graphql.NewList(officeInput),
			},
			"employees": &graphql.InputObjectFieldConfig{
				Type: graphql.NewList(personInput),
			},
			"director": &graphql.InputObjectFieldConfig{
				Type: personInput,
			},
		},
	})
}

var RootQuery *graphql.Object
var RootMutation *graphql.Object

func init() {
	RootQuery = graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"person": &graphql.Field{
				Type: personType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id := params.Args["id"].(string)
					person := officedb.GetPerson(id)
					return person, nil
				},
			},
			"office": &graphql.Field{
				Type: officeType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id := params.Args["id"].(string)
					office := officedb.GetOffice(id)
					return office, nil
				},
			},
			"company": &graphql.Field{
				Type: companyType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id := params.Args["id"].(string)
					company := officedb.GetCompany(id)
					return company, nil
				},
			},
		},
	})
}

func init() {
	RootMutation = graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createPerson": &graphql.Field{
				Type: personType,
				Args: graphql.FieldConfigArgument{
					"person": &graphql.ArgumentConfig{
						Type: personInput,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					personI := parsePerson(params.Args["person"].(map[string]interface{}))
					officedb.AddPerson(personI)
					personO := officedb.GetPerson(personI.ID)
					return personO, nil
				},
			},
			"createOffice": &graphql.Field{
				Type: officeType,
				Args: graphql.FieldConfigArgument{
					"office": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(officeInput),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					officeI := parseOffice(params.Args["office"].(map[string]interface{}))
					officedb.AddOffice(officeI)
					officeO := officedb.GetOffice(officeI.ID)
					return officeO, nil
				},
			},
			"createCompany": &graphql.Field{
				Type: companyType,
				Args: graphql.FieldConfigArgument{
					"company": &graphql.ArgumentConfig{
						Type: companyInput,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					companyI := parseCompany(params.Args["company"].(map[string]interface{}))
					officedb.AddCompany(companyI)
					companyO := officedb.GetCompany(companyI.ID)
					return companyO, nil
				},
			},
		},
	})
}

func parseCompany(cmp map[string]interface{}) officedb.Company {
	cid := cmp["id"].(string)
	offices := parseOffices(cmp["offices"].([]interface{}), &cid)
	var officeIDs []string
	for _, office := range offices {
		officeIDs = append(officeIDs, office.ID)
	}
	company := officedb.Company{
		ID:        cmp["id"].(string),
		Name:      cmp["name"].(string),
		Director:  parsePerson(cmp["director"].(map[string]interface{})),
		Offices:   offices,
		Employees: parsePeople(cmp["employees"].([]interface{}), officeIDs),
	}
	return company
}

func parsePerson(pmp map[string]interface{}) officedb.Person {
	person := officedb.Person{
		ID:          pmp["id"].(string),
		Name:        pmp["name"].(string),
		Designation: pmp["designation"].(string),
		OfficeID:    pmp["officeId"].(string),
	}
	return person
}

func parseOffices(osmp []interface{}, cid *string) []officedb.Office {
	var offices []officedb.Office
	for _, tmp := range osmp {
		omp := tmp.(map[string]interface{})
		office := parseOffice(omp)
		if cid != nil && office.CompanyID != *cid {
			panic("companyId mismatch error")
		}
		offices = append(offices, office)
	}
	return offices
}

func parseOffice(omp map[string]interface{}) officedb.Office {
	var officeIDs []string
	officeIDs = append(officeIDs, omp["id"].(string))
	office := officedb.Office{
		ID:        omp["id"].(string),
		City:      omp["city"].(string),
		CompanyID: omp["companyId"].(string),
		Employees: parsePeople(omp["employees"].([]interface{}), officeIDs),
	}
	return office
}

func parsePeople(pmap []interface{}, officeIDs []string) []officedb.Person {
	var people []officedb.Person
	for _, tmp := range pmap {
		pmp := tmp.(map[string]interface{})
		person := parsePerson(pmp)
		notMatched := true
		if len(officeIDs) == 0 {
			notMatched = false
		}
		for _, officeID := range officeIDs {
			if person.OfficeID == officeID {
				notMatched = false
				break
			}
		}
		if notMatched {
			panic("OfficeID mismatch error")
		}
		people = append(people, person)
	}
	return people
}
