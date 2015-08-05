package volumeattachmentcommands

import "github.com/rackspace/rack/internal/github.com/codegangsta/cli"

var commandPrefix = "servers volume-attachment"
var serviceClientType = "compute"

// Get returns all the commands allowed for a `servers volume-attachment` request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		create,
		get,
		remove,
	}
}
