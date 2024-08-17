package lib_test

import (
	"fmt"
	"testing"

	"github.com/wwwshwww/simple_csv_generator/lib"
)

func TestParse(t *testing.T) {
	fmt.Printf("\n\n==============\n%+v\n==============\n\n", lib.ParseColumns("template.yaml").MustGet())
}
