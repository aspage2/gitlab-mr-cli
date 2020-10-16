package main

import (
	"fmt"

	"gitlab.com/mintel/personal-dev/apage/glmr/stuff"
	"github.com/spf13/viper"
)

func main() {
	viper.SetDefault("Key", "value")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName(".glmr")
	fmt.Printf("Hello, world %d\n", stuff.MyThing(1, 1))
}
