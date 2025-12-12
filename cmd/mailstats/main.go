package main

import (
	"context"
	"time"

	imapHandler "github.com/KevinCFechtel/MailStats/functions/imapHandler"
	Configuration "github.com/KevinCFechtel/MailStats/models/configurationStruct"
	configHandler "github.com/KevinCFechtel/goConfigHandler"
	"github.com/pterm/pterm"
)

func main() {
	mailbox := "INBOX"
	options := []string{"List Mailboxes", "Select Mailbox", "List Top 10 Sender", "List Top 10 biggest Mails", "Exit"}
	pterm.Printfln("Please provide the path to the config file:")
	configFilePath, _ := pterm.DefaultInteractiveTextInput.WithDefaultValue("config.json").Show()
	pterm.Println()
	
	configuration := Configuration.CreateNewConfiguration()
	configHandler.GetConfig("localFile", configFilePath, &configuration, "File not found")
	selectedOption := ""
	for selectedOption != "Exit" {
		pterm.Printfln("Selected Mailbox: %s", mailbox)
		selectedOption, _ = pterm.DefaultInteractiveSelect.WithOptions(options).Show()
		switch selectedOption {
			case "List Mailboxes": 
				listMailboxes(configuration)
			case "Select Mailbox": 
				mailbox, _ = pterm.DefaultInteractiveTextInput.Show()
				pterm.Println()
			case "List Top 10 Sender": 
				listTopTenSender(configuration, mailbox)
			case "List Top 10 biggest Mails": 
				listTopTenBiggsetMails(configuration, mailbox)
		}
	}
}

func listTopTenSender(configuration Configuration.Configuration, mailbox string) {
	start := time.Now()
	introSpinner, _ := pterm.DefaultSpinner.WithShowTimer(true).WithRemoveWhenDone(false).Start("Waiting for results for mailbox: " + mailbox + " ...")
	ctx := context.Background()
	imapServer := imapHandler.NewImapServer(configuration.GetServerURI(), configuration.User, configuration.Pass, configuration.TLS, ctx)
	err := imapServer.Connect()
	if(err != nil) {
		pterm.Fatal.PrintOnError("Failed to connect: " + err.Error(), true)
		return
	}

	err = imapServer.Login()
	if(err != nil) {
		pterm.Fatal.PrintOnError("Failed to login: " + err.Error(), true)
		return
	}

	err = imapServer.GetTopTenSenders(mailbox)
	if(err != nil) {
		pterm.Fatal.PrintOnError("Failed to get top senders: " + err.Error(), true)
		return
	}

	imapServer.Logout()
	elapsed := time.Since(start)
	introSpinner.Success("Completed in " + elapsed.Round(time.Second).String())
}

func listTopTenBiggsetMails(configuration Configuration.Configuration, mailbox string) {
	start := time.Now()
	introSpinner, _ := pterm.DefaultSpinner.WithShowTimer(true).WithRemoveWhenDone(false).Start("Waiting for results for mailbox: " + mailbox + " ...")
	ctx := context.Background()
	imapServer := imapHandler.NewImapServer(configuration.GetServerURI(), configuration.User, configuration.Pass, configuration.TLS, ctx)
	err := imapServer.Connect()
	if(err != nil) {
		pterm.Fatal.PrintOnError("Failed to connect: " + err.Error(), true)
		return
	}

	err = imapServer.Login()
	if(err != nil) {
		pterm.Fatal.PrintOnError("Failed to login: " + err.Error(), true)
		return
	}

	err = imapServer.GetTopTenBiggestMails(mailbox)
	if(err != nil) {
		pterm.Fatal.PrintOnError("Failed to get top senders: " + err.Error(), true)
		return
	}

	imapServer.Logout()
	elapsed := time.Since(start)
	introSpinner.Success("Completed in " + elapsed.Round(time.Second).String())
}

func listMailboxes(configuration Configuration.Configuration) {
	ctx := context.Background()
	imapServer := imapHandler.NewImapServer(configuration.GetServerURI(), configuration.User, configuration.Pass, configuration.TLS, ctx)
	err := imapServer.Connect()
	if(err != nil) {
		pterm.Fatal.PrintOnError("Failed to connect: " + err.Error(), true)
		return
	}

	err = imapServer.Login()
	if(err != nil) {
		pterm.Fatal.PrintOnError("Failed to login: " + err.Error(), true)
		return
	}

	mailboxes, err := imapServer.GetMailboxes()
	if(err != nil) {
		pterm.Fatal.PrintOnError("Failed to list mailboxes: " + err.Error(), true)
		return
	}

	pterm.Printfln("Mailboxes:")
	for _, mailbox := range mailboxes {
		pterm.Printfln("- %s", mailbox)
	}

	imapServer.Logout()
}