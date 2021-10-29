---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "bridgecrew_simple_policy Resource - terraform-provider-bridgecrew"
subcategory: ""
description: |-

---

# bridgecrew_simple_policy (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **category** (String)
- **cloud_provider** (String)
- **conditions** (Block List, Min: 1, Max: 1) (see [below for nested schema](#nestedblock--conditions))
- **guidelines** (String)
- **severity** (String)
- **title** (String)

### Optional

- **benchmarks** (Block Set, Max: 1) (see [below for nested schema](#nestedblock--benchmarks))
- **last_updated** (String)

### Read-Only

- **id** (String) The ID of this resource.

<a id="nestedblock--conditions"></a>
### Nested Schema for `conditions`

Required:

- **attribute** (String)
- **cond_type** (String)
- **operator** (String)
- **resource_types** (List of String)

Optional:

- **value** (String)


<a id="nestedblock--benchmarks"></a>
### Nested Schema for `benchmarks`

Optional:

- **cis_aws_v12** (List of String)
- **cis_aws_v13** (List of String)
- **cis_azure_v11** (List of String)
- **cis_azure_v12** (List of String)
- **cis_azure_v13** (List of String)
- **cis_docker_v11** (List of String)
- **cis_eks_v11** (List of String)
- **cis_gcp_v11** (List of String)
- **cis_gke_v11** (List of String)
- **cis_kubernetes_v15** (List of String)
- **cis_kubernetes_v16** (List of String)