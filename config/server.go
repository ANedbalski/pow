package config

import "time"

type Server struct {
	Addr          string        `mapstructure:"addr"`
	POWDifficulty int           `mapstructure:"pow_difficulty"`
	POWTTL        time.Duration `mapstructure:"pow_ttl"`
}
