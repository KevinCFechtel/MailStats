package main

import (
	"context"
	"strconv"
	"time"

	imapHandler "github.com/KevinCFechtel/MailStats/functions/imapHandler"
	Configuration "github.com/KevinCFechtel/MailStats/models/configurationStruct"
	configHandler "github.com/KevinCFechtel/goConfigHandler"
	"github.com/pterm/pterm"
)

func main() {
	mailbox := "INBOX"
	amount := 10
	options := []string{"List Mailboxes", "Select Mailbox", "Select amount of Mails shown", "List Top Sender", "List Top biggest Mails", "Exit"}
	pterm.Printfln("Please provide the path to the config file:")
	configFilePath, _ := pterm.DefaultInteractiveTextInput.WithDefaultValue("config.json").Show()
	pterm.Println()
	
	configuration := Configuration.CreateNewConfiguration()
	configHandler.GetConfig("localFile", configFilePath, &configuration, "File not found")
	selectedOption := ""
	for selectedOption != "Exit" {
		pterm.Printfln("Selected Mailbox: %s", mailbox)
		pterm.Printfln("Selected Amount: %d", amount)
		selectedOption, _ = pterm.DefaultInteractiveSelect.WithOptions(options).Show()
		switch selectedOption {
			case "List Mailboxes": 
				listMailboxes(configuration)
			case "Select Mailbox": 
				mailbox, _ = pterm.DefaultInteractiveTextInput.Show()
				pterm.Println()
			case "Select amount of Mails shown": 
				amountText, _ := pterm.DefaultInteractiveTextInput.Show()
				amount, _ =  strconv.Atoi(amountText)
				pterm.Println()
			case "List Top Sender": 
				listTopTenSender(configuration, mailbox, amount)
			case "List Top biggest Mails": 
				listTopTenBiggsetMails(configuration, mailbox, amount)
		}
	}
}

func listTopTenSender(configuration Configuration.Configuration, mailbox string, amount int) {
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

	err = imapServer.GetTopTenSenders(mailbox, amount)
	if(err != nil) {
		pterm.Fatal.PrintOnError("Failed to get top senders: " + err.Error(), true)
		return
	}

	imapServer.Logout()
	elapsed := time.Since(start)
	introSpinner.Success("Completed in " + elapsed.Round(time.Second).String())
}

func listTopTenBiggsetMails(configuration Configuration.Configuration, mailbox string, amount int) {
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

	err = imapServer.GetTopTenBiggestMails(mailbox, amount)
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