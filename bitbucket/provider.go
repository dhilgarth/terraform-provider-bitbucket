package bitbucket

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gobb "github.com/ktrysmt/go-bitbucket"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BITBUCKET_USERNAME", nil),
				Description: "Username to authenticate with BitBucket.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BITBUCKET_PASSWORD", nil),
				Description: "Password to authenticate with BitBucket.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"bitbucket_project":    dataSourceBitBucketProject(),
			"bitbucket_repository": dataSourceBitBucketRepository(),
			"bitbucket_workspace":  dataSourceBitBucketWorkspace(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"bitbucket_project":    resourceBitBucketProject(),
			"bitbucket_repository": resourceBitBucketRepository(),
		},

		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, resourceData *schema.ResourceData) (interface{}, diag.Diagnostics) {
	client := gobb.NewBasicAuth(
		resourceData.Get("username").(string),
		resourceData.Get("password").(string),
	)

	return client, nil
}
