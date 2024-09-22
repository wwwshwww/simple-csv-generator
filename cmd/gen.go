/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/spf13/cobra"
	"github.com/wwwshwww/simple-csv-generator/internal/pkg/dummy_producer"
	"github.com/wwwshwww/simple-csv-generator/internal/pkg/template_parser"
)

const (
	dummyIntSpecies             = 10
	dummyFloatSpecies           = 10
	dummyStringSpecies          = 10
	dummyBoolSpecies            = 10
	dummyDatetimeSpecies        = 10
	dummyMultilineStringSpecies = 10
	dummyURLSpecies             = 10

	dummyIntMin                = 0
	dummyIntMax                = 1000
	dummyFloatMin              = 0.0
	dummyFloatMax              = 1000.0
	dummyStringLength          = 6
	dummyMultilineStringLength = 6
	dummyMultilineStringLines  = 3
	dummyURLLength             = 8

	dummyArrayIntElements             = 3
	dummyArrayIntSpecies              = 3
	dummyArrayFloatElements           = 3
	dummyArrayFloatSpecies            = 3
	dummyArrayStringElements          = 3
	dummyArrayStringSpecies           = 3
	dummyArrayBoolElements            = 3
	dummyArrayBoolSpecies             = 3
	dummyArrayDatetimeElements        = 3
	dummyArrayDatetimeSpecies         = 3
	dummyArrayMultilineStringElements = 3
	dummyArrayMultilineStringSpecies  = 3
	dummyArrayURLElements             = 3
	dummyArrayURLSpecies              = 3
)

var (
	dummyDatetimeStart = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	dummyDatetimeEnd   = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := run(templateFilePath, outputPath, rowCount); err != nil {
			fmt.Printf("failed to run: %v\n", err)
		}
	},
}

var templateFilePath string
var outputPath string
var rowCount int

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.Flags().StringVarP(&templateFilePath, "template", "t", "", "Template file path")
	genCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file path")
	genCmd.Flags().IntVarP(&rowCount, "row-count", "r", 10, "Row count to generate")

	if err := genCmd.MarkFlagRequired("template"); err != nil {
		log.Fatalf("failed to mark flag required: %v", err)
	}
	genCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Lookup("output").Value.String() == "" {
			t := time.Now().Format(time.RFC3339)
			if err := cmd.Flags().Set("output", fmt.Sprintf("output_%v.csv", t)); err != nil {
				log.Fatalf("failed to set output flag: %v", err)
			}
		}
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func run(templateFilePath, outputPath string, rowCount int) error {
	source, err := os.Open(templateFilePath)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer source.Close()

	content, err := generateCsvContent(source, rowCount).Get()
	if err != nil {
		log.Fatalf("failed to generate csv content: %v", err)
	}

	output, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer output.Close()

	w := csv.NewWriter(output)
	defer w.Flush()

	err = w.Write(content[0])
	if err != nil {
		log.Fatalf("failed to write header: %v", err)
	}

	for _, record := range content[1:] {
		err = w.Write(record)
		if err != nil {
			log.Fatalf("failed to write record: %v", err)
		}
	}

	fmt.Printf("Done 🎉\nGenerated: %v\n", outputPath)
	return nil
}

func generateCsvContent(source io.Reader, rowCount int) mo.Result[[][]string] {
	columns, err := template_parser.ParseFromYAML(source).Get()
	if err != nil {
		mo.Errf[[][]string]("failed to get columns: %e", err)
	}
	header := lo.Map(columns, func(c template_parser.Column, _ int) string {
		return c.Name
	})
	rows := lo.Times(rowCount, func(_ int) []string {
		return lo.Times(len(columns), func(_ int) string { return "" })
	})

	// * prepare dummies and choices
	nameToDummiesInt := make(map[string][]int)
	nameToDummiesFloat := make(map[string][]float64)
	nameToDummiesBool := make(map[string][]bool)
	nameToDummiesDatetime := make(map[string][]time.Time)
	nameToDummiesString := make(map[string][]string)
	nameToDummiesMultilineString := make(map[string][]string)
	nameToDummiesURL := make(map[string][]string)
	nameToDummiesArrayInt := make(map[string][][]int)
	nameToDummiesArrayFloat := make(map[string][][]float64)
	nameToDummiesArrayBool := make(map[string][][]bool)
	nameToDummiesArrayDatetime := make(map[string][][]time.Time)
	nameToDummiesArrayString := make(map[string][][]string)
	nameToDummiesArrayMultilineString := make(map[string][][]string)
	nameToDummiesArrayURL := make(map[string][][]string)

	nameToChoicesInt := make(map[string][]int)
	nameToChoicesFloat := make(map[string][]float64)
	nameToChoicesBool := make(map[string][]bool)
	nameToChoicesDatetime := make(map[string][]time.Time)
	nameToChoicesString := make(map[string][]string)
	nameToChoicesMultilineString := make(map[string][]string)
	nameToChoicesURL := make(map[string][]string)
	nameToChoicesArrayInt := make(map[string][][]int)
	nameToChoicesArrayFloat := make(map[string][][]float64)
	nameToChoicesArrayBool := make(map[string][][]bool)
	nameToChoicesArrayDatetime := make(map[string][][]time.Time)
	nameToChoicesArrayString := make(map[string][][]string)
	nameToChoicesArrayMultilineString := make(map[string][][]string)
	nameToChoicesArrayURL := make(map[string][][]string)

	for _, c := range columns {
		switch c.Type {
		case template_parser.ColumnTypeInt:
			nameToDummiesInt[c.Name] = dummy_producer.GetDummiesInt(dummyIntSpecies, dummyIntMin, dummyIntMax)
			if v := c.Choices.Values.LeftOrEmpty(); len(v) > 0 {
				r, err := template_parser.UnmarshalIntChoices(v).Get()
				if err != nil {
					return mo.Errf[[][]string]("failed to get choices: %e", err)
				}
				nameToChoicesInt[c.Name] = r
			} else if vs := c.Choices.Values.RightOrEmpty(); len(vs) > 0 {
				fmt.Printf("column [%v]: invalid format for choices.\n", c.Name)
			}
		case template_parser.ColumnTypeFloat:
			nameToDummiesFloat[c.Name] = dummy_producer.GetDummiesFloat(dummyFloatSpecies, dummyFloatMin, dummyFloatMax)
			if v := c.Choices.Values.LeftOrEmpty(); len(v) > 0 {
				r, err := template_parser.UnmarshalFloatChoices(v).Get()
				if err != nil {
					return mo.Errf[[][]string]("failed to get choices: %e", err)
				}
				nameToChoicesFloat[c.Name] = r
			} else if vs := c.Choices.Values.RightOrEmpty(); len(vs) > 0 {
				fmt.Printf("column [%v]: invalid format for choices.\n", c.Name)
			}
		case template_parser.ColumnTypeBool:
			nameToDummiesBool[c.Name] = dummy_producer.GetDummiesBool(dummyBoolSpecies)
			if v := c.Choices.Values.LeftOrEmpty(); len(v) > 0 {
				r, err := template_parser.UnmarshalBoolChoices(v).Get()
				if err != nil {
					return mo.Errf[[][]string]("failed to get choices: %e", err)
				}
				nameToChoicesBool[c.Name] = r
			} else if vs := c.Choices.Values.RightOrEmpty(); len(vs) > 0 {
				fmt.Printf("column [%v]: invalid format for choices.\n", c.Name)
			}
		case template_parser.ColumnTypeDatetime:
			nameToDummiesDatetime[c.Name] = dummy_producer.GetDummiesDatetime(dummyDatetimeSpecies, dummyDatetimeStart, dummyDatetimeEnd)
			if v := c.Choices.Values.LeftOrEmpty(); len(v) > 0 {
				r, err := template_parser.UnmarshalDatetimeChoices(v).Get()
				if err != nil {
					return mo.Errf[[][]string]("failed to get choices: %e", err)
				}
				nameToChoicesDatetime[c.Name] = r
			} else if vs := c.Choices.Values.RightOrEmpty(); len(vs) > 0 {
				fmt.Printf("column [%v]: invalid format for choices.\n", c.Name)
			}
		case template_parser.ColumnTypeString:
			nameToDummiesString[c.Name] = dummy_producer.GetDummiesString(dummyStringSpecies, dummyStringLength)
			if v := c.Choices.Values.LeftOrEmpty(); len(v) > 0 {
				nameToChoicesString[c.Name] = v
			} else if vs := c.Choices.Values.RightOrEmpty(); len(vs) > 0 {
				fmt.Printf("column [%v]: invalid format for choices.\n", c.Name)
			}
		case template_parser.ColumnTypeMultilineString:
			nameToDummiesMultilineString[c.Name] = dummy_producer.GetDummiesMultilineString(dummyMultilineStringSpecies, dummyStringLength, dummyMultilineStringLines)
			if v := c.Choices.Values.LeftOrEmpty(); len(v) > 0 {
				nameToChoicesMultilineString[c.Name] = v
			} else if vs := c.Choices.Values.RightOrEmpty(); len(vs) > 0 {
				fmt.Printf("column [%v]: invalid format for choices.\n", c.Name)
			}
		case template_parser.ColumnTypeURL:
			nameToDummiesURL[c.Name] = dummy_producer.GetDummiesURL(dummyURLSpecies, dummyURLLength)
			if v := c.Choices.Values.LeftOrEmpty(); len(v) > 0 {
				nameToChoicesURL[c.Name] = v
			} else if vs := c.Choices.Values.RightOrEmpty(); len(vs) > 0 {
				fmt.Printf("column [%v]: invalid format for choices.\n", c.Name)
			}
		case template_parser.ColumnTypeArrayInt:
			nameToDummiesArrayInt[c.Name] = lo.Times(dummyArrayIntSpecies, func(_ int) []int {
				return dummy_producer.GetDummiesInt(dummyArrayIntElements, dummyIntMin, dummyIntMax)
			})
			if vs := c.Choices.Values.RightOrEmpty(); len(vs) > 0 {
				rs := lo.Map(vs, func(v []string, _ int) mo.Result[[]int] {
					return template_parser.UnmarshalIntChoices(v)
				})
				choices := make([][]int, 0, len(rs))
				errs := make([]error, 0, len(rs))
				for _, r := range rs {
					if v, err := r.Get(); err != nil {
						errs = append(errs, err)
					} else {
						choices = append(choices, v)
					}
				}
				if len(errs) > 0 {
					return mo.Errf[[][]string]("failed to get choices: %e", errs)
				}
				nameToChoicesArrayInt[c.Name] = choices
			} else if v := c.Choices.Values.LeftOrEmpty(); len(v) > 0 {
				fmt.Printf("column [%v]: invalid format for choices.\n", c.Name)
			}
		case template_parser.ColumnTypeArrayFloat:
			nameToDummiesArrayFloat[c.Name] = lo.Times(dummyArrayFloatSpecies, func(_ int) []float64 {
				return dummy_producer.GetDummiesFloat(dummyArrayFloatElements, dummyFloatMin, dummyFloatMax)
			})
			if vs := c.Choices.Values.RightOrEmpty(); len(vs) > 0 {
				rs := lo.Map(vs, func(v []string, _ int) mo.Result[[]float64] {
					return template_parser.UnmarshalFloatChoices(v)
				})
				choices := make([][]float64, 0, len(rs))
				errs := make([]error, 0, len(rs))
				for _, r := range rs {
					if v, err := r.Get(); err != nil {
						errs = append(errs, err)
					} else {
						choices = append(choices, v)
					}
				}
				if len(errs) > 0 {
					return mo.Errf[[][]string]("failed to get choices: %e", errs)
				}
				nameToChoicesArrayFloat[c.Name] = choices
			} else if v := c.Choices.Values.LeftOrEmpty(); len(v) > 0 {
				fmt.Printf("column [%v]: invalid format for choices.\n", c.Name)
			}
		case template_parser.ColumnTypeArrayBool:
			nameToDummiesArrayBool[c.Name] = lo.Times(dummyArrayBoolSpecies, func(_ int) []bool {
				return dummy_producer.GetDummiesBool(dummyArrayBoolElements)
			})
			if vs := c.Choices.Values.RightOrEmpty(); len(vs) > 0 {
				rs := lo.Map(vs, func(v []string, _ int) mo.Result[[]bool] {
					return template_parser.UnmarshalBoolChoices(v)
				})
				choices := make([][]bool, 0, len(rs))
				errs := make([]error, 0, len(rs))
				for _, r := range rs {
					if v, err := r.Get(); err != nil {
						errs = append(errs, err)
					} else {
						choices = append(choices, v)
					}
				}
				if len(errs) > 0 {
					return mo.Errf[[][]string]("failed to get choices: %e", errs)
				}
				nameToChoicesArrayBool[c.Name] = choices
			} else if v := c.Choices.Values.LeftOrEmpty(); len(v) > 0 {
				fmt.Printf("column [%v]: invalid format for choices.\n", c.Name)
			}
		case template_parser.ColumnTypeArrayDatetime:
			nameToDummiesArrayDatetime[c.Name] = lo.Times(dummyArrayDatetimeSpecies, func(_ int) []time.Time {
				return dummy_producer.GetDummiesDatetime(dummyArrayDatetimeElements, dummyDatetimeStart, dummyDatetimeEnd)
			})
			if vs := c.Choices.Values.RightOrEmpty(); len(vs) > 0 {
				rs := lo.Map(vs, func(v []string, _ int) mo.Result[[]time.Time] {
					return template_parser.UnmarshalDatetimeChoices(v)
				})
				choices := make([][]time.Time, 0, len(rs))
				errs := make([]error, 0, len(rs))
				for _, r := range rs {
					if v, err := r.Get(); err != nil {
						errs = append(errs, err)
					} else {
						choices = append(choices, v)
					}
				}
				if len(errs) > 0 {
					return mo.Errf[[][]string]("failed to get choices: %e", errs)
				}
				nameToChoicesArrayDatetime[c.Name] = choices
			} else if v := c.Choices.Values.LeftOrEmpty(); len(v) > 0 {
				fmt.Printf("column [%v]: invalid format for choices.\n", c.Name)
			}
		case template_parser.ColumnTypeArrayString:
			nameToDummiesArrayString[c.Name] = lo.Times(dummyArrayStringSpecies, func(_ int) []string {
				return dummy_producer.GetDummiesString(dummyArrayStringElements, dummyStringLength)
			})
			if vs := c.Choices.Values.RightOrEmpty(); len(vs) > 0 {
				nameToChoicesArrayString[c.Name] = vs
			} else if v := c.Choices.Values.LeftOrEmpty(); len(v) > 0 {
				fmt.Printf("column [%v]: invalid format for choices.\n", c.Name)
			}
		case template_parser.ColumnTypeArrayMultilineString:
			nameToDummiesArrayMultilineString[c.Name] = lo.Times(dummyArrayMultilineStringSpecies, func(_ int) []string {
				return dummy_producer.GetDummiesMultilineString(dummyArrayMultilineStringElements, dummyStringLength, dummyMultilineStringLines)
			})
			if vs := c.Choices.Values.RightOrEmpty(); len(vs) > 0 {
				nameToChoicesArrayMultilineString[c.Name] = vs
			} else if v := c.Choices.Values.LeftOrEmpty(); len(v) > 0 {
				fmt.Printf("column [%v]: invalid format for choices.\n", c.Name)
			}
		case template_parser.ColumnTypeArrayURL:
			nameToDummiesArrayURL[c.Name] = lo.Times(dummyArrayURLSpecies, func(_ int) []string {
				return dummy_producer.GetDummiesURL(dummyArrayURLElements, dummyURLLength)
			})
			if vs := c.Choices.Values.RightOrEmpty(); len(vs) > 0 {
				nameToChoicesArrayURL[c.Name] = vs
			} else if v := c.Choices.Values.LeftOrEmpty(); len(v) > 0 {
				fmt.Printf("column [%v]: invalid format for choices.\n", c.Name)
			}
		}
	}

	// * select values from dummies or choices
	for rowIdx := range rowCount {
		for colIdx, c := range columns {

			switch c.Type {
			case template_parser.ColumnTypeInt:
				var v int
				if choices, ok := nameToChoicesInt[c.Name]; ok {
					v = dummy_producer.Select(choices)
				} else {
					v = dummy_producer.Select(nameToDummiesInt[c.Name])
				}
				rows[rowIdx][colIdx] = strconv.Itoa(v)
			case template_parser.ColumnTypeFloat:
				var v float64
				if choices, ok := nameToChoicesFloat[c.Name]; ok {
					v = dummy_producer.Select(choices)
				} else {
					v = dummy_producer.Select(nameToDummiesFloat[c.Name])
				}
				rows[rowIdx][colIdx] = strconv.FormatFloat(v, 'f', -1, 64)
			case template_parser.ColumnTypeBool:
				var v bool
				if choices, ok := nameToChoicesBool[c.Name]; ok {
					v = dummy_producer.Select(choices)
				} else {
					v = dummy_producer.Select(nameToDummiesBool[c.Name])
				}
				rows[rowIdx][colIdx] = strconv.FormatBool(v)
			case template_parser.ColumnTypeDatetime:
				var v time.Time
				if choices, ok := nameToChoicesDatetime[c.Name]; ok {
					v = dummy_producer.Select(choices)
				} else {
					v = dummy_producer.Select(nameToDummiesDatetime[c.Name])
				}
				rows[rowIdx][colIdx] = v.Format(time.RFC3339)
			case template_parser.ColumnTypeString:
				var v string
				if choices, ok := nameToChoicesString[c.Name]; ok {
					v = dummy_producer.Select(choices)
				} else {
					v = dummy_producer.Select(nameToDummiesString[c.Name])
				}
				rows[rowIdx][colIdx] = v
			case template_parser.ColumnTypeMultilineString:
				var v string
				if choices, ok := nameToChoicesMultilineString[c.Name]; ok {
					v = dummy_producer.Select(choices)
				} else {
					v = dummy_producer.Select(nameToDummiesMultilineString[c.Name])
				}
				rows[rowIdx][colIdx] = v
			case template_parser.ColumnTypeURL:
				var v string
				if choices, ok := nameToChoicesURL[c.Name]; ok {
					v = dummy_producer.Select(choices)
				} else {
					v = dummy_producer.Select(nameToDummiesURL[c.Name])
				}
				rows[rowIdx][colIdx] = v
			case template_parser.ColumnTypeArrayInt:
				var v []int
				if choices, ok := nameToChoicesArrayInt[c.Name]; ok {
					v = dummy_producer.Select(choices)
				} else {
					v = dummy_producer.Select(nameToDummiesArrayInt[c.Name])
				}
				rows[rowIdx][colIdx] = strings.Join(lo.Map(v, func(e int, _ int) string { return strconv.Itoa(e) }), ",")
			case template_parser.ColumnTypeArrayFloat:
				var v []float64
				if choices, ok := nameToChoicesArrayFloat[c.Name]; ok {
					v = dummy_producer.Select(choices)
				} else {
					v = dummy_producer.Select(nameToDummiesArrayFloat[c.Name])
				}
				rows[rowIdx][colIdx] = strings.Join(lo.Map(v, func(e float64, _ int) string { return strconv.FormatFloat(e, 'f', -1, 64) }), ",")
			case template_parser.ColumnTypeArrayBool:
				var v []bool
				if choices, ok := nameToChoicesArrayBool[c.Name]; ok {
					v = dummy_producer.Select(choices)
				} else {
					v = dummy_producer.Select(nameToDummiesArrayBool[c.Name])
				}
				rows[rowIdx][colIdx] = strings.Join(lo.Map(v, func(e bool, _ int) string { return strconv.FormatBool(e) }), ",")
			case template_parser.ColumnTypeArrayDatetime:
				var v []time.Time
				if choices, ok := nameToChoicesArrayDatetime[c.Name]; ok {
					v = dummy_producer.Select(choices)
				} else {
					v = dummy_producer.Select(nameToDummiesArrayDatetime[c.Name])
				}
				rows[rowIdx][colIdx] = strings.Join(lo.Map(v, func(e time.Time, _ int) string { return e.Format(time.RFC3339) }), ",")
			case template_parser.ColumnTypeArrayString:
				var v []string
				if choices, ok := nameToChoicesArrayString[c.Name]; ok {
					v = dummy_producer.Select(choices)
				} else {
					v = dummy_producer.Select(nameToDummiesArrayString[c.Name])
				}
				rows[rowIdx][colIdx] = strings.Join(lo.Map(v, func(e string, _ int) string { return nestedCsvEscaper(e) }), ",")
			case template_parser.ColumnTypeArrayMultilineString:
				var v []string
				if choices, ok := nameToChoicesArrayMultilineString[c.Name]; ok {
					v = dummy_producer.Select(choices)
				} else {
					v = dummy_producer.Select(nameToDummiesArrayMultilineString[c.Name])
				}
				rows[rowIdx][colIdx] = strings.Join(lo.Map(v, func(e string, _ int) string { return nestedCsvEscaper(e) }), ",")
			case template_parser.ColumnTypeArrayURL:
				var v []string
				if choices, ok := nameToChoicesArrayURL[c.Name]; ok {
					v = dummy_producer.Select(choices)
				} else {
					v = dummy_producer.Select(nameToDummiesArrayURL[c.Name])
				}
				rows[rowIdx][colIdx] = strings.Join(lo.Map(v, func(e string, _ int) string { return nestedCsvEscaper(e) }), ",")
			}
		}
	}

	// * mask values according to creation probability
	for colIdx, c := range columns {
		for rowIdx := range rowCount {
			if dummy_producer.Selector.Float64() > c.CreationProbability {
				rows[rowIdx][colIdx] = "" // mask
			}
		}
	}

	return mo.Ok(append([][]string{header}, rows...))
}

func nestedCsvEscaper(raw string) string {
	var result string
	if strings.Contains(raw, ",") || strings.Contains(raw, "\"") {
		escaped := []rune{'"'}
		for _, s := range raw {
			if s == '"' {
				escaped = append(escaped, '"')
			}
			escaped = append(escaped, s)
		}
		escaped = append(escaped, '"')
		result = string(escaped)
	} else {
		result = raw
	}
	return result
}
