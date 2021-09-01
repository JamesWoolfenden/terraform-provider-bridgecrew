{{define "data_source_repositories_branches"}}---
layout: "bridgecrew"
page_title: "Bridgecrew: data_source_repositories_branches"
sidebar_current: "docs-bridgecrew-data_source_repositories_branches"

description: |-

---

# bridgecrew_repository_branches (Data Source)

Use this datasource to get the details of your managed repositories branches from Bridgecrew.




## Example Usage
```hcl
data "bridgecrew_repositories_branches" "mybranches" {
}
```
{{end}}