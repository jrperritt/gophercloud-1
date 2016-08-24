package object

import "gopkg.in/urfave/cli.v1"

var commandPrefix = "files object"

type ObjectV1Command struct{}

func (_ ObjectV1Command) ServiceClientType() string {
	return "files"
}

// Get returns all the commands allowed for a `servers instance` v2 request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
		upload,
		//download,
		//remove,
	}
}
