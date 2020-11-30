package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"testing"
)

var testconfig = `
lbrycrd:
  A: "rpc://lbry:lbry@localhost:9345"
  B: "rpc://lbry:lbry@localhost:9445"
  C: "rpc://lbry:lbry@localhost:9545"
  D: "rpc://lbry:lbry@localhost:9645"
`

func TestInitializeConfiguration(t *testing.T) {
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(strings.NewReader(testconfig))
	if err != nil {
		t.Fatal(err)
	}

	instances := viper.GetStringMapString("lbrycrd")
	v, ok := instances["a"] // viper puts all keys to lowercase
	if !ok {
		t.Error("not found")
	}
	if v != "rpc://lbry:lbry@localhost:9345" {
		t.Error(fmt.Sprintf("expected %s got %s", "rpc://lbry:lbry@localhost:9345", v))
	}
}
