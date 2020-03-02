package main

import (
	"strings"
	"testing"
)

var input string = "14,41"
var result []string

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
