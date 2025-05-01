package cmd

import (
	"fmt"
	"testing"
)

func TestIterate(t *testing.T) {
	t.Run("test errors", func(t *testing.T) {

		a := fmt.Errorf("deepest error a")
		b := fmt.Errorf("middle b that wraps a: %w", a)
		c := fmt.Errorf("middle c that wraps b: %w", b)
		d := fmt.Errorf("top level d that wraps c: %w", c)

		fmt.Println("--------")
		prettyPrintError(d)
		fmt.Println("--------")
	})
}
