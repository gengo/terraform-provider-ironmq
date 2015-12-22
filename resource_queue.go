package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func queueSchema() *schema.Resource {
	return &schema.Resource{
		Create: updateQueue,
		Read:   readQueue,
		//Update: updateQueue,
		Delete: deleteQueue,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "name of the queue",
			},
		},
	}
}

func updateQueue(data *schema.ResourceData, configured interface{}) error {
	return nil
}

func readQueue(data *schema.ResourceData, configured interface{}) error {
	return nil
}

func deleteQueue(data *schema.ResourceData, configured interface{}) error {
	return nil
}
