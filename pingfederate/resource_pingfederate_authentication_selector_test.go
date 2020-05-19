package pingfederate

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	pf "github.com/iwarapter/pingfederate-sdk-go/pingfederate"
)

func TestAccPingFederateAuthenticationSelectorResource(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		// PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPingFederateAuthenticationSelectorResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPingFederateAuthenticationSelectorResourceConfig("0.0.0.0/0"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPingFederateAuthenticationSelectorResourceExists("pingfederate_authentication_selector.demo"),
					// testAccCheckPingFederateAuthenticationSelectorResourceAttributes(),
				),
			},
			{
				Config: testAccPingFederateAuthenticationSelectorResourceConfig("127.0.0.1/32"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPingFederateAuthenticationSelectorResourceExists("pingfederate_authentication_selector.demo"),
				),
			},
		},
	})
}

func testAccCheckPingFederateAuthenticationSelectorResourceDestroy(s *terraform.State) error {
	return nil
}

func testAccPingFederateAuthenticationSelectorResourceConfig(configUpdate string) string {
	return fmt.Sprintf(`resource "pingfederate_authentication_selector" "demo" {
  name = "wee"
  plugin_descriptor_ref {
	id = "com.pingidentity.pf.selectors.cidr.CIDRAdapterSelector"
  }

  configuration {
	fields {
      name = "Result Attribute Name"
	  value = ""
	}
	tables {
	  name = "Networks"
	  rows {
		fields {
		  name  = "Network Range (CIDR notation)"
		  value = "%s"
		}
	  }
	}
  }
}`, configUpdate)
}

func testAccCheckPingFederateAuthenticationSelectorResourceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" || rs.Primary.ID == "0" {
			return fmt.Errorf("No rule ID is set")
		}

		conn := testAccProvider.Meta().(*pf.PfClient).AuthenticationSelectors
		result, _, err := conn.GetAuthenticationSelector(&pf.GetAuthenticationSelectorInput{Id: rs.Primary.ID})

		if err != nil {
			return fmt.Errorf("Error: AuthenticationSelector (%s) not found", n)
		}

		if *result.Name != rs.Primary.Attributes["name"] {
			return fmt.Errorf("Error: AuthenticationSelector response (%s) didnt match state (%s)", *result.Name, rs.Primary.Attributes["name"])
		}

		return nil
	}
}

func Test_resourcePingFederateAuthenticationSelectorResourceReadData(t *testing.T) {
	cases := []struct {
		Resource pf.AuthenticationSelector
	}{
		{
			Resource: pf.AuthenticationSelector{
				Name: String("foo"),
				Id:   String("foo"),
				PluginDescriptorRef: &pf.ResourceLink{
					Id: String("com.pingidentity.pf.selectors.cidr.CIDRAdapterSelector"),
				},
				Configuration: &pf.PluginConfiguration{
					Fields: &[]*pf.ConfigField{
						{
							Name:      String("Result Attribute Name"),
							Value:     String(""),
							Inherited: Bool(false),
						},
					},
					Tables: &[]*pf.ConfigTable{
						{
							Name:      String("Networks"),
							Inherited: Bool(false),
							Rows: &[]*pf.ConfigRow{
								{
									//DefaultRow: Bool(false),
									Fields: &[]*pf.ConfigField{
										{
											Name:      String("Network Range (CIDR notation)"),
											Value:     String("0.0.0.0/0"),
											Inherited: Bool(false),
										},
									},
								},
							},
						},
					},
				},
				//AttributeContract: &pf.AuthenticationSelectorAttributeContract{},
			},
		},
	}
	for i, tc := range cases {
		t.Run(fmt.Sprintf("tc:%v", i), func(t *testing.T) {

			descs := pf.PluginConfigDescriptor{
				ActionDescriptors: nil,
				Description:       nil,
				Fields:            &[]*pf.FieldDescriptor{},
				Tables: &[]*pf.TableDescriptor{
					{
						Columns: &[]*pf.FieldDescriptor{
							{
								Type: String("TEXT"),
								Name: String("Username"),
							},
						},
						Description:       nil,
						Label:             nil,
						Name:              String("Networks"),
						RequireDefaultRow: nil,
					},
				},
			}

			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				// Test request parameters
				equals(t, req.URL.String(), "/authenticationSelectors/descriptors/com.pingidentity.pf.selectors.cidr.CIDRAdapterSelector")
				// Send response to be tested
				b, _ := json.Marshal(pf.AuthenticationSelectorDescriptor{
					AttributeContract:        nil,
					ClassName:                String("com.pingidentity.pf.selectors.cidr.CIDRAdapterSelector"),
					ConfigDescriptor:         &descs,
					Id:                       String("com.pingidentity.pf.selectors.cidr.CIDRAdapterSelector"),
					Name:                     String("CIDR Authentication Selector"),
					SupportsExtendedContract: nil,
				})
				rw.Write(b)
			}))
			// Close the server when test finishes
			defer server.Close()

			// Use Client & URL from our local test server
			url, _ := url.Parse(server.URL)
			c := pf.NewClient("", "", url, "", server.Client())

			resourceSchema := resourcePingFederateAuthenticationSelectorResourceSchema()
			resourceLocalData := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{})
			resourcePingFederateAuthenticationSelectorResourceReadResult(resourceLocalData, &tc.Resource, c.AuthenticationSelectors)

			if got := *resourcePingFederateAuthenticationSelectorResourceReadData(resourceLocalData); !cmp.Equal(got, tc.Resource) {
				t.Errorf("resourcePingFederateAuthenticationSelectorResourceReadData() = %v", cmp.Diff(got, tc.Resource))
			}
		})
	}
}