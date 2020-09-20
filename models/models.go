package models

type Tutorial struct {
	ID int `graphql:"id"`
	Title string
	Author Author
	Comments []Comment
}

type Author struct {
	Name string
	Tutorials []int
}

type Comment struct {
	Body string
}

func Populate() []Tutorial {
	author := &Author{Name: "Elliot Forbes", Tutorials: []int{1}}
	tutorial := Tutorial{
		ID: 1,
		Title:  "Go GraphQL Tutorial",
		Author: *author,
		Comments: []Comment{
			{Body: "First Comment"},
		},
	}

	var tutorials []Tutorial
	tutorials = append(tutorials, tutorial)

	return tutorials
}