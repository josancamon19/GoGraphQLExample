package main

import (
	"GraphQL/models"
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"log"
)

func main() {
	var commentType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Comment",
		Fields: graphql.Fields{
			"Body": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
	var authorType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Author",
		Fields: graphql.Fields{
			"Name": &graphql.Field{
				Type: graphql.String,
			},
			"Tutorials": &graphql.Field{
				Type: graphql.NewList(graphql.Int),
			},
		},
	})

	var tutorialType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Tutorial", Fields: graphql.Fields{
			"ID": &graphql.Field{
				Type: graphql.Int,
			},
			"Title": &graphql.Field{
				Type: graphql.String,
			},
			"Author": &graphql.Field{
				Type: authorType,
			},
			"Comments": &graphql.Field{
				Type: graphql.NewList(commentType),
			},
		},
	})

	fields := graphql.Fields{
		"tutorial": &graphql.Field{
			Type:        tutorialType,
			Description: "Get tutorial by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					tutorialsFakeData := models.Populate()
					for _, tutorial := range tutorialsFakeData {
						if tutorial.ID == id {
							return tutorial, nil
						}
					}
				}
				return nil, nil
			},
		},
		"list": &graphql.Field{
			Type:        graphql.NewList(tutorialType),
			Description: "Get tutorial list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return models.Populate(), nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Doing query
//	query := `
//    {
//        list {
//            ID
//            Title
//            Comments {
//                Body
//            }
//            Author {
//                Name
//                Tutorials
//            }
//        }
//    }
//`
	query := `
    {
        tutorial(id:1) {
            Title
            Author {
                Name
                Tutorials
            }
        }
    }
`
	params := graphql.Params{Schema: schema, RequestString: query}
	response := graphql.Do(params)
	if len(response.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", response.Errors)
	}

	rJSON, _ := json.Marshal(response)
	fmt.Printf("%s \n", rJSON)
}

func SimpleExample() {
	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
			Description: "Test field",
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Doing query
	query := `
		{
			hello
		}
	`
	params := graphql.Params{Schema: schema, RequestString: query}
	response := graphql.Do(params)
	if len(response.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", response.Errors)
	}

	rJSON, _ := json.Marshal(response)
	fmt.Printf("%s \n", rJSON) // {“data”:{“hello”:”world”}}
}
