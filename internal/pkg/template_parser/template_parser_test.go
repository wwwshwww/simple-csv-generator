package template_parser_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/wwwshwww/simple-csv-generator/internal/pkg/template_parser"
)

func TestParse(t *testing.T) {
	f, err := os.Open("example.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	result := template_parser.ParseFromYAML(f)
	fmt.Printf("\n\n==============\n%+v\n==============\n\n", result.MustGet())
}
