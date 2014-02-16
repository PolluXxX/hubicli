package main

import (
	"os"
)

var cmdInfo = &Command{
	UsageLine: "info",
	Short:     "get information on your hubiC account",
	Long: `
info gives you information on your hubiC account : email, firstname, lastname, ...

`,
}

func init() {
	cmdInfo.Run = runInfo
}

var infoTemplate = `Information about your hubiC account

    Firstname: {{.Firstname}}
    Lastname: {{.Lastname}}
    Email: {{.Email}}
    Creation date: {{.CreationDate}}
    Status: {{.Status}}
    Offer: {{.Offer}}

    Account usage:
        {{.Usage.Used}} used bytes
        {{.Usage.Quota}} total bytes
`

func runInfo(cmd *Command, args []string) {
	if Account.Usage == nil {
		Account.GetUsage()
	}

	tmpl(os.Stdout, infoTemplate, Account)
}
