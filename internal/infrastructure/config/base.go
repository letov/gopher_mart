package config

type Config struct {
	Salt   string
	JwtKey string
}

type preConfig struct {
	Salt   *string
	JwtKey *string
}

type setConfig struct {
	Salt   bool
	JwtKey bool
}

func newPreConfig() preConfig {
	return preConfig{
		Salt:   nil,
		JwtKey: nil,
	}
}

func newSetConfig() setConfig {
	return setConfig{
		Salt:   false,
		JwtKey: false,
	}
}
