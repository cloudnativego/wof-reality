package integrations_test

import (
	"fmt"
	"testing"

	"github.com/cloudfoundry-community/go-cfenv"
	. "github.com/cloudnativego/wof-reality/service"
)

var (
	appEnv, _ = cfenv.Current()
	server    = NewServer(appEnv)
)

func TestIntegration(t *testing.T) {
	fmt.Println("Integration test placeholder")
}
