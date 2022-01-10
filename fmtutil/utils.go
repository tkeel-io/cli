package fmtutil

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/dapr/cli/utils"
	"github.com/gocarina/gocsv"
	"github.com/olekukonko/tablewriter"
	"github.com/tkeel-io/cli/pkg/print"
)

// PrintTable to print in the table format.
func PrintTable(csvContent string) {
	WriteTable(os.Stdout, csvContent)
}

// WriteTable writes the csv table to writer.
func WriteTable(writer io.Writer, csvContent string) {
	table := tablewriter.NewWriter(writer)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetHeaderLine(false)
	table.SetRowLine(false)
	table.SetCenterSeparator("")
	table.SetRowSeparator("")
	table.SetColumnSeparator("")
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	scanner := bufio.NewScanner(strings.NewReader(csvContent))
	header := true

	for scanner.Scan() {
		text := strings.Split(scanner.Text(), ",")

		if header {
			table.SetHeader(text)
			header = false
		} else {
			table.Append(text)
		}
	}

	table.Render()
}

func Output(list interface{}) {
	table, err := gocsv.MarshalString(list)
	if err != nil {
		print.FailureStatusEvent(os.Stdout, err.Error())
		os.Exit(1)
	}

	utils.PrintTable(table)
}
