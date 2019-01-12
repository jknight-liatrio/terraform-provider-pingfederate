package pingfederate

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

//Provider does stuff
//
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Administrator",
				Description: descriptions["username"],
				DefaultFunc: schema.EnvDefaultFunc("PINGFEDERATE_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Pa55word",
				Description: descriptions["password"],
				DefaultFunc: schema.EnvDefaultFunc("PINGFEDERATE_PASSWORD", nil),
			},
			"context": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/pf-admin-api/v1",
				Description: descriptions["context"],
				DefaultFunc: schema.EnvDefaultFunc("PINGFEDERATE_CONTEXT", nil),
			},
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://localhost:9999",
				Description: descriptions["base_url"],
				DefaultFunc: schema.EnvDefaultFunc("PINGFEDERATE_BASEURL", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"pingfederate_oauth_auth_server_settings": resourcePingFederateOauthAuthServerSettingsResource(),
		},
		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"username": "The username for pingfederate API.",
		"password": "The password for pingfederate API.",
		"base_url": "The base url of the pingfederate API.",
		"context":  "The context path of the pingfederate API.",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &Config{
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		BaseURL:  d.Get("base_url").(string),
		Context:  d.Get("context").(string),
	}

	return config.Client()
}

// Takes the result of flatmap.Expand for an array of strings
// and returns a []string
func expandStringList(configured []interface{}) []*string {
	log.Printf("[INFO] expandStringList %d", len(configured))
	vs := make([]*string, 0, len(configured))
	for _, v := range configured {
		val := v.(string)
		if val != "" {
			vs = append(vs, &val)
			log.Printf("[DEBUG] Appending: %s", val)
		}
	}
	return vs
}

// Takes the result of flatmap.Expand for an array of strings
// and returns a []*int
func expandIntList(configured []interface{}) []*int {
	vs := make([]*int, 0, len(configured))
	for _, v := range configured {
		_, ok := v.(int)
		if ok {
			val := v.(int)
			vs = append(vs, &val)
		}
	}
	return vs
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool { return &v }

// Int is a helper routine that allocates a new int value
// to store v and returns a pointer to it.
func Int(v int) *int { return &v }

// Int64 is a helper routine that allocates a new int64 value
// to store v and returns a pointer to it.
func Int64(v int64) *int64 { return &v }

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string { return &v }

func setResourceDataString(d *schema.ResourceData, name string, data *string) error {
	if data != nil {
		if err := d.Set(name, *data); err != nil {
			return err
		}
	}
	return nil
}

func setResourceDataInt(d *schema.ResourceData, name string, data *int) error {
	if data != nil {
		if err := d.Set(name, *data); err != nil {
			return err
		}
	}
	return nil
}

func setResourceDataBool(d *schema.ResourceData, name string, data *bool) error {
	if data != nil {
		if err := d.Set(name, *data); err != nil {
			return err
		}
	}
	return nil
}
