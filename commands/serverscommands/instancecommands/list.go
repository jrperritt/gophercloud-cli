package instancecommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	osServers "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
	"github.com/rackspace/rack/util"
)

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", ""),
	Description: "Lists existing servers",
	Action:      actionList,
	Flags:       commandoptions.CommandFlags(flagsList, keysList),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsList, keysList))
	},
}

func flagsList() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "all-pages",
			Usage: "[optional] Return all servers. Default is to paginate.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional] Only list servers with this name.",
		},
		cli.StringFlag{
			Name:  "changes-since",
			Usage: "[optional] Only list servers that have been changed since this time/date stamp.",
		},
		cli.StringFlag{
			Name:  "image",
			Usage: "[optional] Only list servers that have this image ID.",
		},
		cli.StringFlag{
			Name:  "flavor",
			Usage: "[optional] Only list servers that have this flavor ID.",
		},
		cli.StringFlag{
			Name:  "status",
			Usage: "[optional] Only list servers that have this status.",
		},
		cli.StringFlag{
			Name:  "marker",
			Usage: "[optional] Start listing servers at this server ID.",
		},
		cli.IntFlag{
			Name:  "limit",
			Usage: "[optional] Only return this many servers at most.",
		},
	}
}

var keysList = []string{"ID", "Name", "Status", "PublicIPv4", "PrivateIPv4", "Image", "Flavor"}

type paramsList struct {
	opts     *osServers.ListOpts
	allPages bool
}

type commandList handler.Command

func actionList(c *cli.Context) {
	command := &commandList{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandList) Context() *handler.Context {
	return command.Ctx
}

func (command *commandList) Keys() []string {
	return keysList
}

func (command *commandList) ServiceClientType() string {
	return serviceClientType
}

func (command *commandList) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext
	opts := &osServers.ListOpts{
		ChangesSince: c.String("changes-since"),
		Image:        c.String("image"),
		Flavor:       c.String("flavor"),
		Name:         c.String("name"),
		Status:       c.String("status"),
		Marker:       c.String("marker"),
		Limit:        c.Int("limit"),
	}
	resource.Params = &paramsList{
		opts:     opts,
		allPages: c.Bool("all-pages"),
	}
	return nil
}

func (command *commandList) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsList).opts
	allPages := resource.Params.(*paramsList).allPages
	pager := servers.List(command.Ctx.ServiceClient, opts)
	if allPages {
		pages, err := pager.AllPages()
		if err != nil {
			resource.Err = err
			return
		}
		info, err := servers.ExtractServers(pages)
		if err != nil {
			resource.Err = err
			return
		}
		result := make([]map[string]interface{}, len(info))
		for j, server := range info {
			result[j] = serverSingle(&server)
		}
		resource.Result = result
	} else {
		limit := opts.Limit
		err := pager.EachPage(func(page pagination.Page) (bool, error) {
			info, err := servers.ExtractServers(page)
			if err != nil {
				return false, err
			}
			result := make([]map[string]interface{}, len(info))
			for j, server := range info {
				result[j] = serverSingle(&server)
			}
			resource.Result = result
			if len(info) >= limit {
				return false, nil
			}
			limit -= len(info)
			command.Ctx.Results <- resource
			return true, nil
		})
		if err != nil {
			resource.Err = err
			return
		}
	}
}
