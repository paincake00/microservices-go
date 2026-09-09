package mongo

type Option func(c *Config)

func WithNetworkName(network string) Option {
	return func(c *Config) {
		c.NetworkName = network
	}
}

func WithContainerName(name string) Option {
	return func(c *Config) {
		c.ContainerName = name
	}
}

func WithImageName(name string) Option {
	return func(c *Config) {
		c.ImageName = name
	}
}

func WithDatabase(db string) Option {
	return func(c *Config) {
		c.Database = db
	}
}

func WithUsername(username string) Option {
	return func(c *Config) {
		c.Username = username
	}
}

func WithPassword(password string) Option {
	return func(c *Config) {
		c.Password = password
	}
}

func WithAuthDB(db string) Option {
	return func(c *Config) {
		c.AuthDB = db
	}
}

func WithLogger(logger Logger) Option {
	return func(c *Config) {
		c.Logger = logger
	}
}
