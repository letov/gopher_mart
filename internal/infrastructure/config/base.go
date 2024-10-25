package config

type Config struct {
	Salt         string
	JwtKey       string
	Ampq         string
	AccrualUrl   string
	DBConnection string
	Addr         string
}

type preConfig struct {
	Salt         *string
	JwtKey       *string
	Ampq         *string
	AccrualUrl   *string
	DBConnection *string
	Addr         *string
}

type setConfig struct {
	Salt         bool
	JwtKey       bool
	Ampq         bool
	AccrualUrl   bool
	DBConnection bool
	Addr         bool
}

func newPreConfig() preConfig {
	return preConfig{
		Salt:         nil,
		JwtKey:       nil,
		Ampq:         nil,
		AccrualUrl:   nil,
		DBConnection: nil,
		Addr:         nil,
	}
}

func newSetConfig() setConfig {
	return setConfig{
		Salt:         false,
		JwtKey:       false,
		Ampq:         false,
		AccrualUrl:   false,
		DBConnection: false,
		Addr:         false,
	}
}
