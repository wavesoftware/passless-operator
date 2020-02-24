package controller

import (
	"github.com/wavesoftware/passless-operator/pkg/controller/passless"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, passless.Add)
}
