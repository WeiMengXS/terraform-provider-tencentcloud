package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudPostgresqlParameterTemplatesDataSource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccStepSetRegion(t, "ap-guangzhou")
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlParameterTemplatesDataSource,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-guangzhou")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_postgresql_parameter_templates.parameter_templates"),
					resource.TestCheckResourceAttr("data.tencentcloud_postgresql_parameter_templates.parameter_templates", "filters.#", "2"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_parameter_templates.parameter_templates", "list.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_postgresql_parameter_templates.parameter_templates", "list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_postgresql_parameter_templates.parameter_templates", "list.0.template_name", "tf_test_pg_temp_ds"),
				),
			},
		},
	})
}

const testAccPostgresqlParameterTemplatesDataSource = `

resource "tencentcloud_postgresql_parameter_template" "temp1" {
	template_name = "tf_test_pg_temp_ds"
	db_major_version = "13"
	db_engine = "postgresql"
	template_description = "For_tf_test"
  
	modify_param_entry_set {
	  name = "lc_time"
	  expected_value = "POSIX"
	}
	modify_param_entry_set {
	  name = "timezone"
	  expected_value = "PRC"
	}
  }

data "tencentcloud_postgresql_parameter_templates" "parameter_templates" {
  filters {
	name = "TemplateName"
	values = [tencentcloud_postgresql_parameter_template.temp1.template_name]
  }
  filters {
	name = "DBEngine"
	values = [tencentcloud_postgresql_parameter_template.temp1.db_engine]
  }
  order_by = "CreateTime"
  order_by_type = "desc"
}

`
