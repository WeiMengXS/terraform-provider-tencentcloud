/*
Provides a resource to create a wedata rule_template

Example Usage

```hcl
resource "tencentcloud_wedata_rule_template" "rule_template" {
  type                = 2
  name                = "fo test"
  quality_dim         = 3
  source_object_type  = 2
  description         = "for tf test"
  source_engine_types = [3]
  multi_source_flag   = false
  sql_expression      = "c2VsZWN0ICogZnJvbSBkYg=="
  project_id          = "1840731346428280832"
  where_flag          = false
}
```

Import

wedata rule_template can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_rule_template.rule_template rule_template_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudWedataRuleTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataRuleTemplateCreate,
		Read:   resourceTencentCloudWedataRuleTemplateRead,
		Update: resourceTencentCloudWedataRuleTemplateUpdate,
		Delete: resourceTencentCloudWedataRuleTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Template type. `1` means System template, `2` means Custom template.",
			},

			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Template name.",
			},

			"quality_dim": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Quality inspection dimensions. `1` Accuracy, `2` Uniqueness, `3` Completeness, `4` Consistency, `5` Timeliness, `6` Effectiveness.",
			},

			"source_object_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Source data object type. `1` Constant `2` Offline table level Offline field level.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Description of Template.",
			},

			"source_engine_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "The engine type corresponding to the source.",
			},

			"multi_source_flag": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to associate other library tables.",
			},

			"sql_expression": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "SQL Expression.",
			},

			"project_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Project ID.",
			},

			"where_flag": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "If add where.",
			},
		},
	}
}

func resourceTencentCloudWedataRuleTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_rule_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request        = wedata.NewCreateRuleTemplateRequest()
		response       = wedata.NewCreateRuleTemplateResponse()
		ruleTemplateId uint64
	)
	if v, ok := d.GetOkExists("type"); ok {
		request.Type = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("quality_dim"); ok {
		request.QualityDim = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("source_object_type"); ok {
		request.SourceObjectType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_engine_types"); ok {
		sourceEngineTypesSet := v.(*schema.Set).List()
		for i := range sourceEngineTypesSet {
			sourceEngineTypes := sourceEngineTypesSet[i].(int)
			request.SourceEngineTypes = append(request.SourceEngineTypes, helper.IntUint64(sourceEngineTypes))
		}
	}

	if v, ok := d.GetOkExists("multi_source_flag"); ok {
		request.MultiSourceFlag = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("sql_expression"); ok {
		request.SqlExpression = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("where_flag"); ok {
		request.WhereFlag = helper.Bool(v.(bool))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().CreateRuleTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create wedata ruleTemplate failed, reason:%+v", logId, err)
		return err
	}

	ruleTemplateId = *response.Response.Data
	d.SetId(helper.UInt64ToStr(ruleTemplateId))

	return resourceTencentCloudWedataRuleTemplateRead(d, meta)
}

func resourceTencentCloudWedataRuleTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_rule_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := WedataService{client: meta.(*TencentCloudClient).apiV3Conn}

	ruleTemplateId := d.Id()

	ruleTemplate, err := service.DescribeWedataRuleTemplateById(ctx, ruleTemplateId)
	if err != nil {
		return err
	}

	if ruleTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WedataRuleTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if ruleTemplate.Type != nil {
		_ = d.Set("type", ruleTemplate.Type)
	}

	if ruleTemplate.Name != nil {
		_ = d.Set("name", ruleTemplate.Name)
	}

	if ruleTemplate.QualityDim != nil {
		_ = d.Set("quality_dim", ruleTemplate.QualityDim)
	}

	if ruleTemplate.SourceObjectType != nil {
		_ = d.Set("source_object_type", ruleTemplate.SourceObjectType)
	}

	if ruleTemplate.Description != nil {
		_ = d.Set("description", ruleTemplate.Description)
	}

	if ruleTemplate.SourceEngineTypes != nil {
		_ = d.Set("source_engine_types", ruleTemplate.SourceEngineTypes)
	}

	if ruleTemplate.MultiSourceFlag != nil {
		_ = d.Set("multi_source_flag", ruleTemplate.MultiSourceFlag)
	}

	if ruleTemplate.SqlExpression != nil {
		_ = d.Set("sql_expression", ruleTemplate.SqlExpression)
	}

	if ruleTemplate.WhereFlag != nil {
		_ = d.Set("where_flag", ruleTemplate.WhereFlag)
	}

	return nil
}

func resourceTencentCloudWedataRuleTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_rule_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := wedata.NewModifyRuleTemplateRequest()

	ruleTemplateId := d.Id()

	request.TemplateId = helper.StrToUint64Point(ruleTemplateId)

	immutableArgs := []string{
		"type", "name", "quality_dim", "source_object_type",
		"description", "source_engine_types", "multi_source_flag",
		"sql_expression", "project_id", "where_flag",
	}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("type") {
		if v, ok := d.GetOkExists("type"); ok {
			request.Type = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("quality_dim") {
		if v, ok := d.GetOkExists("quality_dim"); ok {
			request.QualityDim = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("source_object_type") {
		if v, ok := d.GetOkExists("source_object_type"); ok {
			request.SourceObjectType = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("source_engine_types") {
		if v, ok := d.GetOk("source_engine_types"); ok {
			sourceEngineTypesSet := v.(*schema.Set).List()
			for i := range sourceEngineTypesSet {
				sourceEngineTypes := sourceEngineTypesSet[i].(int)
				request.SourceEngineTypes = append(request.SourceEngineTypes, helper.IntUint64(sourceEngineTypes))
			}
		}
	}

	if d.HasChange("multi_source_flag") {
		if v, ok := d.GetOkExists("multi_source_flag"); ok {
			request.MultiSourceFlag = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("sql_expression") {
		if v, ok := d.GetOk("sql_expression"); ok {
			request.SqlExpression = helper.String(v.(string))
		}
	}

	if d.HasChange("project_id") {
		if v, ok := d.GetOk("project_id"); ok {
			request.ProjectId = helper.String(v.(string))
		}
	}

	if d.HasChange("where_flag") {
		if v, ok := d.GetOkExists("where_flag"); ok {
			request.WhereFlag = helper.Bool(v.(bool))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().ModifyRuleTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update wedata ruleTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWedataRuleTemplateRead(d, meta)
}

func resourceTencentCloudWedataRuleTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_rule_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := WedataService{client: meta.(*TencentCloudClient).apiV3Conn}
	ruleTemplateId := d.Id()

	if err := service.DeleteWedataRuleTemplateById(ctx, ruleTemplateId); err != nil {
		return err
	}

	return nil
}
