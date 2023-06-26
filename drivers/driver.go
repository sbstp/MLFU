package drivers

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/samber/lo"
)

var (
	drivers = make(map[string]Driver)
)

type Config struct {
	Driver   string `json:"driver"`
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoadConfig(path string) (*Config, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(buf, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

type Driver interface {
	Name() string
	AddMagnetURL(config *Config, magnet string) error
}

func registerDriver(d Driver) {
	drivers[strings.ToLower(d.Name())] = d
}

func ListDrivers() []string {
	return lo.Map(lo.Values(drivers), func(d Driver, _ int) string {
		return d.Name()
	})
}

func GetDriver(name string) Driver {
	return drivers[strings.ToLower(name)]
}
