package main

import (
	"bytes"
	"fmt"
	"strings"

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
		Exists: queueExists,

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
	info, err := mq.ConfigCreateQueue(queueInfoFromData(data), &cfg)
	if err != nil {
		return err
	}
	refreshState(data, cfg.ProjectId, info)
	return nil
}

func updateQueue(data *schema.ResourceData, meta interface{}) error {
	cfg := meta.(config.Settings)
	info := queueInfoFromData(data)
	q := mq.ConfigNew(info.Name, &cfg)
	info, err := q.Update(info)
	refreshState(data, cfg.ProjectId, info)
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
	refreshState(data, cfg.ProjectId, info)
	return nil
}

func deleteQueue(data *schema.ResourceData, meta interface{}) error {
	cfg := meta.(config.Settings)
	q := mq.ConfigNew(queueInfoFromData(data).Name, &cfg)
	return q.Delete()
}

func queueExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	cfg := meta.(config.Settings)
	q := mq.ConfigNew(queueInfoFromData(data).Name, &cfg)
	_, err := q.Info()
	if err != nil {
		// TODO: avoid this hacky detection once mq client library supports QueueNotExist err.
		if strings.Contains(err.Error(), "Queue not found") {
			return false, nil
		}
		return false, err
	}
	return true, nil
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

// refreshState reflects the current configuration of the queue to "data" based
// on "project" and "info".
func refreshState(data *schema.ResourceData, project string, info mq.QueueInfo) {
	data.SetId(fmt.Sprintf("%s/%s", project, info.Name))
	if info.Type == nil || *info.Type == "pull" {
		data.Set("type", "pull")
		data.Set("push", nil)
		return
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
