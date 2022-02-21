package monzo

type Config struct {
	PolitenessInterval uint // Minimum milliseconds between requests to a domain
	ConcurrentFetchers uint
}

func GetConfig() Config {
	return Config{
		PolitenessInterval: 100,
		ConcurrentFetchers: 4,
	}
}
