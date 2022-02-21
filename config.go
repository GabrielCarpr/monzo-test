package monzo

type Config struct {
	PolitenessInterval uint // Minimum milliseconds between requests to a domain
}

func GetConfig() Config {
	return Config{
		PolitenessInterval: 100,
	}
}
