package bridgecrew

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/karlseguin/typed"
)

func resourcePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyCreate,
		ReadContext:   resourcePolicyRead,
		UpdateContext: resourcePolicyUpdate,
		DeleteContext: resourcePolicyDelete,
		Schema: map[string]*schema.Schema{
			"cloud_provider": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: false,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"title": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"severity": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					switch val.(string) {
					case
						"critical",
						"high",
						"low",
						"medium":
						return
					}
					errs = append(errs, fmt.Errorf("%q Must be one of critical, high, medium or low", val))
					return
				},
			},
			"category": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					switch val.(string) {
					case
						"logging",
						"elasticsearch",
						"general",
						"storage",
						"encryption",
						"networking",
						"monitoring",
						"kubernetes",
						"serverless",
						"backup_and_recovery",
						"iam",
						"secrets",
						"public",
						"general_security":
						return
					}
					errs = append(errs,
						fmt.Errorf("%q Must be one of logging, elasticsearch, general, storage, encryption,"+
							" networking, monitoring, kubernetes, serverless, backup_and_recovery, backup_and_recovery, public,"+
							" general_security or iam", val))
					return
				},
			},
			"guidelines": {
				Type:     schema.TypeString,
				Required: true,
			},
			"conditions": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"code"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_types": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cond_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"attribute": {
							Type:     schema.TypeString,
							Required: true,
						},

						"operator": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"benchmarks": {
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cis_azure_v11": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cis_azure_v12": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cis_azure_v13": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cis_aws_v12": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cis_aws_v13": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cis_kubernetes_v15": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cis_kubernetes_v16": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cis_gcp_v11": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cis_gke_v11": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cis_docker_v11": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cis_eks_v11": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"code": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"conditions"},
			},
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := &http.Client{Timeout: 60 * time.Second}
	var diags diag.Diagnostics

	myPolicy := setPolicy(d)

	jsPolicy, err := json.Marshal(myPolicy)
	if err != nil {
		log.Fatal("json could no be written")
	}

	configure := m.(ProviderConfig)
	url := configure.URL + "/api/v1/policies"
	payload := strings.NewReader(string(jsPolicy))

	req, _ := http.NewRequest("POST", url, payload)
	highlight(url)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("authorization", configure.Token)

	res, err := client.Do(req)
	if err != nil {
		log.Print(err)
		log.Fatal("POST FAILED")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		//find out what the results of the post was in the message
		log.Print(err.Error())
		myResults := body
		log.Print(myResults)
		log.Fatal("IO Failure")
	}

	//set the ID from the post into the current object
	clean, err := strconv.Unquote(string(body))
	d.SetId(clean)

	resourcePolicyRead(ctx, d, m)

	return diags
}

func setPolicy(d *schema.ResourceData) Policy {
	myPolicy := Policy{}
	myPolicy.Benchmarks = setBenchmark(d)

	myPolicy.Category = d.Get("category").(string)
	myCode := d.Get("code").(string)
	if len(myCode) != 0 {
		myPolicy.Code = d.Get("code").(string)
	}

	myPolicy.Provider = d.Get("cloud_provider").(string)
	myPolicy.Severity = d.Get("severity").(string)
	myPolicy.Title = d.Get("title").(string)

	conditions := setConditions(d)
	myPolicy.Conditions = conditions[0]

	myPolicy.Guidelines = d.Get("guidelines").(string)
	return myPolicy
}

func setConditions(d *schema.ResourceData) []Conditions {
	conditions := make([]Conditions, 0, 1)

	myConditions := d.Get("conditions").([]interface{})
	for _, myCondition := range myConditions {
		temp := myCondition.(map[string]interface{})
		var Condition Conditions
		Condition.Value = temp["value"].(string)
		Condition.CondType = temp["cond_type"].(string)
		Condition.Attribute = temp["attribute"].(string)
		Condition.Operator = temp["operator"].(string)

		var myResources []string
		myResources = CastToStringList(temp["resource_types"].([]interface{}))
		Condition.ResourceTypes = myResources

		conditions = append(conditions, Condition)
	}
	return conditions
}

func setBenchmark(d *schema.ResourceData) Benchmark {
	myBenchmark := (d.Get("benchmarks").(*schema.Set)).List()

	var myItem Benchmark
	s := myBenchmark[0].(map[string]interface{})
	myItem.Cisawsv12 = CastToStringList(s["cis_aws_v12"].([]interface{}))
	myItem.Cisawsv13 = CastToStringList(s["cis_aws_v13"].([]interface{}))
	myItem.Cisazurev11 = CastToStringList(s["cis_azure_v11"].([]interface{}))
	myItem.Cisazurev12 = CastToStringList(s["cis_azure_v12"].([]interface{}))
	myItem.Cisazurev13 = CastToStringList(s["cis_azure_v13"].([]interface{}))
	myItem.Cisgcpv11 = CastToStringList(s["cis_gcp_v11"].([]interface{}))
	myItem.Ciskubernetesv15 = CastToStringList(s["cis_kubernetes_v15"].([]interface{}))
	myItem.Ciskubernetesv16 = CastToStringList(s["cis_kubernetes_v16"].([]interface{}))
	myItem.Cisdockerv11 = CastToStringList(s["cis_docker_v11"].([]interface{}))
	myItem.Ciseksv11 = CastToStringList(s["cis_eks_v11"].([]interface{}))
	myItem.Cisgkev11 = CastToStringList(s["cis_gke_v11"].([]interface{}))
	return myItem
}

// CastToStringList is a helper to work with conversion of types
// If there's a better way (most likely)?
func CastToStringList(temp []interface{}) []string {
	var versions []string
	for _, version := range temp {
		versions = append(versions, version.(string))
	}
	return versions
}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 60 * time.Second}

	policyID := d.Id()

	configure := m.(ProviderConfig)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/policies/%s", configure.URL, policyID), nil)

	// add authorization header to the req
	req.Header.Add("authorization", configure.Token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		log.Print("Failed to make get")
		log.Fatal(err.Error())
	}

	r, err := client.Do(req)

	if err != nil {
		log.Fatal("Failed at client.Do")
	}
	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)
	typedjson, err := typed.Json(body)
	if err != nil {
		log.Fatal("Failed at unmarshalling with typed")
	}

	d.Set("cloud_provider", strings.ToLower(typedjson["provider"].(string)))
	d.Set("title", typedjson["title"].(string))
	d.Set("severity", strings.ToLower(typedjson["severity"].(string)))
	d.Set("category", strings.ToLower(typedjson["category"].(string)))
	d.Set("guidelines", typedjson["guidelines"])
	d.Set("conditions", typedjson["conditions"])
	d.Set("benchmarks", typedjson["benchmarks"])
	d.Set("code", typedjson["code"])

	var diags diag.Diagnostics

	return diags
}

// highlight is just to help with manual debugging, so you can find the lines
func highlight(myPolicy interface{}) {
	log.Print("XXXXXXXXXXX")
	log.Print(myPolicy)
	log.Print("XXXXXXXXXXX")
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 60 * time.Second}

	policyID := d.Id()
	if policyChange(d) {
		myPolicy := setPolicy(d)

		jsPolicy, err := json.Marshal(myPolicy)
		if err != nil {
			highlight("Failed to Marshall")
			log.Fatal(err.Error())
		}

		configure := m.(ProviderConfig)

		payload := strings.NewReader(string(jsPolicy))
		req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/policies/%s", configure.URL, policyID), payload)

		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("authorization", configure.Token)

		res, err := client.Do(req)
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)

		if err != nil {
			//find out what the results of the post was in the message
			log.Print(err.Error())
			myResults := body
			log.Print(myResults)
			log.Fatal("IO Failure")
		}

		//just retrieved to check the result
		clean, err := strconv.Unquote(string(body))
		highlight(clean)
		highlight(policyID)

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}
	return resourcePolicyRead(ctx, d, m)
}

func policyChange(d *schema.ResourceData) bool {
	return d.HasChange("conditions") ||
		d.HasChange("cloud_provider") ||
		d.HasChange("title") ||
		d.HasChange("severity") ||
		d.HasChange("category") ||
		d.HasChange("guidelines") ||
		d.HasChange("benchmarks") ||
		d.HasChange("code") //||
	//d.HasChange("last_updated")
}

func resourcePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 60 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	policyID := d.Id()
	configure := m.(ProviderConfig)
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/policies/%s", configure.URL, policyID), nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("authorization", configure.Token)

	res, err := client.Do(req)
	if err != nil {
		log.Print(err)
		log.Fatal("DELETE FAILED")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if string(body) != "Deleted policy" {
		highlight(string(body))
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
