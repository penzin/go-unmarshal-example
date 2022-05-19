package main

import (
	"encoding/json"
	"fmt"
)

type Attribute interface {
	GetUniqueName() string
}

type AbstractAttribute struct {
	Name string `json:"name"`
}

type StringAttribute struct {
	AbstractAttribute
	Value string `json:"value"`
}

type IntAttribute struct {
	AbstractAttribute
	Value int `json:"value"`
}

type Product struct {
	Name       string `json:"name"`
	Attributes []Attribute `json:"attributes"`
}

func main() {
	product := Product{
		Name: "Bread",
		Attributes: []Attribute{
			StringAttribute{
				AbstractAttribute: AbstractAttribute{Name: "title"},
				Value:             "Baton",
			},
			IntAttribute{
				AbstractAttribute: AbstractAttribute{Name: "weight"},
				Value:             5,
			},
		},
	}

	fmt.Println(product)

	serializedProduct, err := json.Marshal(product)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(serializedProduct))

	var deserializedProduct Product

	err = json.Unmarshal(serializedProduct, &deserializedProduct)

	if err != nil {
		panic(err)
	}

	fmt.Println(deserializedProduct)
}

func (a StringAttribute) GetUniqueName() string {
	return "String attribute name"
}

func (a IntAttribute) GetUniqueName() string {
	return "Int attribute name"
}
