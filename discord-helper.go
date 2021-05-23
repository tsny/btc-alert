package main

import (
	"strings"

	"github.com/olekukonko/tablewriter"
)

var helpArray = [][]string{
	{"help, commands", "Displays this block"},
	{"get", "Gets summary of the securities price and movements"},
	{"whois", "Gets summary of the stocks underlying company"},
	{"sub", "Subscribes the asking user to the target price"},
	{"chart, graph", "Shows basic graph of the security's price over 2 days"},
	{"quiet", "Toggles whether alerts notify everyone"},
	{"host", "Returns the host name (computer) of the bot"},
}

// Note: could remove SetBorder and SetRowLine if char count needs to be reduced
func formatHelpText() string {
	str := &strings.Builder{}
	table := tablewriter.NewWriter(str)
	table.SetHeader([]string{"Command", "Desc"})
	table.AppendBulk(helpArray)
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetHeaderLine(true)
	table.Render()
	return "```" + str.String() + "```"
}
