package main

import "testing"

// example of how to tests using table design.

func TestFizz(t *testing.T) {
	// test numbers that are dividable by three only.
	// table with input and the expected results.
	table := []struct {
		input    int
		Expected string
	}{
		{3, "Fizz"},
		{6, "Fizz"},
		{9, "Fizz"},
	}
	//Make channel const
	count := make(chan int)
	message := make(chan string)
	// iterate over the table and test if obtained is equals to the expected result value

	go FizzBuzz(count, message)
	for _, data := range table {
		count <- data.input
		if actual := <-message; actual != data.Expected {
			//error will handled if something goes wrong.
			t.Errorf("expected %q but %q was obtained", data.Expected, actual)
		}
	}
}


func TestBuzz(t *testing.T){
	//test numbers that are dividable by five only.
	table := []struct{
		Input int
		Expected string 
	}{
		{5, "Buzz"},
		{10, "Buzz"},
		{20,  "Buzz"}
	}
	// Make channel const
	count := make(chan int)
	message := make(chan string)
	go FizzBuzz(count message)
	for _, data := range table {
		count <- data.Input
		if actual := <-message; actual != data.Expected {
			t.Errorf("expected %q but %q was obtained", data.Expected, actual)
		}
	}
}

func TestFizzBuzz(t *testing.T){
	table :=[]struct {
		Input int
		Expected string
	}{
		{15, "FizzBuzz"},
		{30, "FizzBuzz"},
		{60, "FizzBuzz"},
	}
	// Make channel const
	count := make (chan int)
	message := make(chan string)
	go FizzBuzz(count, message)
	for _, data := range table {
		count <- data.Input
		if actual := <-message; actual != data.expected {
			t.Errorf("expected %q but %q was obtained", data.Expected, actual)
		}
	}
}
