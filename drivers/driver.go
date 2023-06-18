package drivers

import (
	"github.com/samber/lo"
)

var (
	drivers = make(map[string]Driver)
)

type Config struct {
	URL      string
	Username string
	Password string
}

type Driver interface {
	Name() string
	AddMagnetURL(config *Config, magnet string) error
}

func registerDriver(driver Driver) {
	drivers[driver.Name()] = driver
}

func ListDrivers() []string {
	return lo.Keys(drivers)
}

func GetDriver(name string) Driver {
	return drivers[name]
}
