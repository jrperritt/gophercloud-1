package flavor

import "github.com/codegangsta/cli"

var commandPrefix = "servers flavor"

type FlavorV2Command struct{}

func (_ FlavorV2Command) ServiceClientType() string {
	return "compute"
}

// Get returns all the commands allowed for a `servers flavor` v2 request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
	}
}
