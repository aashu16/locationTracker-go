package main

import (
	"log"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
)

var input string = "14,41"
var result []string
var form *Response

func BenchmarkSplitStringToNothing(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_ = strings.Split(input, ",")
	}
}

func BenchmarkSplitStringToSlice(b *testing.B) {
	b.ReportAllocs()
	var splitString []string
	for n := 0; n < b.N; n++ {
		splitString = strings.Split(input, ",")
	}
	result = splitString
}

func BenchmarkSplitStringToKnownSlice(b *testing.B) {
	b.ReportAllocs()
	splitString := make([]string, 2)
	for n := 0; n < b.N; n++ {
		splitString = strings.Split(input, ",")
	}
	result = splitString
}

func BenchmarkSplitStringToArray(b *testing.B) {
	b.ReportAllocs()
	splitString := make([]string, 0, 2)
	for n := 0; n < b.N; n++ {
		splitString = strings.Split(input, ",")
	}
	result = splitString
}

func BenchmarkGinRoute(b *testing.B) {
	b.ReportAllocs()

	db, err := sqlx.Connect("pgx", "user=postgres password=postgres dbname=entrysport sslmode=disable")
	if err != nil {
		log.Fatalln("Cannot connect to database...")
	}
	var env = &Env{db: db}

	var forms *Response

	for n := 0; n < b.N; n++ {
		forms = env.getForm()
	}
	form = forms
}
