// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Resource interface {
	IsResource()
}

type Event struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Tags        []*Tag      `json:"tags"`
	URL         string      `json:"url"`
	Publisher   *Function   `json:"publisher"`
	Subscribers []*Function `json:"subscribers"`
}

func (Event) IsResource() {}

type Function struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Tags        []*Tag `json:"tags"`
	URL         string `json:"url"`
}

func (Function) IsResource() {}

type Tag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}