package tfe

import (
	"fmt"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTFEPolicySetParameter_basic(t *testing.T) {
	parameter := &tfe.PolicySetParameter{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTFEPolicySetParameterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTFEPolicySetParameter_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTFEPolicySetParameterExists(
						"tfe_policy_set_parameter.foobar", parameter),
					testAccCheckTFEPolicySetParameterAttributes(parameter),
					resource.TestCheckResourceAttr(
						"tfe_policy_set_parameter.foobar", "key", "key_test"),
					resource.TestCheckResourceAttr(
						"tfe_policy_set_parameter.foobar", "value", "value_test"),
					resource.TestCheckResourceAttr(
						"tfe_policy_set_parameter.foobar", "sensitive", "false"),
				),
			},
		},
	})
}

func TestAccTFEPolicySetParameter_update(t *testing.T) {
	parameter := &tfe.PolicySetParameter{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTFEPolicySetParameterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTFEPolicySetParameter_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTFEPolicySetParameterExists(
						"tfe_policy_set_parameter.foobar", parameter),
					testAccCheckTFEPolicySetParameterAttributes(parameter),
					resource.TestCheckResourceAttr(
						"tfe_policy_set_parameter.foobar", "key", "key_test"),
					resource.TestCheckResourceAttr(
						"tfe_policy_set_parameter.foobar", "value", "value_test"),
					resource.TestCheckResourceAttr(
						"tfe_policy_set_parameter.foobar", "sensitive", "false"),
				),
			},

			{
				Config: testAccTFEPolicySetParameter_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTFEPolicySetParameterExists(
						"tfe_policy_set_parameter.foobar", parameter),
					testAccCheckTFEPolicySetParameterAttributesUpdate(parameter),
					resource.TestCheckResourceAttr(
						"tfe_policy_set_parameter.foobar", "key", "key_updated"),
					resource.TestCheckResourceAttr(
						"tfe_policy_set_parameter.foobar", "value", "value_updated"),
					resource.TestCheckResourceAttr(
						"tfe_policy_set_parameter.foobar", "sensitive", "true"),
				),
			},
		},
	})
}

func TestAccTFEPolicySetParameter_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTFEPolicySetParameterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTFEPolicySetParameter_basic,
			},

			{
				ResourceName:        "tfe_policy_set_parameter.foobar",
				ImportState:         true,
				ImportStateIdPrefix: "policy-set-test/",
				ImportStateVerify:   true,
			},
		},
	})
}

func testAccCheckTFEPolicySetParameterExists(
	n string, parameter *tfe.PolicySetParameter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		tfeClient := testAccProvider.Meta().(*tfe.Client)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		v, err := tfeClient.PolicySetParameters.Read(ctx, rs.Primary.Attributes["policy_set_id"], rs.Primary.ID)
		if err != nil {
			return err
		}

		*parameter = *v

		return nil
	}
}

func testAccCheckTFEPolicySetParameterAttributes(
	parameter *tfe.PolicySetParameter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if parameter.Key != "key_test" {
			return fmt.Errorf("Bad key: %s", parameter.Key)
		}

		if parameter.Value != "value_test" {
			return fmt.Errorf("Bad value: %s", parameter.Value)
		}

		if parameter.Sensitive != false {
			return fmt.Errorf("Bad sensitive: %t", parameter.Sensitive)
		}

		return nil
	}
}

func testAccCheckTFEPolicySetParameterAttributesUpdate(
	parameter *tfe.PolicySetParameter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if parameter.Key != "key_updated" {
			return fmt.Errorf("Bad key: %s", parameter.Key)
		}

		if parameter.Value != "" {
			return fmt.Errorf("Bad value: %s", parameter.Value)
		}

		if parameter.Sensitive != true {
			return fmt.Errorf("Bad sensitive: %t", parameter.Sensitive)
		}

		return nil
	}
}

func testAccCheckTFEPolicySetParameterDestroy(s *terraform.State) error {
	tfeClient := testAccProvider.Meta().(*tfe.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tfe_policy_set_parameter" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		_, err := tfeClient.PolicySetParameters.Read(ctx, rs.Primary.Attributes["policy_set_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("PolicySetParameter %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccTFEPolicySetParameter_basic = `
resource "tfe_organization" "foobar" {
  name  = "tst-terraform"
  email = "admin@company.com"
}

resource "tfe_policy_set" "foobar" {
  name         = "policy-set-test"
  organization = "${tfe_organization.foobar.id}"
}

resource "tfe_policy_set_parameter" "foobar" {
  key          = "key_test"
  value        = "value_test"
  policy_set_id = "${tfe_policy_set.foobar.id}"
}`

const testAccTFEPolicySetParameter_update = `
resource "tfe_organization" "foobar" {
  name  = "tst-terraform"
  email = "admin@company.com"
}

resource "tfe_policy_set" "foobar" {
  name         = "policy-set-test"
  organization = "${tfe_organization.foobar.id}"
}

resource "tfe_policy_set_parameter" "foobar" {
  key          = "key_updated"
  value        = "value_updated"
  sensitive    = true
  policy_set_id = "${tfe_policy_set.foobar.id}"
}`
