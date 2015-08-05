package imagecommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/fatih/structs"
	osImages "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/compute/v2/images"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/compute/v2/images"
	"github.com/rackspace/rack/util"
)

var get = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", "[--id <serverID> | --name <serverName> | --stdin id]"),
	Description: "Retreives an image",
	Action:      actionGet,
	Flags:       commandoptions.CommandFlags(flagsGet, keysGet),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsGet, keysGet))
	},
}

func flagsGet() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` or `name` isn't provided] The ID of the image.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `id` or `stdin` isn't provided] The name of the image.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` or `name` isn't provided] The field being piped to STDIN. Valid values are: id",
		},
	}
}

var keysGet = []string{"ID", "Name", "Status", "Progress", "MinDisk", "MinRAM", "Created", "Updated", "Metadata"}

type paramsGet struct {
	image string
}

type commandGet handler.Command

func actionGet(c *cli.Context) {
	command := &commandGet{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandGet) Context() *handler.Context {
	return command.Ctx
}

func (command *commandGet) Keys() []string {
	return keysGet
}

func (command *commandGet) ServiceClientType() string {
	return serviceClientType
}

func (command *commandGet) HandleFlags(resource *handler.Resource) error {
	resource.Params = &paramsGet{}
	return nil
}

func (command *commandGet) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsGet).image = item
	return nil
}

func (command *commandGet) HandleSingle(resource *handler.Resource) error {
	id, err := command.Ctx.IDOrName(osImages.IDFromName)
	resource.Params.(*paramsGet).image = id
	return err
}

func (command *commandGet) Execute(resource *handler.Resource) {
	imageID := resource.Params.(*paramsGet).image
	image, err := images.Get(command.Ctx.ServiceClient, imageID).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = structs.Map(image)
}

func (command *commandGet) StdinField() string {
	return "id"
}

func (command *commandGet) PreCSV(resource *handler.Resource) {
	resource.FlattenMap("Metadata")
}

func (command *commandGet) PreTable(resource *handler.Resource) {
	command.PreCSV(resource)
}
