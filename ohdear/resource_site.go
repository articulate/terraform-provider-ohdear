package ohdear

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	//"github.com/hashicorp/terraform/helper/validation"
	"github.com/smallnest/goreq"
	//"runtime"
	"github.com/davecgh/go-spew/spew"

)

type Site struct {
	Url    string `json:"url"`
	TeamId string `json:"team_id"`
}

func resourceSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteCreate,
		Read:   resourceSiteRead,
		Update: resourceSiteUpdate,
		Delete: resourceSiteDelete,
		Exists: resourceSiteExists,
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "URL of the site to be checked",
			},
			"team_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the team for this site",
			},
		},
	}
}


func resourceSiteExists(d *schema.ResourceData, m interface{}) (bool, error) {
	_, _, err := goreq.New().Get("https://ohdear.app/api/sites").
		SetHeader("Authorization", "Bearer "+m.(Config).Token).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		End()

	if (err == nil) {
		return true, nil
	} else {
		return false, nil
	}
}

func resourceSiteCreate(d *schema.ResourceData, m interface{}) error {
	newSite := Site{
		Url: d.Get("url").(string),
		TeamId: d.Get("team_id").(string),
	}

	encoded, _ := json.Marshal(newSite)

	_, body, _ := goreq.New().Post("https://ohdear.app/api/sites").
		SetHeader("Authorization", "Bearer "+m.(Config).Token).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SendRawString(string(encoded)).
		End()

	var datai interface{}
	err := json.Unmarshal([]byte(body), &datai)
	data, _ := datai.(map[string]interface{})

	if (err == nil) {
		s := fmt.Sprintf("%f", data["id"].(float64))
		d.SetId(s)
		//runtime.Breakpoint()
		spew.Dump(data)
	} else {
		fmt.Println(err)
	}
	return nil
}

func resourceSiteUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSiteRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSiteDelete(d *schema.ResourceData, m interface{}) error {
	newSite := Site{
		Url: d.Get("url").(string),
		TeamId: d.Get("team_id").(string),
	}

	encoded, _ := json.Marshal(newSite)

	resp, _, _ := goreq.New().Delete("https://ohdear.app/api/sites/" + d.Id()).
		SetHeader("Authorization", "Bearer "+m.(Config).Token).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SendRawString(string(encoded)).
		End()
	fmt.Println(&resp)

	return nil
}
