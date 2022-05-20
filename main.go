package main

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type Product struct {
	Name       string      `json:"name"`
	Attributes []Attribute `json:"attributes"`
}

type Attribute interface {
	GetUniqueName() string
}

type AbstractAttribute struct {
	Name string `json:"name"`
}

type StringAttribute struct {
	AbstractAttribute `mapstructure:",squash"`
	Value             string `json:"value"`
}

type IntAttribute struct {
	AbstractAttribute `mapstructure:",squash"`
	Value             int `json:"value"`
}

func (p *Product) UnmarshalJSON(data []byte) error {

	var raw map[string]interface{}

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	p.Name = raw["name"].(string)
	p.Attributes = []Attribute{}

	for _, abstractAttribute := range raw["attributes"].([]interface{}) {
		attributeName := abstractAttribute.(map[string]interface{})["name"]

		switch attributeName {
			case "title":
				var sa StringAttribute
				err := mapstructure.Decode(abstractAttribute, &sa)
				if err != nil {
					panic(err)
				}
				p.Attributes = append(p.Attributes, sa)
			case "weight":
				var sa IntAttribute
				err := mapstructure.Decode(abstractAttribute, &sa)
				if err != nil {
					panic(err)
				}
				p.Attributes = append(p.Attributes, sa)
		}
	}

	return nil
}

func main() {
	product := Product{
		Name: "Bread",
		Attributes: []Attribute{
			StringAttribute{
				AbstractAttribute: AbstractAttribute{Name:  "title"},
				Value: "Baton",
			},
			IntAttribute{
				AbstractAttribute: AbstractAttribute{Name: "weight"},
				Value: 5,
			},
		},
	}

	fmt.Println("Initial object:", product)

	serializedProduct, err := json.Marshal(product)
	if err != nil {
		panic(err)
	}
	fmt.Println("Serialized object:", string(serializedProduct))

	var deserializedProduct Product
	err = json.Unmarshal(serializedProduct, &deserializedProduct)
	if err != nil {
		panic(err)
	}
	fmt.Println("Deserialized object", deserializedProduct)
}

func (a StringAttribute) GetUniqueName() string {
	return "String attribute name"
}

func (a IntAttribute) GetUniqueName() string {
	return "Int attribute name"
}
