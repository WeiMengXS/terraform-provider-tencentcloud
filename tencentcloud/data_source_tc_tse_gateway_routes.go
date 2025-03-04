/*
Use this data source to query detailed information of tse gateway_routes

Example Usage

```hcl
data "tencentcloud_tse_gateway_routes" "gateway_routes" {
  gateway_id   = "gateway-ddbb709b"
  service_name = "test"
  route_name   = "keep-routes"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTseGatewayRoutes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTseGatewayRoutesRead,
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "gateway ID.",
			},

			"service_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "service name.",
			},

			"route_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "route name.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "result.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"route_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "route list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "service ID.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "service name.",
									},
									"methods": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "method list.",
									},
									"paths": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "path list.",
									},
									"hosts": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "host list.",
									},
									"protocols": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "protocol list.",
									},
									"preserve_host": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "whether to keep the host when forwarding to the backend.",
									},
									"https_redirect_status_code": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "https redirection status code.",
									},
									"strip_path": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "whether to strip path when forwarding to the backend.",
									},
									"created_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "created time.",
									},
									"force_https": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "whether to enable forced HTTPS, no longer use.",
									},
									"service_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "service name.",
									},
									"service_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "service ID.",
									},
									"destination_ports": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
										Computed:    true,
										Description: "destination port for Layer 4 matching.",
									},
									"headers": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "the headers of route.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "key of header.",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "value of header.",
												},
											},
										},
									},
								},
							},
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "total count.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTseGatewayRoutesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tse_gateway_routes.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("gateway_id"); ok {
		paramMap["GatewayId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_name"); ok {
		paramMap["ServiceName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("route_name"); ok {
		paramMap["RouteName"] = helper.String(v.(string))
	}

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result *tse.KongServiceRouteList
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeTseGatewayRoutesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = response
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result.RouteList))
	kongServiceRouteListMap := map[string]interface{}{}
	if result != nil {

		if result.RouteList != nil {
			routeListList := []interface{}{}
			for _, routeList := range result.RouteList {
				routeListMap := map[string]interface{}{}

				if routeList.ID != nil {
					routeListMap["id"] = routeList.ID
				}

				if routeList.Name != nil {
					routeListMap["name"] = routeList.Name
				}

				if routeList.Methods != nil {
					routeListMap["methods"] = routeList.Methods
				}

				if routeList.Paths != nil {
					routeListMap["paths"] = routeList.Paths
				}

				if routeList.Hosts != nil {
					routeListMap["hosts"] = routeList.Hosts
				}

				if routeList.Protocols != nil {
					routeListMap["protocols"] = routeList.Protocols
				}

				if routeList.PreserveHost != nil {
					routeListMap["preserve_host"] = routeList.PreserveHost
				}

				if routeList.HttpsRedirectStatusCode != nil {
					routeListMap["https_redirect_status_code"] = routeList.HttpsRedirectStatusCode
				}

				if routeList.StripPath != nil {
					routeListMap["strip_path"] = routeList.StripPath
				}

				if routeList.CreatedTime != nil {
					routeListMap["created_time"] = routeList.CreatedTime
				}

				if routeList.ForceHttps != nil {
					routeListMap["force_https"] = routeList.ForceHttps
				}

				if routeList.ServiceName != nil {
					routeListMap["service_name"] = routeList.ServiceName
				}

				if routeList.ServiceID != nil {
					routeListMap["service_id"] = routeList.ServiceID
				}

				if routeList.DestinationPorts != nil {
					routeListMap["destination_ports"] = routeList.DestinationPorts
				}

				if routeList.Headers != nil {
					headersMap := map[string]interface{}{}

					if routeList.Headers.Key != nil {
						headersMap["key"] = routeList.Headers.Key
					}

					if routeList.Headers.Value != nil {
						headersMap["value"] = routeList.Headers.Value
					}

					routeListMap["headers"] = []interface{}{headersMap}
				}

				routeListList = append(routeListList, routeListMap)
				ids = append(ids, *routeList.ID)
			}

			kongServiceRouteListMap["route_list"] = routeListList
		}

		if result.TotalCount != nil {
			kongServiceRouteListMap["total_count"] = result.TotalCount
		}

		_ = d.Set("result", []interface{}{kongServiceRouteListMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), kongServiceRouteListMap); e != nil {
			return e
		}
	}
	return nil
}
