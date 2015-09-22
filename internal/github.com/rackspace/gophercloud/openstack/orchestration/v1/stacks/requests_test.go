package stacks

import (
	"testing"

	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper/client"
)

func TestCreateStack(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateSuccessfully(t, CreateOutput)
	template := new(Template)
	template.Bin = []byte(`
		{
			"heat_template_version": "2013-05-23",
			"description": "Simple template to test heat commands",
			"parameters": {
				"flavor": {
					"default": "m1.tiny",
					"type": "string"
				}
			}
		}`)
	createOpts := CreateOpts{
		Name:            "stackcreated",
		Timeout:         60,
		TemplateOpts:    template,
		DisableRollback: Disable,
	}
	actual, err := Create(fake.ServiceClient(), createOpts).Extract()
	th.AssertNoErr(t, err)

	expected := CreateExpected
	th.AssertDeepEquals(t, expected, actual)
}

func TestAdoptStack(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateSuccessfully(t, CreateOutput)

	adoptOpts := AdoptOpts{
		AdoptStackData:  `{environment{parameters{}}}`,
		Name:            "stackcreated",
		Timeout:         60,
		DisableRollback: Disable,
	}
	actual, err := Adopt(fake.ServiceClient(), adoptOpts).Extract()
	th.AssertNoErr(t, err)

	expected := CreateExpected
	th.AssertDeepEquals(t, expected, actual)
}

func TestListStack(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListSuccessfully(t, FullListOutput)

	count := 0
	err := List(fake.ServiceClient(), nil).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractStacks(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ListExpected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestGetStack(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetSuccessfully(t, GetOutput)

	actual, err := Get(fake.ServiceClient(), "postman_stack", "16ef0584-4458-41eb-87c8-0dc8d5f66c87").Extract()
	th.AssertNoErr(t, err)

	expected := GetExpected
	th.AssertDeepEquals(t, expected, actual)
}

func TestUpdateStack(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleUpdateSuccessfully(t)

	template := new(Template)
	template.Bin = []byte(`
		{
			"heat_template_version": "2013-05-23",
			"description": "Simple template to test heat commands",
			"parameters": {
				"flavor": {
					"default": "m1.tiny",
					"type": "string"
				}
			}
		}`)
	updateOpts := UpdateOpts{
		TemplateOpts: template,
	}
	err := Update(fake.ServiceClient(), "gophercloud-test-stack-2", "db6977b2-27aa-4775-9ae7-6213212d4ada", updateOpts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestDeleteStack(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteSuccessfully(t)

	err := Delete(fake.ServiceClient(), "gophercloud-test-stack-2", "db6977b2-27aa-4775-9ae7-6213212d4ada").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestPreviewStack(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandlePreviewSuccessfully(t, GetOutput)

	template := new(Template)
	template.Bin = []byte(`
		{
			"heat_template_version": "2013-05-23",
			"description": "Simple template to test heat commands",
			"parameters": {
				"flavor": {
					"default": "m1.tiny",
					"type": "string"
				}
			}
		}`)
	previewOpts := PreviewOpts{
		Name:            "stackcreated",
		Timeout:         60,
		TemplateOpts:    template,
		DisableRollback: Disable,
	}
	actual, err := Preview(fake.ServiceClient(), previewOpts).Extract()
	th.AssertNoErr(t, err)

	expected := PreviewExpected
	th.AssertDeepEquals(t, expected, actual)
}

func TestAbandonStack(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleAbandonSuccessfully(t, AbandonOutput)

	actual, err := Abandon(fake.ServiceClient(), "postman_stack", "16ef0584-4458-41eb-87c8-0dc8d5f66c8").Extract()
	th.AssertNoErr(t, err)

	expected := AbandonExpected
	th.AssertDeepEquals(t, expected, actual)
}
