/*
Provides a resource to create a scf function_event_invoke_config

Example Usage

```hcl
resource "tencentcloud_scf_function_event_invoke_config" "function_event_invoke_config" {
  function_name = "keep-1676351130"
  namespace     = "default"
  async_trigger_config {
    retry_config {
      retry_num = 2
    }
    msg_ttl = 24
  }
}
```

Import

scf function_event_invoke_config can be imported using the id, e.g.

```
terraform import tencentcloud_scf_function_event_invoke_config.function_event_invoke_config function_name#namespace
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudScfFunctionEventInvokeConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudScfFunctionEventInvokeConfigCreate,
		Read:   resourceTencentCloudScfFunctionEventInvokeConfigRead,
		Update: resourceTencentCloudScfFunctionEventInvokeConfigUpdate,
		Delete: resourceTencentCloudScfFunctionEventInvokeConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Function name.",
			},

			"namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Default:     "default",
				Description: "Function namespace. Default value: default.",
			},
			"async_trigger_config": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Async retry configuration information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"retry_config": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Async retry configuration of function upon user error.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"retry_num": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Number of retry attempts.",
									},
								},
							},
						},
						"msg_ttl": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Message retention period.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudScfFunctionEventInvokeConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_function_event_invoke_config.create")()
	defer inconsistentCheck(d, meta)()

	functionName := d.Get("function_name").(string)
	namespace := d.Get("namespace").(string)

	d.SetId(functionName + FILED_SP + namespace)

	return resourceTencentCloudScfFunctionEventInvokeConfigUpdate(d, meta)
}

func resourceTencentCloudScfFunctionEventInvokeConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_function_event_invoke_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ScfService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	functionName := idSplit[0]
	namespace := idSplit[1]

	FunctionEventInvokeConfig, err := service.DescribeScfFunctionEventInvokeConfigById(ctx, namespace, functionName)
	if err != nil {
		return err
	}

	if FunctionEventInvokeConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ScfFunctionEventInvokeConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if FunctionEventInvokeConfig != nil {
		asyncTriggerConfigMap := map[string]interface{}{}

		if FunctionEventInvokeConfig.RetryConfig != nil {
			retryConfigList := []interface{}{}
			for _, retryConfig := range FunctionEventInvokeConfig.RetryConfig {
				retryConfigMap := map[string]interface{}{}

				if retryConfig.RetryNum != nil {
					retryConfigMap["retry_num"] = retryConfig.RetryNum
				}

				retryConfigList = append(retryConfigList, retryConfigMap)
			}

			asyncTriggerConfigMap["retry_config"] = retryConfigList
		}

		if FunctionEventInvokeConfig.MsgTTL != nil {
			asyncTriggerConfigMap["msg_ttl"] = FunctionEventInvokeConfig.MsgTTL
		}

		_ = d.Set("async_trigger_config", []interface{}{asyncTriggerConfigMap})
	}

	_ = d.Set("function_name", functionName)

	_ = d.Set("namespace", namespace)

	return nil
}

func resourceTencentCloudScfFunctionEventInvokeConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_function_event_invoke_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := scf.NewUpdateFunctionEventInvokeConfigRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	functionName := idSplit[0]
	namespace := idSplit[1]

	request.Namespace = &namespace
	request.FunctionName = &functionName

	if dMap, ok := helper.InterfacesHeadMap(d, "async_trigger_config"); ok {
		asyncTriggerConfig := scf.AsyncTriggerConfig{}
		if v, ok := dMap["retry_config"]; ok {
			for _, item := range v.([]interface{}) {
				retryConfigMap := item.(map[string]interface{})
				retryConfig := scf.RetryConfig{}
				if v, ok := retryConfigMap["retry_num"]; ok {
					retryConfig.RetryNum = helper.IntInt64(v.(int))
				}
				asyncTriggerConfig.RetryConfig = append(asyncTriggerConfig.RetryConfig, &retryConfig)
			}
		}
		if v, ok := dMap["msg_ttl"]; ok {
			asyncTriggerConfig.MsgTTL = helper.IntInt64(v.(int))
		}
		request.AsyncTriggerConfig = &asyncTriggerConfig
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseScfClient().UpdateFunctionEventInvokeConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update scf FunctionEventInvokeConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudScfFunctionEventInvokeConfigRead(d, meta)
}

func resourceTencentCloudScfFunctionEventInvokeConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_function_event_invoke_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
