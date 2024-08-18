package template_parser_test

import (
	"fmt"
	"testing"

	"github.com/wwwshwww/simple_csv_generator/internal/pkg/template_parser"
)

func TestParse(t *testing.T) {
	result := template_parser.ParseColumns("example_template.yaml")
	fmt.Printf("\n\n==============\n%+v\n==============\n\n", result.MustGet())
}
