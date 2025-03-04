---
subcategory: "Tencent Cloud Service Engine(TSE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tse_cngw_service"
sidebar_current: "docs-tencentcloud-resource-tse_cngw_service"
description: |-
  Provides a resource to create a tse cngw_service
---

# tencentcloud_tse_cngw_service

Provides a resource to create a tse cngw_service

## Example Usage

```hcl
resource "tencentcloud_tse_cngw_service" "cngw_service" {
  gateway_id = "gateway-ddbb709b"
  name       = "terraform-test"
  path       = "/test"
  protocol   = "http"
  retries    = 5
  tags = {
    "created" = "terraform"
  }
  timeout       = 6000
  upstream_type = "IPList"

  upstream_info {
    algorithm                   = "round-robin"
    auto_scaling_cvm_port       = 80
    auto_scaling_group_id       = "asg-519acdug"
    auto_scaling_hook_status    = "Normal"
    auto_scaling_tat_cmd_status = "Normal"
    port                        = 0
    slow_start                  = 20

    targets {
      host   = "192.168.0.1"
      port   = 80
      weight = 100
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String) gateway ID.
* `name` - (Required, String) service name.
* `path` - (Required, String) path.
* `protocol` - (Required, String) protocol. Reference value:`https`, `http`, `tcp`, `udp`.
* `retries` - (Required, Int) retry times.
* `timeout` - (Required, Int) time out, unit:ms.
* `upstream_info` - (Required, List) service config information.
* `upstream_type` - (Required, String) service type. Reference value:`Kubernetes`, `Registry`, `IPList`, `HostIP`, `Scf`.
* `tags` - (Optional, Map) Tag description list.

The `targets` object supports the following:

* `host` - (Required, String) host.
* `port` - (Required, Int) port.
* `weight` - (Required, Int) weight.
* `source` - (Optional, String) source of target.

The `upstream_info` object supports the following:

* `algorithm` - (Optional, String) load balance algorithm,default: `round-robin`, `least-connections` and `consisten_hashing` also support.
* `auto_scaling_cvm_port` - (Optional, Int) auto scaling group port of cvm.
* `auto_scaling_group_id` - (Optional, String) auto scaling group ID of cvm.
* `auto_scaling_hook_status` - (Optional, String) hook status in auto scaling group of cvm.
* `auto_scaling_tat_cmd_status` - (Optional, String) tat cmd status in auto scaling group of cvm.
* `host` - (Optional, String) an IP address or domain name.
* `namespace` - (Optional, String) namespace.
* `port` - (Optional, Int) backend service port.valid values: `1` to `65535`.
* `real_source_type` - (Optional, String) exact source service type.
* `scf_lambda_name` - (Optional, String) scf lambda name.
* `scf_lambda_qualifier` - (Optional, String) scf lambda version.
* `scf_namespace` - (Optional, String) scf lambda namespace.
* `scf_type` - (Optional, String) scf lambda type.
* `service_name` - (Optional, String) the name of the service in registry or kubernetes.
* `slow_start` - (Optional, Int) slow start time, unit: `second`, when it is enabled, weight of the node is increased from 1 to the target value gradually.
* `source_id` - (Optional, String) service source ID.
* `source_name` - (Optional, String) the name of source service.
* `source_type` - (Optional, String) source service type.
* `targets` - (Optional, List) provided when service type is IPList.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `service_id` - service id.


## Import

tse cngw_service can be imported using the id, e.g.

```
terraform import tencentcloud_tse_cngw_service.cngw_service cngw_service_id
```

