package main

import (
	"GraphQL/models"
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"log"
)

var tutorials []models.Tutorial

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

var fields = graphql.Fields{
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

var mutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"create": &graphql.Field{
			Description: "Create a new tutorial",
			Type:        tutorialType,
			Args: graphql.FieldConfigArgument{
				"Title": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tutorial := models.Tutorial{
					Title: p.Args["Title"].(string),
				}
				tutorials = append(tutorials, tutorial)
				return tutorial, nil
			},
		},
	},
})

func main() {
	// TODO include database with gorm ORM package
	tutorials = models.Populate()
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{
		Query:    graphql.NewObject(rootQuery),
		Mutation: mutationType,
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Query
	query := `
    mutation {
        create(Title: "Hello World") {
            Title
        }
    }
	`
	applyQuery(schema, query)
	query = `
	{
		list {
			ID
			Title
		}
	}
	`
	applyQuery(schema, query)

}
func applyQuery(schema graphql.Schema, query string) {
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON)
}
