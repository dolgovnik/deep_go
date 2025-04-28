package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

var (
	propertiesTag = "properties"
	omitemptyTag  = "omitempty"
)

func Serialize(person Person) string {
	persT := reflect.TypeOf(person)	
	persV := reflect.ValueOf(person)	

	res := strings.Builder{}

	for i := 0; i < persT.NumField(); i++ {
		tag := strings.Split(persT.Field(i).Tag.Get(propertiesTag), ",")

		omitempty := false
		if len(tag) > 1 && tag[1] == omitemptyTag{
			omitempty = true
		}

		val := persV.Field(i)
		if val.IsZero() && omitempty {
			continue
		}

		if res.Len() > 0 {
			res.WriteByte('\n')
		}

		res.WriteString(tag[0])
		res.WriteByte('=')
		res.WriteString(fmt.Sprintf("%v", val))
	}

	return res.String()

}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
