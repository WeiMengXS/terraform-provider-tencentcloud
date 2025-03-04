package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixVpcResumeSnapshotInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcResumeSnapshotInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_resume_snapshot_instance.resume_snapshot_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_resume_snapshot_instance.resume_snapshot_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcResumeSnapshotInstance = `

resource "tencentcloud_vpc_resume_snapshot_instance" "resume_snapshot_instance" {
  snapshot_policy_id = "sspolicy-1t6cobbv"
  snapshot_file_id = "ssfile-test"
  instance_id = "policy-1t6cob"
}

`
