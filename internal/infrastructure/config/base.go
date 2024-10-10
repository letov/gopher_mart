package config

type Config struct {
	Salt string
}

type preConfig struct {
	Salt *string
}

type setConfig struct {
	Salt bool
}

func newPreConfig() preConfig {
	return preConfig{
		Salt: nil,
	}
}

func newSetConfig() setConfig {
	return setConfig{
		Salt: false,
	}
}
