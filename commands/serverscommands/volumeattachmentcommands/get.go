package volumeattachmentcommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/fatih/structs"
	osVolumeAttach "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/compute/v2/extensions/volumeattach"
	"github.com/rackspace/rack/util"
)

var get = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", "[--server-id <serverID> | --server-name <serverName>] --id <attachmentID> "),
	Description: "Gets an existing volume attachment",
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
			Usage: "[required] The ID of the attachment.",
		},
		cli.StringFlag{
			Name:  "server-id",
			Usage: "[optional; required if `server-name` isn't provided] The server ID of the attachment.",
		},
		cli.StringFlag{
			Name:  "server-name",
			Usage: "[optional; required if `server-id` isn't provided] The server name of the attachment.",
		},
	}
}

var keysGet = []string{"ID", "Device", "VolumeID", "ServerID"}

type paramsGet struct {
	volumeID string
	serverID string
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
	serverID, err := serverIDorName(command.Ctx)
	if err != nil {
		return err
	}

	err = command.Ctx.CheckFlagsSet([]string{"id"})
	if err != nil {
		return err
	}

	resource.Params = &paramsGet{
		volumeID: command.Ctx.CLIContext.String("id"),
		serverID: serverID,
	}
	return nil
}

func (command *commandGet) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsGet)
	volumeAttachment, err := osVolumeAttach.Get(command.Ctx.ServiceClient, params.serverID, params.volumeID).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = structs.Map(volumeAttachment)
}
