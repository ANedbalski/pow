package config

type Client struct {
	Server        string `mapstructure:"server_path"`
	POWDifficulty int    `mapstructure:"pow_difficulty"`
	MaxIterations int    `mapstructure:"max_iterations"`
}
