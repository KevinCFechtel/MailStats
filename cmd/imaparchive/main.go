package main

import (
	"context"
	"flag"
	"fmt"

	imapHandler "github.com/KevinCFechtel/ImapArchive/functions/imapHandler"
	Configuration "github.com/KevinCFechtel/ImapArchive/models/configurationStruct"
	Logger "github.com/KevinCFechtel/ImapArchive/models/logger"
	configHandler "github.com/KevinCFechtel/goConfigHandler"
	"github.com/go-co-op/gocron/v2"
)

func main() {
	var configFilePath string
	var runAsService bool
	flag.StringVar(&configFilePath, "configFile", "config.json", "File Path to config File!")
	flag.BoolVar(&runAsService, "runAsService", false, "Run imaparchive as a service")
	flag.Parse()
	
	ctx := context.Background()

	configuration := Configuration.CreateNewConfiguration()
	configHandler.GetConfig("localFile", configFilePath, &configuration, "File not found")
	
	logger := Logger.NewLogger(configuration.Logurl)
	
	if(runAsService) {
		s, err := gocron.NewScheduler()
		if err != nil {
			logger.LogThis("Failed to start scheduler: " + err.Error(), true)
		}
		defer func() { _ = s.Shutdown() }()
		_, err = s.NewJob(
			gocron.CronJob(
				configuration.CronScheduleConfig,
				false,
			),
			gocron.NewTask(
				func() {
					arichvedMails, err := archiveMails(configuration, logger, ctx)
					if(err != nil) {
						logger.LogThis("Run unsuccesfully, moved " + fmt.Sprint(arichvedMails) + " Messages", false)
					} else {
						logger.LogThis("Run succesfully, moved " + fmt.Sprint(arichvedMails) + " Messages", false)
					}
				},
			),
		)
		if err != nil {
			logger.LogThis("Failed to config scheduler job: " + err.Error(), true)
		}
		s.Start()
		select {} // wait forever
	} else {
		arichvedMails, err := archiveMails(configuration, logger, ctx)
		if(err != nil) {
			logger.LogThis("Run unsuccesfully, moved " + fmt.Sprint(arichvedMails) + " Messages", false)
		} else {
			logger.LogThis("Run succesfully, moved " + fmt.Sprint(arichvedMails) + " Messages", false)
		}
	}
}

func archiveMails(configuration Configuration.Configuration, logger *Logger.Logger, ctx context.Context) (int,error) {
	countArchiveMessages := 0
	imapServer := imapHandler.NewImapServer(configuration.GetServerURI(), configuration.User, configuration.Pass, configuration.TLS, logger, ctx)
	err := imapServer.Connect()
	if(err != nil) {
		return 0, err
	}

	err = imapServer.Login()
	if(err != nil) {
		return 0, err
	}
	
	for _, s := range configuration.MailboxDurations {
		arichvedMessages, _ := imapServer.ArchiveMessages(s.SourceMailbox, s.Duration, s.DestMailbox)
		countArchiveMessages = countArchiveMessages + arichvedMessages
	}

	imapServer.Logout()
	return countArchiveMessages, nil
}