package softlayer

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

func TestAccSoftLayerSSHKey_Basic(t *testing.T) {
	var key datatypes.Security_Ssh_Key

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSoftLayerSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSoftLayerSSHKeyConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSoftLayerSSHKeyExists("softlayer_ssh_key.testacc_foobar", &key),
					testAccCheckSoftLayerSSHKeyAttributes(&key),
					resource.TestCheckResourceAttr(
						"softlayer_ssh_key.testacc_foobar", "label", "testacc_foobar"),
					resource.TestCheckResourceAttr(
						"softlayer_ssh_key.testacc_foobar", "public_key", testAccValidPublicKey),
					resource.TestCheckResourceAttr(
						"softlayer_ssh_key.testacc_foobar", "notes", "first_note"),
				),
			},

			{
				Config: testAccCheckSoftLayerSSHKeyConfig_updated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSoftLayerSSHKeyExists("softlayer_ssh_key.testacc_foobar", &key),
					resource.TestCheckResourceAttr(
						"softlayer_ssh_key.testacc_foobar", "label", "changed_name"),
					resource.TestCheckResourceAttr(
						"softlayer_ssh_key.testacc_foobar", "public_key", testAccValidPublicKey),
					resource.TestCheckResourceAttr(
						"softlayer_ssh_key.testacc_foobar", "notes", "changed_note"),
				),
			},
		},
	})
}

func testAccCheckSoftLayerSSHKeyDestroy(s *terraform.State) error {
	service := services.GetSecuritySshKeyService(testAccProvider.Meta().(*session.Session))

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "softlayer_ssh_key" {
			continue
		}

		keyId, _ := strconv.Atoi(rs.Primary.ID)

		// Try to find the key
		_, err := service.Id(keyId).GetObject()

		if err == nil {
			return fmt.Errorf("SSH key %d still exists", keyId)
		}
	}

	return nil
}

func testAccCheckSoftLayerSSHKeyAttributes(key *datatypes.Security_Ssh_Key) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if *key.Label != "testacc_foobar" {
			return fmt.Errorf("Bad name: %s", *key.Label)
		}

		return nil
	}
}

func testAccCheckSoftLayerSSHKeyExists(n string, key *datatypes.Security_Ssh_Key) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Record ID is set")
		}

		keyId, _ := strconv.Atoi(rs.Primary.ID)

		service := services.GetSecuritySshKeyService(testAccProvider.Meta().(*session.Session))
		foundKey, err := service.Id(keyId).GetObject()

		if err != nil {
			return err
		}

		if strconv.Itoa(int(*foundKey.Id)) != rs.Primary.ID {
			return fmt.Errorf("Record %d not found", keyId)
		}

		*key = foundKey

		return nil
	}
}

var testAccCheckSoftLayerSSHKeyConfig_basic = fmt.Sprintf(`
resource "softlayer_ssh_key" "testacc_foobar" {
    label = "testacc_foobar"
    notes = "first_note"
    public_key = "%s"
}`, testAccValidPublicKey)

var testAccCheckSoftLayerSSHKeyConfig_updated = fmt.Sprintf(`
resource "softlayer_ssh_key" "testacc_foobar" {
    label = "changed_name"
    notes = "changed_note"
    public_key = "%s"
}`, testAccValidPublicKey)

var testAccValidPublicKey = strings.TrimSpace(`
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCKVmnMOlHKcZK8tpt3MP1lqOLAcqcJzhsvJcjscgVERRN7/9484SOBJ3HSKxxNG5JN8owAjy5f9yYwcUg+JaUVuytn5Pv3aeYROHGGg+5G346xaq3DAwX6Y5ykr2fvjObgncQBnuU5KHWCECO/4h8uWuwh/kfniXPVjFToc+gnkqA+3RKpAecZhFXwfalQ9mMuYGFxn+fwn8cYEApsJbsEmb0iJwPiZ5hjFC8wREuiTlhPHDgkBLOiycd20op2nXzDbHfCHInquEe/gYxEitALONxm0swBOwJZwlTDOB7C6y2dzlrtxr1L59m7pCkWI4EtTRLvleehBoj3u7jB4usR
`)
