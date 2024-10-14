package config

import (
	"flag"
)

func newArgs() preConfig {
	pre := preConfig{
		Salt:   flag.String("s", "", "Salt desc"),
		JwtKey: flag.String("j", "", "JwtKey desc"),
	}

	set := newSetConfig()

	flag.Parse()

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "s":
			set.Salt = true
		case "j":
			set.JwtKey = true
		}
	})

	if !set.Salt {
		pre.Salt = nil
	}
	if !set.JwtKey {
		pre.JwtKey = nil
	}

	return pre
}
