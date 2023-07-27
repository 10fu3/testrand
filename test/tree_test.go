package test

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
	"testrand/reader"
	"testrand/reader/eval"
)

func TestFailedTree(t *testing.T) {
	sample := strings.NewReader("1\n")
	read := reader.New(bufio.NewReader(sample))
	result, err := read.Read()
	if err != nil {
		t.Fatal(err)
	}
	if result.Type() == "cons_cell" {
		t.Fatal("type error")
	}
}

func TestSimpleTree(t *testing.T) {
	sample := strings.NewReader("(1 2 34 55 666)\n")
	read := reader.New(bufio.NewReader(sample))
	result, err := read.Read()
	if err != nil {
		t.Fatal(err)
	}
	if result.Type() != "cons_cell" {
		t.Fatal("type error")
	}
	consCell := (result).(eval.ConsCell)
	for _, i := range []string{"1", "2", "34", "55", "666"} {
		if consCell.Type() != "cons_cell" {
			t.Fatal("type error")
		}
		if consCell.GetCar().String() != i {
			t.Fatalf("different value be: %s actually: %s", i, consCell.GetCar().String())
		}
		consCell = consCell.GetCdr().(eval.ConsCell)
	}
	fmt.Println(result.String())
}

func TestDotTree(t *testing.T) {
	sample := strings.NewReader("(a . (b . c))\n")
	read := reader.New(bufio.NewReader(sample))
	result, err := read.Read()
	if err != nil {
		t.Fatal(err)
	}
	if result.Type() != "cons_cell" {
		t.Fatal("type error")
	}
	if result.String() != "(a b . c)" {
		t.Fatal("be: (a b . c) but got " + result.String())
	}
}
