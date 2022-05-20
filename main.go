package main

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type Product struct {
	ProductName string      `json:"product-name" mapstructure:"product-name"`
	Attributes  []Attribute `json:"attributes"`
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

	var rawProductData, rawProductDataWithoutAttributes map[string]interface{}

	err := json.Unmarshal(data, &rawProductData)
	if err != nil {
		return err
	}

	rawProductDataWithoutAttributes = make(map[string]interface{})
	for key, value := range rawProductData {
		if key == "attributes" {
			continue
		}
		rawProductDataWithoutAttributes[key] = value
	}

	err = mapstructure.Decode(rawProductDataWithoutAttributes, &p)
	if err != nil {
		return err
	}

	p.Attributes = []Attribute{}

	for _, abstractAttribute := range rawProductData["attributes"].([]interface{}) {
		p.Attributes = append(p.Attributes, deserializeAttribute(abstractAttribute))
	}

	return nil
}

func main() {
	product := Product{
		ProductName: "Bread",
		Attributes: []Attribute{
			StringAttribute{
				AbstractAttribute: AbstractAttribute{Name: "title"},
				Value:             "Baguette",
			},
			IntAttribute{
				AbstractAttribute: AbstractAttribute{Name: "weight"},
				Value:             500,
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

func deserializeAttribute(unknownAttribute interface{}) Attribute {
	attributeName := unknownAttribute.(map[string]interface{})["name"]

	if attributeName == "title" {
		var sa StringAttribute
		err := mapstructure.Decode(unknownAttribute, &sa)
		if err != nil {
			panic(err)
		}
		return sa
	}

	var sa IntAttribute
	err := mapstructure.Decode(unknownAttribute, &sa)
	if err != nil {
		panic(err)
	}
	return sa
}
