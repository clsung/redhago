package redhago

// Config describe what we want to test
type Config struct {
	RedisServer []RedisConfig `json:"redis,omitempty"`
}

type RedisConfig struct {
	Address        string `json:"address"`
	Password       string `json:"password,omitempty"`
	MaxIdleConn    int    `json:"max_idleconn,omitempty"`
	MaxIdleTimeout int    `json:"max_idletimeout,omitempty"`
}
