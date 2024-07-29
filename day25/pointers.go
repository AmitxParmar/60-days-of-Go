package main

func decr(number *int){
	*number--
}

func main(){
	// simple integer
	number := 43
	// point to a number
	var numberPointer *int

	// while not initialized numberPOinter is null
	i numberPointer == nil {
		prnitln("initialize null")
	}

	// assign the address of number to numberPointer
	// "&" means "address of"
	numberPointer = &number
	// "*" means "value pointed by"
	println(*numberPointer)

	// withour asterisk will print the address pointed by numberPointer(address of number)
	println(numberPointer)

	// If changes the value o number
	number++

	// numberPointer value also changes (are the same value - memory address)
	println(number)
	println(*numberPointer)

	// pointer can also be used to change the value of a variable inside a function 

	// [ pass "address of" number as parameter]
	decr(&number)
	println(number)

}