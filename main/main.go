package main

import (
	"net/http"

	"../office_db"
	"../schema"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func main() {
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{Query: office_schema.RootQuery, Mutation: office_schema.RootMutation})
	defer officedb.CloseDB()

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
	http.Handle("/graphql", h)
	http.ListenAndServe(":8083", nil)
}
