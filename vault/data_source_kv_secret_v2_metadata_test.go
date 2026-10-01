// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-vault/internal/consts"
	"github.com/hashicorp/terraform-provider-vault/testutil"
)

func TestDataSourceKVV2SecretMetadata(t *testing.T) {
	t.Parallel()
	mount := acctest.RandomWithPrefix("tf-kv")
	name := acctest.RandomWithPrefix("secret")
	resourceName := "data.vault_kv_secret_v2_metadata.test"

	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(context.Background(), t),
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testDataSourceKVV2SecretMetadataConfig(mount, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, consts.FieldMount, mount),
					resource.TestCheckResourceAttr(resourceName, consts.FieldName, name),
					resource.TestCheckResourceAttr(resourceName, consts.FieldPath, fmt.Sprintf("%s/metadata/%s", mount, name)),
					resource.TestCheckResourceAttr(resourceName, consts.FieldCASRequired, "false"),
					resource.TestCheckResourceAttr(resourceName, "current_version", "1"),
					resource.TestCheckResourceAttr(resourceName, consts.FieldCustomMetadata+".owner", "platform"),
					resource.TestCheckResourceAttr(resourceName, consts.FieldDeleteVersionAfter, "0s"),
					resource.TestCheckResourceAttr(resourceName, consts.FieldMaxVersions, "5"),
					resource.TestCheckResourceAttrSet(resourceName, "last_updated_by.actor"),
					resource.TestCheckResourceAttr(resourceName, "last_updated_by.operation", "create"),
					resource.TestCheckResourceAttr(resourceName, "oldest_version", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "created_time"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_time"),
					resource.TestCheckResourceAttr(resourceName, consts.FieldVersions+".%", "1"),
					resource.TestCheckResourceAttrSet(resourceName, consts.FieldVersions+".1"),
				),
			},
		},
	})
}

func testDataSourceKVV2SecretMetadataConfig(mount, name string) string {
	return fmt.Sprintf(`
%s

resource "vault_kv_secret_v2" "test" {
  mount               = vault_mount.kvv2.path
  name                = %q
  delete_all_versions = true
  data_json            = jsonencode({ value = "test" })

  custom_metadata {
    max_versions = 5
    data = {
      owner = "platform"
    }
  }
}

data "vault_kv_secret_v2_metadata" "test" {
  mount = vault_mount.kvv2.path
  name  = vault_kv_secret_v2.test.name
}
`, kvV2MountConfig(mount), name)
}
