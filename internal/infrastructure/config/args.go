package config

import (
	"flag"
)

func newArgs() preConfig {
	pre := preConfig{
		Salt:         flag.String("s", "", "Salt desc"),
		JwtKey:       flag.String("jk", "", "JwtKey desc"),
		Ampq:         flag.String("ampq", "", "Ampq desc"),
		AccrualUrl:   flag.String("r", "", "AccrualUrl desc"),
		DBConnection: flag.String("d", "", "DBConnection desc"),
		Addr:         flag.String("a", "", "Addr desc"),
	}

	set := newSetConfig()

	flag.Parse()

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "s":
			set.Salt = true
		case "jk":
			set.JwtKey = true
		case "ampq":
			set.Ampq = true
		case "r":
			set.AccrualUrl = true
		case "d":
			set.DBConnection = true
		case "a":
			set.Addr = true
		}
	})

	if !set.Salt {
		pre.Salt = nil
	}
	if !set.JwtKey {
		pre.JwtKey = nil
	}
	if !set.Ampq {
		pre.Ampq = nil
	}
	if !set.AccrualUrl {
		pre.AccrualUrl = nil
	}
	if !set.DBConnection {
		pre.DBConnection = nil
	}
	if !set.Addr {
		pre.Addr = nil
	}

	return pre
}
