package bridgecrew

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"

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
				ConflictsWith: []string{"file"},
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
			"file": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"conditions"},
			},

			"last_updated": {
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

	myPolicy, err := setPolicy(d)

	if err != nil {
		return diag.FromErr(err)
	}

	jsPolicy, err := json.Marshal(myPolicy)
	if err != nil {
		return diag.FromErr(err)
	}

	configure := m.(ProviderConfig)
	url := configure.URL + "/api/v1/policies"
	payload := strings.NewReader(string(jsPolicy))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("authorization", configure.Token)

	res, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}

	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	//set the ID from the post into the current object
	clean, err := strconv.Unquote(string(body))
	d.SetId(clean)

	resourcePolicyRead(ctx, d, m)

	return diags
}

func setPolicy(d *schema.ResourceData) (Policy, error) {
	myPolicy := Policy{}
	myPolicy.Benchmarks = setBenchmark(d)
	myPolicy.Category = d.Get("category").(string)

	filename, hasFilename := d.GetOk("file")

	//if the filename is set then this is a yaml policy
	if hasFilename {
		code, err := loadFileContent(filename.(string))
		if err != nil {
			return myPolicy, fmt.Errorf("unable to load %q: %w", filename.(string), err)
		}
		myPolicy.Code = string(code)
	} else {
		//Not a yaml policy
		conditions, err := setConditions(d)
		//Don't set if not set
		if err != nil {
			return myPolicy, fmt.Errorf("unable set conditions %q", err)
		}
		myPolicy.Conditions = conditions[0]
	}

	myPolicy.Provider = d.Get("cloud_provider").(string)
	myPolicy.Severity = d.Get("severity").(string)
	myPolicy.Title = d.Get("title").(string)
	myPolicy.Guidelines = d.Get("guidelines").(string)

	return myPolicy, nil
}

// loadFileContent returns contents of a file in a given path
func loadFileContent(v string) ([]byte, error) {
	filename, err := homedir.Expand(v)
	if err != nil {
		return nil, err
	}
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}

func setConditions(d *schema.ResourceData) ([]Conditions, error) {
	conditions := make([]Conditions, 0, 1)

	myConditions := d.Get("conditions").([]interface{})

	if len(myConditions) > 0 {
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
	} else {
		return nil, errors.New("no Conditions Set")
	}

	return conditions, nil
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
	if err != nil {
		return diag.FromErr(err)
	}

	// add authorization header to the req
	req.Header.Add("authorization", configure.Token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}

	//goland:noinspection GoUnhandledErrorResult
	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)
	typedjson, err := typed.Json(body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("cloud_provider", strings.ToLower(typedjson["provider"].(string)))
	d.Set("title", typedjson["title"].(string))
	d.Set("severity", strings.ToLower(typedjson["severity"].(string)))
	d.Set("category", strings.ToLower(typedjson["category"].(string)))

	err = d.Set("guidelines", typedjson["guideline"])
	if err != nil {
		return diag.FromErr(err)
	}

	//myconditions should be an array it currently a map
	//hence this fudge
	//todo: once you start passing around condition arrays
	//this can go
	myConditions := make([]interface{}, 1)
	myConditions[0] = typedjson["conditionQuery"]
	err = d.Set("conditions", myConditions)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("benchmarks", typedjson["benchmarks"])
	if err != nil {
		return diag.FromErr(err)
	}

	if typedjson["file"] != nil {
		err = d.Set("file", typedjson["file"].(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	var diags diag.Diagnostics

	return diags
}

// highlight is just to help with manual debugging, so you can find the lines
//goland:noinspection SpellCheckingInspection
func highlight(myPolicy interface{}) {
	log.Print("XXXXXXXXXXX")
	log.Print(myPolicy)
	log.Print("XXXXXXXXXXX")
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := &http.Client{Timeout: 60 * time.Second}

	policyID := d.Id()
	if policyChange(d) {
		myPolicy, err := setPolicy(d)

		if err != nil {
			return diag.FromErr(err)
		}

		jsPolicy, err := json.Marshal(myPolicy)
		if err != nil {
			return diag.FromErr(err)
		}

		configure := m.(ProviderConfig)

		payload := strings.NewReader(string(jsPolicy))
		req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/policies/%s", configure.URL, policyID), payload)

		if err != nil {
			return diag.FromErr(err)
		}

		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("authorization", configure.Token)

		res, err := client.Do(req)
		if err != nil {
			return diag.FromErr(err)
		}

		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return diag.FromErr(err)
		}

		//just retrieved so I should check the result
		clean, err := strconv.Unquote(string(body))
		highlight(clean)

		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Failed at clean quotes %s \n", err.Error()),
			})
			return diags
		}

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
		d.HasChange("file") //||
	//d.HasChange("last_updated")
}

func resourcePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 60 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	policyID := d.Id()
	configure := m.(ProviderConfig)
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/policies/%s", configure.URL, policyID), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("authorization", configure.Token)

	res, err := client.Do(req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed at client.Do %s \n", err.Error()),
		})
		return diags
	}

	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		highlight(string(body))
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to delete %s \n", err.Error()),
		})
		return diags
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
