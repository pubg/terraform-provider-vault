// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-vault/internal/consts"
	"github.com/hashicorp/terraform-provider-vault/internal/provider"
)

func kvSecretV2MetadataDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: provider.ReadContextWrapper(kvSecretV2MetadataDataSourceRead),

		Schema: map[string]*schema.Schema{
			consts.FieldMount: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Path where KV-V2 engine is mounted",
			},
			consts.FieldName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Full name of the secret, excluding the mount and metadata prefix",
			},
			consts.FieldPath: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Full path where the KV-V2 secret metadata is read.",
			},
			consts.FieldCASRequired: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether check-and-set is required for writes to this secret.",
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time at which the secret was created.",
			},
			"current_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Current version of the secret.",
			},
			consts.FieldCustomMetadata: {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "User-provided metadata for the secret.",
			},
			consts.FieldDeleteVersionAfter: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Duration after which new versions are deleted.",
			},
			consts.FieldMaxVersions: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum number of versions retained for the secret.",
			},
			"last_updated_by": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Identity and operation that last updated the secret.",
			},
			"oldest_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Oldest retained version of the secret.",
			},
			"updated_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time at which the secret was last updated.",
			},
			consts.FieldVersions: {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Map of version numbers to JSON-encoded version metadata.",
			},
		},
	}
}

func kvSecretV2MetadataDataSourceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := provider.GetClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	path := getKVV2Path(d.Get(consts.FieldMount).(string), d.Get(consts.FieldName).(string), consts.FieldMetadata)
	log.Printf("[DEBUG] Reading KVV2 secret metadata at %q from Vault", path)

	secret, err := client.Logical().Read(path)
	if err != nil {
		return diag.Errorf("error reading secret metadata %q from Vault: %s", path, err)
	}
	if secret == nil {
		return diag.Errorf("no secret metadata found at %q", path)
	}

	if err := d.Set(consts.FieldPath, path); err != nil {
		return diag.FromErr(err)
	}

	fields := []string{
		consts.FieldCASRequired,
		"created_time",
		"current_version",
		consts.FieldCustomMetadata,
		consts.FieldDeleteVersionAfter,
		consts.FieldMaxVersions,
		"last_updated_by",
		"oldest_version",
		"updated_time",
		consts.FieldVersions,
	}
	for _, field := range fields {
		value := secret.Data[field]
		if valueMap, ok := value.(map[string]interface{}); ok {
			value = serializeDataMapToString(valueMap)
		}
		if err := d.Set(field, value); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(path)
	return nil
}
