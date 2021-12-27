package burn

import (
	"os"
)

const (
	Dev  string = "development"
	Prod string = "production"
	Test string = "test"
)

var Env = Dev
var Root string

func setENV(e string) {
	if len(e) > 0 {
		Env = e
	}
}

func init() {
	setENV(os.Getenv("BURN_ENV"))
	var err error
	Root, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}
