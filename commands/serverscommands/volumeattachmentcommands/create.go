package volumeattachmentcommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/fatih/structs"
	osVolumeAttach "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/compute/v2/extensions/volumeattach"
	"github.com/rackspace/rack/util"
)

var create = cli.Command{
	Name:        "create",
	Usage:       util.Usage(commandPrefix, "create", "[--server-id <serverID> | --server-name <serverName>] [--volume-id <volumeID> | --volume-name <volumeName> | --stdin volume-id]"),
	Description: "Creates a new volume attachment on the server",
	Action:      actionCreate,
	Flags:       commandoptions.CommandFlags(flagsCreate, keysCreate),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsCreate, keysCreate))
	},
}

func flagsCreate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "volume-id",
			Usage: "[optional; required if `stdin` or volume-name isn't provided] The ID of the volume to attach.",
		},
		cli.StringFlag{
			Name:  "volume-name",
			Usage: "[optional; required if `stdin` or `volume-id` isn't provided] The ID of the volume to attach.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `volume-id` or `volume-name` isn't provided] The field being piped into STDIN. Valid values are: volume-id",
		},
		cli.StringFlag{
			Name:  "server-id",
			Usage: "[optional; required if `server-name` isn't provided] The server ID to which attach the volume.",
		},
		cli.StringFlag{
			Name:  "server-name",
			Usage: "[optional; required if `server-id` isn't provided] The server name to which attach the volume.",
		},
		cli.StringFlag{
			Name:  "device",
			Usage: "[optional] The name of the device to which the volume will attach. Default is 'auto'.",
		},
	}
}

var keysCreate = []string{"ID", "Device", "VolumeID", "ServerID"}

type paramsCreate struct {
	opts     *osVolumeAttach.CreateOpts
	serverID string
}

type commandCreate handler.Command

func actionCreate(c *cli.Context) {
	command := &commandCreate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandCreate) Context() *handler.Context {
	return command.Ctx
}

func (command *commandCreate) Keys() []string {
	return keysCreate
}

func (command *commandCreate) ServiceClientType() string {
	return serviceClientType
}

func (command *commandCreate) HandleFlags(resource *handler.Resource) error {
	serverID, err := serverIDorName(command.Ctx)
	if err != nil {
		return err
	}

	c := command.Ctx.CLIContext
	opts := &osVolumeAttach.CreateOpts{
		Device: c.String("device"),
	}

	resource.Params = &paramsCreate{
		opts:     opts,
		serverID: serverID,
	}
	return nil
}

func (command *commandCreate) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsCreate).opts.VolumeID = item
	return nil
}

func (command *commandCreate) HandleSingle(resource *handler.Resource) error {
	volumeID, err := volumeIDorName(command.Ctx)
	if err != nil {
		return err
	}

	resource.Params.(*paramsCreate).opts.VolumeID = volumeID
	return nil
}

func (command *commandCreate) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsCreate)
	volumeAttachment, err := osVolumeAttach.Create(command.Ctx.ServiceClient, params.serverID, params.opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = structs.Map(volumeAttachment)
}

func (command *commandCreate) StdinField() string {
	return "volume-id"
}
