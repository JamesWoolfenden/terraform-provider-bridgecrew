---
layout: "bridgecrew"
page_title: "Bridgecrew: data_source_users"
sidebar_current: "docs-bridgecrew-data_source_users"

description: |-
Get a list of all your Bridgecrew platform users.

---

# bridgecrew_users

Use this datasource to get the details of your users from Bridgecrew.




## Example Usage
```hcl
data "bridgecrew_users" "myusers" {}
```
<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- **id** (String) The ID of this resource.

### Read-Only

- **users** (List of Object) (see [below for nested schema](#nestedatt--users))

<a id="nestedatt--users"></a>
### Nested Schema for `users`

Read-Only:

- **accounts** (List of String)
- **customername** (String)
- **email** (String)
- **lastmodified** (Number)
- **role** (String)