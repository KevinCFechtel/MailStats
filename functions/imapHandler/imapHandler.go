package imapHandler

import (
	"context"
	"time"

	imap "github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type ImapServer struct {
    server  string
	user    string
    pass    string
	tls     bool
	cliente *client.Client
	context context.Context
}


func NewImapServer( server string, user string, pass string, tls bool, ctx context.Context) *ImapServer {
    imapServer := &ImapServer{}
	imapServer.server = server
    imapServer.user = user
    imapServer.pass = pass
	imapServer.tls = tls
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
		return err
	}

	imapServer.cliente = imapClient
	return nil
}

func (imapServer *ImapServer) Login() error {
    if err := imapServer.cliente.Login(imapServer.user, imapServer.pass); err != nil {
		return err
    }
	return nil
}

func (imapServer *ImapServer) setLabelBox(label string) (*imap.MailboxStatus, error) {
    mailbox, err := imapServer.cliente.Select(label, false)
    if err != nil {
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
			}
		}()
		seqSetArchive := new(imap.SeqSet)
		for message := range messages {
			if message == nil {
			} else {
				seqSetArchive.AddNum(message.SeqNum)
				//armored, err := helper.SignCleartextMessageArmored(unlockedKeyObj, passphrase, message.Body)
			}
		}

		err = imapServer.cliente.Move(seqSetArchive, destMailbox)
		if err != nil {
		}
	}

	return len(uids), nil
}

func (imapServer *ImapServer) GetMailboxes() ([]string, error) {
	mailboxes := []string{}
	mailboxChan := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- imapServer.cliente.List("", "*", mailboxChan)
	}()

	for m := range mailboxChan {
		mailboxes = append(mailboxes, m.Name)
	}

	if err := <-done; err != nil {
		return nil, err
	}
	return mailboxes, nil
}

func (imapServer *ImapServer) SelectMailbox(mailbox string) (*imap.MailboxStatus, error) {
	selectedMailbox, err := imapServer.cliente.Select(mailbox, false)
	if err != nil {
		return nil, err
	}
	return selectedMailbox, nil
}

func (imapServer *ImapServer) GetTopTenSenders(mailbox string) (map[string]int, error) {
	// Get the top ten senders from the specified mailbox
	

	return nil, nil
}

func (imapServer *ImapServer) GetTopTenBiggestMails(mailbox string) (map[string]int, error) {

	return nil, nil
}

func (imapServer *ImapServer) Logout() error {
    if err := imapServer.cliente.Close(); err != nil {
		return err
    }
    
    if err := imapServer.cliente.Logout(); err != nil {
		return err
    }
	return nil;
}