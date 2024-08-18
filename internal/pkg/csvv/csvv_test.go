package csvv_test

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/samber/lo"
)

func TestXxx(t *testing.T) {
	records := [][]string{
		{"first_name", "\",last,\",_name", "username"},
		{"Rob", "Pike", "rob"},
		{"Ken", "Thompson", "ken"},
		{"Robert", "Griesemer", "gri"},
	}

	w := csv.NewWriter(os.Stdout)

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n\n==============\n%+v\n==============\n\n", time.Now().Format(time.RFC3339))

	fmt.Printf("\n\n==============\n%+v\n==============\n\n",
		lo.Times(10, func(_ int) []string {
			return lo.Times(3, func(_ int) string { return "a" })
		}),
	)
}
