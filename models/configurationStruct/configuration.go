package Configuration

type Configuration struct {
	Server                    string            `json:"server"`
	Port                      string            `json:"port"`
	User                      string            `json:"user"`
	Pass                      string            `json:"pass"`
	TLS                       bool              `json:"tls"`
	MailboxDurations          []MailboxDuration `json:"durations"`
	Logurl					  string            `json:"logurl"`
	CronScheduleConfig        string            `json:"cronScheduleConfig"`
	Err                       string
}

type MailboxDuration struct {
	DestMailbox   string `json:"destMailbox"`
	SourceMailbox string `json:"sourceMailbox"`
	Duration      int    `json:"duration"`
}

func CreateNewConfiguration() Configuration {
	config := Configuration{}

	return config
}

func (config *Configuration) GetServerURI() string {
	return config.Server + ":" + config.Port
}