package main

import (
	"fmt"

	"github.com/AaronSaikovski/gostarter/constants"
	"github.com/AaronSaikovski/gostarter/pkg/samplemodule"
	"github.com/AaronSaikovski/gostarter/types"
)

func main() {
	fmt.Printf("This is the main function.\n")

	for i := 0; i < constants.LoopConstant; i++ {
		fmt.Printf("print using loop const \n")
	}

	samplemodule.SampleFunction()

	fmt.Print(types.Sample{SampleString: "hello from struct", SampleInt: 1})
}
