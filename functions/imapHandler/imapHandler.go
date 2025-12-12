package imapHandler

import (
	"context"
	"time"

	Logger "github.com/KevinCFechtel/ImapArchive/models/logger"
	imap "github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type ImapServer struct {
    server  string
	user    string
    pass    string
	tls     bool
	cliente *client.Client
	logger *Logger.Logger
	context context.Context
}


func NewImapServer( server string, user string, pass string, tls bool, logger *Logger.Logger, ctx context.Context) *ImapServer {
    imapServer := &ImapServer{}
	imapServer.server = server
    imapServer.user = user
    imapServer.pass = pass
	imapServer.tls = tls
	imapServer.logger = logger
	imapServer.context = ctx

    return imapServer
}

func (imapServer *ImapServer) Connect() error {
    var (
		imapClient *client.Client
		err error
	)

	if(imapServer.tls) {
		imapClient, err = client.DialTLS(imapServer.server, nil)
	} else {
		imapClient, err = client.Dial(imapServer.server)
	}

	if err != nil {
		imapServer.logger.LogThis("Failed to Dial: " + err.Error(), true)
		return err
	}

	imapServer.cliente = imapClient
	return nil
}

func (imapServer *ImapServer) Login() error {
    if err := imapServer.cliente.Login(imapServer.user, imapServer.pass); err != nil {
		imapServer.logger.LogThis("Failed to Login: " + err.Error(), true)
		return err
    }
	return nil
}

func (imapServer *ImapServer) setLabelBox(label string) (*imap.MailboxStatus, error) {
    mailbox, err := imapServer.cliente.Select(label, false)
    if err != nil {
		imapServer.logger.LogThis("Failed to Select Mailbox: " + err.Error(), true)
		return nil, err
    }
    return mailbox, err
}

func (imapServer *ImapServer) ArchiveMessages(sourceMailbox string, duration int, destMailbox string) (int, error) {
	beforeCriteria := time.Now()
	if duration != 0 {
		daysToPreserve := duration * -1
		beforeCriteria = time.Now().AddDate(0, 0, daysToPreserve)
	}
	
	_, err := imapServer.setLabelBox(sourceMailbox)
	if(err != nil) {
		return 0, err
	}

	criteria := imap.NewSearchCriteria()
	criteria.Before = beforeCriteria

	uids, err := imapServer.cliente.UidSearch(criteria)
	if err != nil {
		imapServer.logger.LogThis("Failed to Search UID: " + err.Error(), true)
		return 0, err
	}
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uids...)
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags, imap.FetchInternalDate, section.FetchItem()}
	messages := make(chan *imap.Message)
			
	if len(uids) > 0 {
		go func() {
			if err := imapServer.cliente.UidFetch(seqSet, items, messages); err != nil {
				imapServer.logger.LogThis("Failed to Fetch UID: " + err.Error(), true)
			}
		}()
		seqSetArchive := new(imap.SeqSet)
		for message := range messages {
			if message == nil {
				imapServer.logger.LogThis("Server didn't returned message", true)
			} else {
				seqSetArchive.AddNum(message.SeqNum)
				//armored, err := helper.SignCleartextMessageArmored(unlockedKeyObj, passphrase, message.Body)
			}
		}

		err = imapServer.cliente.Move(seqSetArchive, destMailbox)
		if err != nil {
			imapServer.logger.LogThis("Failed to Move: " + err.Error(), true)
		}
	}

	return len(uids), nil
}

func (imapServer *ImapServer) Logout() {
    if err := imapServer.cliente.Close(); err != nil {
		imapServer.logger.LogThis("Failed to Close: " + err.Error(), true)
    }
    
    if err := imapServer.cliente.Logout(); err != nil {
		imapServer.logger.LogThis("Failed to Logout: " + err.Error(), true)
    }
}