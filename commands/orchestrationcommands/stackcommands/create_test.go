package stackcommands

import (
	"flag"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	osStacks "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/orchestration/v1/stacks"
	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
)

func TestCreateContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandCreate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestCreateKeys(t *testing.T) {
	cmd := &commandCreate{}
	expected := keysCreate
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestCreateServiceClientType(t *testing.T) {
	cmd := &commandCreate{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestCreateHandleFlags(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("disable-rollback", "", "")
	flagset.String("template-file", "", "")
	flagset.String("environment-file", "", "")
	flagset.String("timeout", "", "")
	flagset.String("parameters", "", "")
	flagset.Set("disable-rollback", "true")
	flagset.Set("template-file", "mytemplate.yaml")
	flagset.Set("environment-file", "myenvironment.yaml")
	flagset.Set("timeout", "300")
	flagset.Set("parameters", "img=foo,flavor=bar")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandCreate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	disableRollback := true
	templateOpts := new(osStacks.Template)
	templateOpts.URL = "mytemplate.yaml"
	environmentOpts := new(osStacks.Environment)
	environmentOpts.URL = "myenvironment.yaml"
	expected := &handler.Resource{
		Params: &paramsCreate{
			opts: &osStacks.CreateOpts{
				TemplateOpts:    templateOpts,
				EnvironmentOpts: environmentOpts,
				DisableRollback: &disableRollback,
				Timeout:         300,
				Parameters: map[string]string{
					"img":    "foo",
					"flavor": "bar",
				},
			},
		},
	}
	actual := &handler.Resource{}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, *expected.Params.(*paramsCreate).opts, *actual.Params.(*paramsCreate).opts)
}

func TestCreateHandleSingle(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("name", "", "")
	flagset.Set("name", "stack1")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandCreate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsCreate{
			opts: &osStacks.CreateOpts{
				Name: "stack1",
			},
		},
	}
	actual := &handler.Resource{
		Params: &paramsCreate{
			opts: &osStacks.CreateOpts{},
		},
	}
	err := cmd.HandleSingle(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsCreate).opts.Name, actual.Params.(*paramsCreate).opts.Name)
}

/* Disabled since there seems to be issues with stack list
func TestCreateExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/stacks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
            w.WriteHeader(http.StatusAccepted)
		    w.Header().Add("Content-Type", "application/json")
		    fmt.Fprintf(w, `{"stack": {"id": "3095aefc-09fb-4bc7-b1f0-f21a304e864c"}}`)
        } else if r.Method == "GET" {
            w.WriteHeader(http.StatusOK)
            w.Header().Add("Content-Type", "application/json")
            fmt.Fprint(w, `{"stacks": {"creation_time": "2014-06-03T20:59:46Z"}}`)
        }
	})
	cmd := &commandCreate{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
    templateOpts := new(osStacks.Template)
    templateOpts.Bin = 	[]byte(`"heat_template_version": "2014-10-16"`)
	actual := &handler.Resource{
		Params: &paramsCreate{
			opts: &osStacks.CreateOpts{
                Name: "stack1",
                TemplateOpts: templateOpts,
                Parameters: map[string]string{
                    "img":    "foo",
                    "flavor": "bar",
                },
			},
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}
*/
