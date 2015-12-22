package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/iron-io/iron_go3/config"
	"github.com/iron-io/iron_go3/mq"
)

func queueSchema() *schema.Resource {
	return &schema.Resource{
		Create: createQueue,
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

func createQueue(data *schema.ResourceData, meta interface{}) error {
	cfg := meta.(config.Settings)
	name := data.Get("name").(string)
	_, err := mq.ConfigCreateQueue(mq.QueueInfo{
		Name: name,
	}, &cfg)
	return err
}

func updateQueue(data *schema.ResourceData, meta interface{}) error {
	return nil
}

func readQueue(data *schema.ResourceData, ironcfg interface{}) error {
	return nil
}

func deleteQueue(data *schema.ResourceData, ironcfg interface{}) error {
	return nil
}
