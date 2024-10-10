package config

import (
	"flag"
)

func newArgs() preConfig {
	pre := preConfig{
		Salt: flag.String("s", "", "Salt desc"),
	}

	set := newSetConfig()

	flag.Parse()

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "s":
			set.Salt = true
		}
	})

	if !set.Salt {
		pre.Salt = nil
	}

	return pre
}
