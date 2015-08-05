package portcommands

import (
	"fmt"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	osPorts "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/networking/v2/ports"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/networking/v2/ports"
	"github.com/rackspace/rack/util"
)

var remove = cli.Command{
	Name:        "delete",
	Usage:       util.Usage(commandPrefix, "delete", ""),
	Description: "Deletes a port",
	Action:      actionDelete,
	Flags:       commandoptions.CommandFlags(flagsDelete, keysDelete),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsDelete, keysDelete))
	},
}

func flagsDelete() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `name` or `stdin` isn't provided] The ID of the port to delete.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` or `id` isn't provided] The name of the port to delete.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` or `id` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
	}
}

var keysDelete = []string{"ID", "Name", "Network ID", "Status", "MAC Address", "Device ID", "Device Owner", "Up", "Fixed IPs", "Security Groups"}

type paramsDelete struct {
	portID string
}

type commandDelete handler.Command

func actionDelete(c *cli.Context) {
	command := &commandDelete{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandDelete) Context() *handler.Context {
	return command.Ctx
}

func (command *commandDelete) Keys() []string {
	return keysDelete
}

func (command *commandDelete) ServiceClientType() string {
	return serviceClientType
}

func (command *commandDelete) HandleFlags(resource *handler.Resource) error {
	resource.Params = &paramsDelete{}
	return nil
}

func (command *commandDelete) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsDelete).portID = item
	return nil
}

func (command *commandDelete) HandleSingle(resource *handler.Resource) error {
	portID, err := command.Ctx.IDOrName(osPorts.IDFromName)
	if err != nil {
		return err
	}
	resource.Params.(*paramsDelete).portID = portID
	return nil
}

func (command *commandDelete) Execute(resource *handler.Resource) {
	portID := resource.Params.(*paramsDelete).portID
	err := ports.Delete(command.Ctx.ServiceClient, portID).ExtractErr()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = fmt.Sprintf("Successfully deleted port [%s]\n", portID)
}

func (command *commandDelete) StdinField() string {
	return "id"
}
