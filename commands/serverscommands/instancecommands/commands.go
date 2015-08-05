package instancecommands

import "github.com/rackspace/rack/internal/github.com/codegangsta/cli"

var commandPrefix = "servers instance"
var serviceClientType = "compute"

// Get returns all the commands allowed for a `servers instance` request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		create,
		get,
		update,
		remove,
		reboot,
		rebuild,
		resize,
		listAddresses,
		listAddressesByNetwork,
		getMetadata,
		setMetadata,
		updateMetadata,
		deleteMetadata,
	}
}
