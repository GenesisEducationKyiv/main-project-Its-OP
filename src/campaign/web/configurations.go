package web

type FileConfiguration struct {
	EmailStorageFile string
}

type SendgridConfiguration struct {
	ApiKey      string
	SenderName  string
	SenderEmail string
}

type ProviderConfiguration struct {
	Hostname string
	Schema   string
}
