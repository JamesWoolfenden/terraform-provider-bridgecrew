package bridgecrew

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEnforcementRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnforcementRuleRead,
		Schema: map[string]*schema.Schema{
			"accountid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account ID",
			},
			"codecategories": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"supplychain": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"softfailthreshold": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hardfailthreshold": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"commentsbotthreshold": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"secrets": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"softfailthreshold": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hardfailthreshold": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"commentsbotthreshold": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"iac": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"softfailthreshold": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hardfailthreshold": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"commentsbotthreshold": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"images": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"softfailthreshold": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hardfailthreshold": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"commentsbotthreshold": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"opensource": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"softfailthreshold": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hardfailthreshold": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"commentsbotthreshold": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

//goland:noinspection GoUnusedParameter,GoLinter,GoLinter
func dataSourceEnforcementRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	id := d.Get("accountid").(string)
	request := "%s/enforcement-rules/account/" + id
	params := RequestParams{request, "v1", "GET"}

	configure := m.(ProviderConfig)
	client, req, diagnostics, done := authClient(params, configure, nil)

	if done {
		return diagnostics
	}

	r, err := client.Do(req)

	if err != nil {
		log.Fatal("Failed at client.Do")
	}

	//goland:noinspection GoUnhandledErrorResult
	defer r.Body.Close()

	Enforcement := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&Enforcement)

	if err != nil {
		if err.Error() == "EOF" {
			temp := "no data found"
			err = errors.New(temp)
			log.Print(temp)
		} else {
			log.Println("Failed to parse data")
		}
		return diag.FromErr(err)
	}

	if err := flattenEnforcementRule(Enforcement, d); err != nil {
		return err
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diagnostics
}

func flattenEnforcementRule(enforcement map[string]interface{}, d *schema.ResourceData) diag.Diagnostics {
	accountid := enforcement["accountId"]

	if err := d.Set("accountid", accountid.(string)); err != nil {
		return diag.FromErr(err)
	}

	codecategories := make([]interface{}, 0)
	mycat := make(map[string]interface{})
	mycode := enforcement["codeCategories"].(map[string]interface{})

	supplies := setcategories(mycode, "SUPPLY_CHAIN")
	secrets := setcategories(mycode, "SECRETS")
	iac := setcategories(mycode, "IAC")
	images := setcategories(mycode, "IMAGES")
	opensource := setcategories(mycode, "OPEN_SOURCE")

	mycat["supplychain"] = supplies
	mycat["secrets"] = secrets
	mycat["iac"] = iac
	mycat["images"] = images
	mycat["opensource"] = opensource

	codecategories = append(codecategories, mycat)

	if err := d.Set("codecategories", codecategories); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
