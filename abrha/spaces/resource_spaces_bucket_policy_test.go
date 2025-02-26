package spaces_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/spaces"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	testAccAbrhaSpacesBucketPolicy_TestRegion = "nyc3"
)

func TestAccAbrhaBucketPolicy_basic(t *testing.T) {
	name := acceptance.RandomTestName()

	bucketPolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"s3:*","Resource":"*"}]}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaSpacesBucketPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaSpacesBucketPolicy(name, bucketPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketPolicy("abrha_spaces_bucket_policy.policy", bucketPolicy),
					resource.TestCheckResourceAttr("abrha_spaces_bucket_policy.policy", "region", testAccAbrhaSpacesBucketPolicy_TestRegion),
				),
			},
		},
	})
}

func TestAccAbrhaBucketPolicy_update(t *testing.T) {
	name := acceptance.RandomTestName()

	initialBucketPolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"s3:*","Resource":"*"}]}`
	updatedBucketPolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"s3:*","Resource":"*"},{"Effect":"Allow","Principal":"*","Action":"s3:GetObject","Resource":"*"}]}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaSpacesBucketPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaSpacesBucketPolicy(name, initialBucketPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketPolicy("abrha_spaces_bucket_policy.policy", initialBucketPolicy),
					resource.TestCheckResourceAttr("abrha_spaces_bucket_policy.policy", "region", testAccAbrhaSpacesBucketPolicy_TestRegion),
				),
			},
			{
				Config: testAccAbrhaSpacesBucketPolicy(name, updatedBucketPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketPolicy("abrha_spaces_bucket_policy.policy", updatedBucketPolicy),
					resource.TestCheckResourceAttr("abrha_spaces_bucket_policy.policy", "region", testAccAbrhaSpacesBucketPolicy_TestRegion),
				),
			},
		},
	})
}

func TestAccAbrhaBucketPolicy_invalidJson(t *testing.T) {
	name := acceptance.RandomTestName()

	bucketPolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"s3:*","Resource":"*"}}`
	expectError := regexp.MustCompile(`"policy" contains an invalid JSON`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaSpacesBucketPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAbrhaSpacesBucketPolicy(name, bucketPolicy),
				ExpectError: expectError,
			},
		},
	})
}

func TestAccAbrhaBucketPolicy_emptyPolicy(t *testing.T) {
	name := acceptance.RandomTestName()

	expectError := regexp.MustCompile(`policy must not be empty`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaSpacesBucketPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAbrhaSpacesBucketEmptyPolicy(name),
				ExpectError: expectError,
			},
		},
	})
}

func TestAccAbrhaBucketPolicy_unknownBucket(t *testing.T) {
	expectError := regexp.MustCompile(`bucket 'unknown' does not exist`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaSpacesBucketPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAbrhaSpacesBucketUnknownBucket(),
				ExpectError: expectError,
			},
		},
	})
}

func testAccGetS3PolicyConn() (*s3.S3, error) {
	client, err := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).SpacesClient(testAccAbrhaSpacesBucketPolicy_TestRegion)
	if err != nil {
		return nil, err
	}

	s3conn := s3.New(client)

	return s3conn, nil
}

func testAccCheckAbrhaSpacesBucketPolicy(n string, expectedPolicy string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No S3 Bucket Policy ID is set")
		}

		s3conn, err := testAccGetS3PolicyConn()
		if err != nil {
			return err
		}

		response, err := s3conn.GetBucketPolicy(
			&s3.GetBucketPolicyInput{
				Bucket: aws.String(rs.Primary.Attributes["bucket"]),
			})
		if err != nil {
			return fmt.Errorf("S3Bucket policy error: %s", err)
		}

		actualPolicy := aws.StringValue(response.Policy)
		equivalent := spaces.CompareSpacesBucketPolicy(expectedPolicy, actualPolicy)
		if !equivalent {
			return fmt.Errorf("Expected policy to be '%v', got '%v'", expectedPolicy, actualPolicy)
		}
		return nil
	}
}

func testAccCheckAbrhaSpacesBucketPolicyDestroy(s *terraform.State) error {
	s3conn, err := testAccGetS3PolicyConn()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "abrha_spaces_bucket_policy":
			_, err := s3conn.GetBucketPolicy(&s3.GetBucketPolicyInput{
				Bucket: aws.String(rs.Primary.Attributes["bucket"]),
			})
			if err == nil {
				return fmt.Errorf("Spaces Bucket policy still exists: %s", rs.Primary.ID)
			}

		case "abrha_spaces_bucket":
			_, err = s3conn.HeadBucket(&s3.HeadBucketInput{
				Bucket: aws.String(rs.Primary.ID),
			})
			if err == nil {
				return fmt.Errorf("Spaces Bucket still exists: %s", rs.Primary.ID)
			}

		default:
			continue
		}
	}

	return nil
}

func testAccAbrhaSpacesBucketPolicy(name string, policy string) string {
	return fmt.Sprintf(`
resource "abrha_spaces_bucket" "policy_bucket" {
  region        = "%s"
  name          = "%s"
  force_destroy = true
}

resource "abrha_spaces_bucket_policy" "policy" {
  region = abrha_spaces_bucket.policy_bucket.region
  bucket = abrha_spaces_bucket.policy_bucket.name
  policy = <<EOF
%s
EOF
}


`, testAccAbrhaSpacesBucketPolicy_TestRegion, name, policy)
}

func testAccAbrhaSpacesBucketEmptyPolicy(name string) string {
	return fmt.Sprintf(`
resource "abrha_spaces_bucket" "policy_bucket" {
  region        = "%s"
  name          = "%s"
  force_destroy = true
}

resource "abrha_spaces_bucket_policy" "policy" {
  region = abrha_spaces_bucket.policy_bucket.region
  bucket = abrha_spaces_bucket.policy_bucket.name
  policy = ""
}


`, testAccAbrhaSpacesBucketPolicy_TestRegion, name)
}

func testAccAbrhaSpacesBucketUnknownBucket() string {
	return fmt.Sprintf(`
resource "abrha_spaces_bucket_policy" "policy" {
  region = "%s"
  bucket = "unknown"
  policy = "{}"
}

`, testAccAbrhaSpacesBucketPolicy_TestRegion)
}
