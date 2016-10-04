package loadbalancer

import (
	"github.com/gophercloud/cli/lib/traits"
	"gopkg.in/urfave/cli.v1"
)

var commandPrefix = "load-balancing load-balancer"

type LoadbalancerV2Command struct {
	traits.Commandable
	traits.Networkingable
}

// Get returns all the commands allowed for a `load-balancing load-balancer` v2 request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		//get,
		//update,
		create,
		//remove,
	}
}
