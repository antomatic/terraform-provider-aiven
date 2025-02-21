package aiven

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceKafka() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceServiceRead,
		Description: "The Kafka data source provides information about the existing Aiven Kafka services.",
		Schema:      resourceSchemaAsDatasourceSchema(aivenKafkaSchema(), "project", "service_name"),
	}
}
