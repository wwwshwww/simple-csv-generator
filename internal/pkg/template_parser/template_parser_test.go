package template_parser_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/wwwshwww/simple_csv_generator/internal/pkg/template_parser"
)

func TestParse(t *testing.T) {
	f, err := os.Open("example_template.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	result := template_parser.ParseFromYAML(f)
	fmt.Printf("\n\n==============\n%+v\n==============\n\n", result.MustGet())
}
