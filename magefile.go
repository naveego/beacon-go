// +build mage

package main

import (
	"github.com/magefile/mage/sh"
)

// Generate regenerates the swagger-based code in this package.
func Generate() error {

	if err := sh.Run("autorest", "README.md"); err != nil {
		return err
	}

	sh.Run("go", "fmt")

	return nil
}
