/*
Provides a resource to create a monitor tmpCvmAgent

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_monitor_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "tf_monitor_subnet"
  cidr_block        = "10.0.1.0/24"
}


resource "tencentcloud_monitor_tmp_instance" "foo" {
  instance_name       = "tf-tmp-instance"
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = tencentcloud_subnet.subnet.id
  data_retention_time = 30
  zone                = var.availability_zone
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_monitor_tmp_cvm_agent" "foo" {
  instance_id = tencentcloud_monitor_tmp_instance.foo.id
  name        = "tf-agent"
}

```
Import

monitor tmpCvmAgent can be imported using the id, e.g.
```
$ terraform import tencentcloud_monitor_tmp_cvm_agent.tmpCvmAgent instance_id#agent_id
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
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMonitorTmpCvmAgent() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMonitorTmpCvmAgentRead,
		Create: resourceTencentCloudMonitorTmpCvmAgentCreate,
		//Update: resourceTencentCloudMonitorTmpCvmAgentUpdate,
		Delete: resourceTencentCloudMonitorTmpCvmAgentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance id.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Agent name.",
			},

			"agent_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Agent id.",
			},
		},
	}
}

func resourceTencentCloudMonitorTmpCvmAgentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_cvm_agent.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = monitor.NewCreatePrometheusAgentRequest()
		response *monitor.CreatePrometheusAgentResponse
	)

	var instanceId string

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().CreatePrometheusAgent(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create monitor tmpCvmAgent failed, reason:%+v", logId, err)
		return err
	}

	tmpCvmAgentId := *response.Response.AgentId

	d.SetId(strings.Join([]string{instanceId, tmpCvmAgentId}, FILED_SP))
	return resourceTencentCloudMonitorTmpCvmAgentRead(d, meta)
}

func resourceTencentCloudMonitorTmpCvmAgentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmpCvmAgent.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	tmpCvmAgent, err := service.DescribeMonitorTmpCvmAgent(ctx, ids[0], ids[1])

	if err != nil {
		return err
	}

	if tmpCvmAgent == nil {
		d.SetId("")
		return fmt.Errorf("resource `tmpCvmAgent` %s does not exist", ids[1])
	}

	if tmpCvmAgent.InstanceId != nil {
		_ = d.Set("instance_id", tmpCvmAgent.InstanceId)
	}

	if tmpCvmAgent.Name != nil {
		_ = d.Set("name", tmpCvmAgent.Name)
	}

	if tmpCvmAgent.AgentId != nil {
		_ = d.Set("agent_id", tmpCvmAgent.AgentId)
	}

	return nil
}

func resourceTencentCloudMonitorTmpCvmAgentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_cvm_agent.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
