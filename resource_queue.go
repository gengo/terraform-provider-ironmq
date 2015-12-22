package main

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/iron-io/iron_go3/config"
	"github.com/iron-io/iron_go3/mq"
)

func pushConfigHash(v interface{}) int {
	cfg := v.(map[string]interface{})
	var buf bytes.Buffer
	for _, item := range cfg["subscribers"].([]interface{}) {
		s := item.(map[string]interface{})
		buf.WriteString(s["url"].(string))
		buf.WriteRune(0)
	}
	return hashcode.String(buf.String())
}

func queueSchema() *schema.Resource {
	subscriber := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "URL of a queue subscriber endpoint",
			},
		},
	}
	pushConfig := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"subscribers": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     subscriber,
			},
		},
	}
	return &schema.Resource{
		Create: createQueue,
		Read:   readQueue,
		Update: updateQueue,
		Delete: deleteQueue,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "name of the queue",
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "pull",
				Description:  "type of the queue",
				ValidateFunc: validateQueueType,
			},
			"push": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     pushConfig,
				Set:      pushConfigHash,
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
	cfg := meta.(config.Settings)
	info := queueInfoFromData(data)
	q := mq.ConfigNew(info.Name, &cfg)
	_, err := q.Update(info)
	return err
}

func readQueue(data *schema.ResourceData, meta interface{}) error {
	cfg := meta.(config.Settings)
	name := queueInfoFromData(data).Name
	q := mq.ConfigNew(name, &cfg)
	info, err := q.Info()
	if err != nil {
		return err
	}

	data.SetId(fmt.Sprintf("%s/%s", cfg.ProjectId, name))
	if info.Type == nil || *info.Type == "pull" {
		data.Set("type", "pull")
		data.Set("push", nil)
		return nil
	}

	data.Set("type", *info.Type)

	var subscribers []interface{}
	for _, s := range info.Push.Subscribers {
		subscribers = append(subscribers, map[string]interface{}{
			"url": s.URL,
		})
	}
	push := map[string]interface{}{
		"subscribers": subscribers,
	}
	data.Set("push", schema.NewSet(pushConfigHash, []interface{}{push}))
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
	typ := data.Get("type").(string)
	var push *mq.PushInfo
	if typ != "pull" {
		push = new(mq.PushInfo)
		m := data.Get("push").(*schema.Set).List()[0].(map[string]interface{})
		for _, item := range m["subscribers"].([]interface{}) {
			s := item.(map[string]interface{})
			push.Subscribers = append(push.Subscribers, mq.QueueSubscriber{
				URL: s["url"].(string),
			})
		}
	}
	return mq.QueueInfo{
		Name: name,
		Type: &typ,
		Push: push,
	}
}

func validateQueueType(value interface{}, key string) ([]string, []error) {
	switch value.(string) {
	case "pull", "unicast", "multicast":
		break
	default:
		return nil, []error{
			fmt.Errorf("%q must be either pull, unicast or multicast", key),
		}
	}
	return nil, nil
}
