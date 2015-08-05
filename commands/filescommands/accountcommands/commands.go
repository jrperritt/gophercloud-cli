package accountcommands

import "github.com/rackspace/rack/internal/github.com/codegangsta/cli"

var commandPrefix = "files account"
var serviceClientType = "object-store"

// Get returns all the commands allowed for a `files account` request.
func Get() []cli.Command {
	return []cli.Command{
		setMetadata,
		updateMetadata,
		getMetadata,
		deleteMetadata,
	}
}
