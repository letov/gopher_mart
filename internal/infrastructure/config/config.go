package config

import (
	"reflect"
)

type configSource struct {
	Args   preConfig
	Env    preConfig
	Dotenv preConfig
}

func NewConfig() *Config {
	cs := configSource{
		Args:   newArgs(),
		Env:    newEnv(),
		Dotenv: newDotenv(),
	}

	return &Config{
		Salt:   getPriorConfigValue(cs, "Salt").(string),
		JwtKey: getPriorConfigValue(cs, "JwtKey").(string),
	}
}

func getPriorConfigValue(cs configSource, fieldName string) interface{} {
	ev := getConfigValue(cs.Env, fieldName)
	if ev != nil {
		return ev
	}

	av := getConfigValue(cs.Args, fieldName)
	if av != nil {
		return av
	}

	return getConfigValue(cs.Dotenv, fieldName)
}

func getConfigValue(pre preConfig, fieldName string) interface{} {
	value := reflect.ValueOf(pre)
	fp := value.FieldByName(fieldName)
	if fp.IsNil() {
		return nil
	}
	return reflect.Indirect(fp).Interface()
}
