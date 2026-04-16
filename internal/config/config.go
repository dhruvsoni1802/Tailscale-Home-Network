package config

type Config struct {
	Hostname   string
	StateDir   string
	AuthKey    string
	StorageDir string
}

func Default() Config {
	return Config{
		Hostname:   "storage-node",
		StateDir:   "./tsnet-state",
		StorageDir: "./storage-data",
	}
}