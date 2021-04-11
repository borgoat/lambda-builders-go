package lambdabuilders

import (
	"log"
)

func Build(params *Params) error {
	err := GenericCall("LambdaBuilder.build", params, nil)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
