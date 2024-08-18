package cmd

import "testing"

func TestGen(t *testing.T) {
	if err := run("example.yaml", "out.csv", 100); err != nil {
		t.Fatal(err)
	}
}
