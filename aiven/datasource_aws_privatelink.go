// Copyright (c) 2021 Aiven, Helsinki, Finland. https://aiven.io/
package aiven

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceAWSPrivatelink() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceAWSPrivatelinkRead,
		Description: "The AWS Privatelink resource allows the creation and management of Aiven AWS Privatelink for a services.",
		Schema:      resourceSchemaAsDatasourceSchema(aivenAWSPrivatelinkSchema, "project", "service_name"),
	}
}

func datasourceAWSPrivatelinkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	projectName := d.Get("project").(string)
	serviceName := d.Get("service_name").(string)
	d.SetId(buildResourceID(projectName, serviceName))

	return resourceAWSPrivatelinkRead(ctx, d, m)
}
