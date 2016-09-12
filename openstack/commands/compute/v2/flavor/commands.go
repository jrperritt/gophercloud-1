package flavor

import "gopkg.in/urfave/cli.v1"

var commandPrefix = "compute flavor"

type FlavorV2Command struct{}

func (_ FlavorV2Command) ServiceClientType() string {
	return "compute"
}

// Get returns all the commands allowed for a `compute flavor` v2 request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
	}
}
