package errors

import (
	"fmt"
	"os"
)

func Handle(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}
