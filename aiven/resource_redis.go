package aiven

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func redisSchema() map[string]*schema.Schema {
	s := serviceCommonSchema()
	s[ServiceTypeRedis] = &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Redis server provided values",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{},
		},
	}
	s[ServiceTypeRedis+"_user_config"] = generateServiceUserConfiguration(ServiceTypeRedis)

	return s
}

func resourceRedis() *schema.Resource {
	return &schema.Resource{
		Description:   "The Redis resource allows the creation and management of Aiven Redis services.",
		CreateContext: resourceServiceCreateWrapper(ServiceTypeRedis),
		ReadContext:   resourceServiceRead,
		UpdateContext: resourceServiceUpdate,
		DeleteContext: resourceServiceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceServiceState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: redisSchema(),
	}
}
