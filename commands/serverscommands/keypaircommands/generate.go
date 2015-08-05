package keypaircommands

import (
	"fmt"
	"strings"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/fatih/structs"
	osKeypairs "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/compute/v2/keypairs"
	"github.com/rackspace/rack/util"
)

var generate = cli.Command{
	Name:        "generate",
	Usage:       util.Usage(commandPrefix, "generate", "[--name <keypairName> | --stdin name]"),
	Description: "Generates a keypair",
	Action:      actionGenerate,
	Flags:       commandoptions.CommandFlags(flagsGenerate, keysGenerate),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsGenerate, keysGenerate))
	},
}

func flagsGenerate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name of the keypair",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped into STDIN. Valid values are: name",
		},
	}
}

var keysGenerate = []string{"Name", "Fingerprint", "PublicKey", "PrivateKey"}

type paramsGenerate struct {
	opts *osKeypairs.CreateOpts
}

type commandGenerate handler.Command

func actionGenerate(c *cli.Context) {
	command := &commandGenerate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandGenerate) Context() *handler.Context {
	return command.Ctx
}

func (command *commandGenerate) Keys() []string {
	return keysGenerate
}

func (command *commandGenerate) ServiceClientType() string {
	return serviceClientType
}

func (command *commandGenerate) HandleFlags(resource *handler.Resource) error {
	resource.Params = &paramsGenerate{
		opts: &osKeypairs.CreateOpts{},
	}
	return nil
}

func (command *commandGenerate) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsGenerate).opts.Name = item
	return nil
}

func (command *commandGenerate) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"name"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsGenerate).opts.Name = command.Ctx.CLIContext.String("name")
	return err
}

func (command *commandGenerate) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsGenerate).opts
	keypair, err := keypairs.Create(command.Ctx.ServiceClient, opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = structs.Map(keypair)
}

func (command *commandGenerate) StdinField() string {
	return "name"
}

func (command *commandGenerate) Table(resource *handler.Resource) {
	output := []string{"PROPERTY\tVALUE",
		"Name\t\t%s",
		"Fingerprint\t%s",
		"PublicKey\t%s",
		"PrivateKey:\n%s",
	}
	kp := resource.Result.(map[string]interface{})
	resource.Result = fmt.Sprintf(strings.Join(output, "\n"), kp["Name"], kp["Fingerprint"], kp["PublicKey"], kp["PrivateKey"])
}
