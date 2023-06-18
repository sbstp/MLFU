package main

import (
	"log"

	"github.com/sbstp/MLFU/drivers"
)

func main() {
	driver := drivers.GetDriver(drivers.ListDrivers()[0])
	err := driver.AddMagnetURL(&drivers.Config{
		URL:      "http://192.168.0.200:17001",
		Username: "dummy",
		Password: "dummy",
	}, "")
	if err != nil {
		log.Fatal(err)
	}
}
