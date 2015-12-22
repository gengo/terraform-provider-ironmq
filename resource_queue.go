package main

import (
	"fmt"

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
	_, err := mq.ConfigCreateQueue(queueInfoFromData(data), &cfg)
	if err != nil {
		return err
	}
	if err := readQueue(data, meta); err != nil {
		return err
	}
	return nil
}

func updateQueue(data *schema.ResourceData, meta interface{}) error {
	return nil
}

func readQueue(data *schema.ResourceData, meta interface{}) error {
	cfg := meta.(config.Settings)
	name := queueInfoFromData(data).Name
	data.SetId(fmt.Sprintf("%s/%s", cfg.ProjectId, name))
	return nil
}

func deleteQueue(data *schema.ResourceData, meta interface{}) error {
	cfg := meta.(config.Settings)
	q := mq.ConfigNew(queueInfoFromData(data).Name, &cfg)
	return q.Delete()
}

// queueInfoFromData constructs expected queueInfo based on a resource state given by Terraform.
//
// It assumes that "data" is valid against queueSchema.
func queueInfoFromData(data *schema.ResourceData) mq.QueueInfo {
	name := data.Get("name").(string)
	return mq.QueueInfo{
		Name: name,
	}
}
