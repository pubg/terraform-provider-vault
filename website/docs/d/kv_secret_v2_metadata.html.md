---
layout: "vault"
page_title: "Vault: vault_kv_secret_v2_metadata data source"
sidebar_current: "docs-vault-datasource-kv-secret-v2-metadata"
description: |-
  Reads metadata for a KV-V2 secret from a given path in Vault
---

# vault\_kv\_secret\_v2\_metadata

Reads metadata and version history for a KV-V2 secret without reading its secret data.

## Example Usage

```hcl
data "vault_kv_secret_v2_metadata" "example" {
  mount = "kvv2"
  name  = "apps/example"
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Optional) The namespace of the target resource.
  The value should not contain leading or trailing forward slashes.
  The `namespace` is always relative to the provider's configured [namespace](/docs/providers/vault/index.html#namespace).
  *Available only for Vault Enterprise*.

* `mount` - (Required) Path where the KV-V2 engine is mounted.

* `name` - (Required) Full name of the secret, excluding the mount and `metadata/` prefix.

## Required Vault Capabilities

Use of this data source requires the `read` capability on the secret's `metadata/` path.

## Attributes Reference

The following attributes are exported:

* `path` - Full path where the KV-V2 secret metadata is read.
* `cas_required` - Whether check-and-set is required for writes to this secret.
* `created_time` - Time at which the secret was created.
* `current_version` - Current version of the secret.
* `custom_metadata` - User-provided metadata for the secret.
* `delete_version_after` - Duration after which new versions are deleted.
* `max_versions` - Maximum number of versions retained for the secret.
* `last_updated_by` - Identity and operation that last updated the secret.
* `oldest_version` - Oldest retained version of the secret.
* `updated_time` - Time at which the secret was last updated.
* `versions` - Map of version numbers to JSON-encoded version metadata.
