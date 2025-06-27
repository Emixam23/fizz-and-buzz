package main

import (
	"gitlab.com/emixam23/fizz-and-buzz/internal"
)

func main() {

	appService, err := internal.New()
	if err != nil {
		panic(err)
	}

	if err := appService.Start(); err != nil {
		panic(err)
	}

}
