package util

import (
	"fmt"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	Config := LoadConfig("../")
	fmt.Printf("Type : %v, Value : %v\n", reflect.TypeOf(Config.AccessTokenDuration), Config.AccessTokenDuration)
}
