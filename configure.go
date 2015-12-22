package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/iron-io/iron_go3/config"
)

var (
	configurationSchema = map[string]*schema.Schema{
		"env": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "env name in iron.json",
			Default:     "",
		},
	}
)

func configure(data *schema.ResourceData) (client interface{}, err error) {
	defer func() {
		// config.ConfigWithEnv does panic on runtime error. Need to recover.
		if p := recover(); p != nil {
			err = fmt.Errorf("iron_mq configuration failure: %v", p)
		}
	}()
	cfg := config.ConfigWithEnv("iron_mq", data.Get("env").(string))
	return cfg, nil
}
