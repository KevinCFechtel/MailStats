package imapHandler

import (
	"context"
	"sort"

	imap "github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/pterm/pterm"
)

type ImapServer struct {
    server  string
	user    string
    pass    string
	tls     bool
	cliente *client.Client
	context context.Context
}

type entryCount  struct {
    val int
    key string
}

type entrySize  struct {
    val uint32
    key string
}

type entriesSize []entrySize
func (s entriesSize) Len() int { return len(s) }
func (s entriesSize) Less(i, j int) bool { return s[i].val < s[j].val }
func (s entriesSize) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type entriesCount []entryCount
func (s entriesCount) Len() int { return len(s) }
func (s entriesCount) Less(i, j int) bool { return s[i].val < s[j].val }
func (s entriesCount) Swap(i, j int) { s[i], s[j] = s[j], s[i] }


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

func (imapServer *ImapServer) GetTopTenSenders(mailbox string) (error) {
	// Get the top ten senders from the specified mailbox
	senders := make(map[string]int)
	_, err := imapServer.setLabelBox(mailbox)
	if(err != nil) {
		return err
	}

	criteria := imap.NewSearchCriteria()
	//criteria.WithoutFlags = []string{"\\Deleted"}

	uids, err := imapServer.cliente.UidSearch(criteria)
	if err != nil {
		return err
	}
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uids...)
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags, imap.FetchInternalDate, section.FetchItem()}
	messages := make(chan *imap.Message)
			
	if len(uids) > 0 {
		go func() {
			if err := imapServer.cliente.UidFetch(seqSet, items, messages); err != nil {
				pterm.Fatal.PrintOnError("Failed to get top senders: " + err.Error(), true)
			}
		}()
		for message := range messages {
			if message != nil {
				if len(message.Envelope.Sender) > 0 {
					senders[message.Envelope.Sender[0].Address()]++
				}
			}
		}
	}

	var es entriesCount
    for k, v := range senders {
        es = append(es, entryCount{val: v, key: k})
    }

    sort.Sort(sort.Reverse(es))

	pterm.Printfln("Top 10 Senders in mailbox %s:", mailbox)
    for count, e := range es {
		if count < 10 {
			pterm.Printfln("%s: %d mails", e.key, e.val)
		}
    }

	return nil
}

func (imapServer *ImapServer) GetTopTenBiggestMails(mailbox string) error {
	// Get the top ten senders from the specified mailbox
	messageSizes := make(map[string]uint32)
	_, err := imapServer.setLabelBox(mailbox)
	if(err != nil) {
		return err
	}

	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{"\\Deleted"}

	uids, err := imapServer.cliente.UidSearch(criteria)
	if err != nil {
		return err
	}
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uids...)
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags, imap.FetchInternalDate, section.FetchItem(), imap.FetchRFC822Size}
	messages := make(chan *imap.Message)
			
	if len(uids) > 0 {
		go func() {
			if err := imapServer.cliente.UidFetch(seqSet, items, messages); err != nil {
			}
		}()
		for message := range messages {
			if message != nil {
				messageSizes[message.Envelope.Subject] = message.Size
			}
		}
	}

	var es entriesSize
    for k, v := range messageSizes {
        es = append(es, entrySize{val: v, key: k})
    }

    sort.Sort(sort.Reverse(es))

	pterm.Printfln("Top 10 biggest Mails in mailbox %s:", mailbox)
    for count, e := range es {
		if count < 10 {
			pterm.Printfln("Subject: %s: %d MB size", e.key, e.val / (1024 * 1024))
		}
    }

	return nil
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