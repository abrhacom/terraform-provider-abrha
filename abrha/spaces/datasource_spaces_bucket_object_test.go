package spaces_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAbrhaSpacesBucketObject_basic(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceOnlyConf, conf := testAccDataSourceAbrhaSpacesObjectConfig_basic(name)

	var rObj s3.GetObjectOutput
	var dsObj s3.GetObjectOutput

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories:         acceptance.TestAccProviderFactories,
		CheckDestroy:              testAccCheckAbrhaBucketDestroy,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: resourceOnlyConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketObjectExists("abrha_spaces_bucket_object.object", &rObj),
				),
			},
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesObjectDataSourceExists("data.abrha_spaces_bucket_object.obj", &dsObj),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_object.obj", "content_length", "11"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_object.obj", "content_type", "binary/octet-stream"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_object.obj", "etag", "b10a8db164e0754105b7a99be72e3fe5"),
					resource.TestMatchResourceAttr("data.abrha_spaces_bucket_object.obj", "last_modified",
						regexp.MustCompile("^[a-zA-Z]{3}, [0-9]+ [a-zA-Z]+ [0-9]{4} [0-9:]+ [A-Z]+$")),
					resource.TestCheckNoResourceAttr("data.abrha_spaces_bucket_object.obj", "body"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaSpacesBucketObject_readableBody(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceOnlyConf, conf := testAccDataSourceAbrhaSpacesObjectConfig_readableBody(name)

	var rObj s3.GetObjectOutput
	var dsObj s3.GetObjectOutput

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories:         acceptance.TestAccProviderFactories,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: resourceOnlyConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketObjectExists("abrha_spaces_bucket_object.object", &rObj),
				),
			},
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesObjectDataSourceExists("data.abrha_spaces_bucket_object.obj", &dsObj),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_object.obj", "content_length", "3"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_object.obj", "content_type", "text/plain"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_object.obj", "etag", "a6105c0a611b41b08f1209506350279e"),
					resource.TestMatchResourceAttr("data.abrha_spaces_bucket_object.obj", "last_modified",
						regexp.MustCompile("^[a-zA-Z]{3}, [0-9]+ [a-zA-Z]+ [0-9]{4} [0-9:]+ [A-Z]+$")),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_object.obj", "body", "yes"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaSpacesBucketObject_allParams(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceOnlyConf, conf := testAccDataSourceAbrhaSpacesObjectConfig_allParams(name)

	var rObj s3.GetObjectOutput
	var dsObj s3.GetObjectOutput

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories:         acceptance.TestAccProviderFactories,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: resourceOnlyConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketObjectExists("abrha_spaces_bucket_object.object", &rObj),
				),
			},
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesObjectDataSourceExists("data.abrha_spaces_bucket_object.obj", &dsObj),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_object.obj", "content_length", "21"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_object.obj", "content_type", "application/unknown"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_object.obj", "etag", "723f7a6ac0c57b445790914668f98640"),
					resource.TestMatchResourceAttr("data.abrha_spaces_bucket_object.obj", "last_modified",
						regexp.MustCompile("^[a-zA-Z]{3}, [0-9]+ [a-zA-Z]+ [0-9]{4} [0-9:]+ [A-Z]+$")),
					resource.TestCheckResourceAttrSet("data.abrha_spaces_bucket_object.obj", "version_id"),
					resource.TestCheckNoResourceAttr("data.abrha_spaces_bucket_object.obj", "body"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_object.obj", "cache_control", "no-cache"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_object.obj", "content_disposition", "attachment"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_object.obj", "content_encoding", "identity"),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_object.obj", "content_language", "en-GB"),
					// Encryption is off
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_object.obj", "expiration", ""),
					// Currently unsupported in abrha_spaces_bucket_object resource
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_object.obj", "expires", ""),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_object.obj", "website_redirect_location", ""),
					resource.TestCheckResourceAttr("data.abrha_spaces_bucket_object.obj", "metadata.%", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaSpacesBucketObject_LeadingSlash(t *testing.T) {
	var rObj s3.GetObjectOutput
	var dsObj1, dsObj2, dsObj3 s3.GetObjectOutput
	resourceName := "abrha_spaces_bucket_object.object"
	dataSourceName1 := "data.abrha_spaces_bucket_object.obj1"
	dataSourceName2 := "data.abrha_spaces_bucket_object.obj2"
	dataSourceName3 := "data.abrha_spaces_bucket_object.obj3"
	name := acceptance.RandomTestName()
	resourceOnlyConf, conf := testAccDataSourceAbrhaSpacesObjectConfig_leadingSlash(name)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories:         acceptance.TestAccProviderFactories,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: resourceOnlyConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketObjectExists(resourceName, &rObj),
				),
			},
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesObjectDataSourceExists(dataSourceName1, &dsObj1),
					resource.TestCheckResourceAttr(dataSourceName1, "content_length", "3"),
					resource.TestCheckResourceAttr(dataSourceName1, "content_type", "text/plain"),
					resource.TestCheckResourceAttr(dataSourceName1, "etag", "a6105c0a611b41b08f1209506350279e"),
					resource.TestMatchResourceAttr(dataSourceName1, "last_modified",
						regexp.MustCompile("^[a-zA-Z]{3}, [0-9]+ [a-zA-Z]+ [0-9]{4} [0-9:]+ [A-Z]+$")),
					resource.TestCheckResourceAttr(dataSourceName1, "body", "yes"),
					testAccCheckAbrhaSpacesObjectDataSourceExists(dataSourceName2, &dsObj2),
					resource.TestCheckResourceAttr(dataSourceName2, "content_length", "3"),
					resource.TestCheckResourceAttr(dataSourceName2, "content_type", "text/plain"),
					resource.TestCheckResourceAttr(dataSourceName2, "etag", "a6105c0a611b41b08f1209506350279e"),
					resource.TestMatchResourceAttr(dataSourceName2, "last_modified",
						regexp.MustCompile("^[a-zA-Z]{3}, [0-9]+ [a-zA-Z]+ [0-9]{4} [0-9:]+ [A-Z]+$")),
					resource.TestCheckResourceAttr(dataSourceName2, "body", "yes"),
					testAccCheckAbrhaSpacesObjectDataSourceExists(dataSourceName3, &dsObj3),
					resource.TestCheckResourceAttr(dataSourceName3, "content_length", "3"),
					resource.TestCheckResourceAttr(dataSourceName3, "content_type", "text/plain"),
					resource.TestCheckResourceAttr(dataSourceName3, "etag", "a6105c0a611b41b08f1209506350279e"),
					resource.TestMatchResourceAttr(dataSourceName3, "last_modified",
						regexp.MustCompile("^[a-zA-Z]{3}, [0-9]+ [a-zA-Z]+ [0-9]{4} [0-9:]+ [A-Z]+$")),
					resource.TestCheckResourceAttr(dataSourceName3, "body", "yes"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaSpacesBucketObject_MultipleSlashes(t *testing.T) {
	var rObj1, rObj2 s3.GetObjectOutput
	var dsObj1, dsObj2, dsObj3 s3.GetObjectOutput
	resourceName1 := "abrha_spaces_bucket_object.object1"
	resourceName2 := "abrha_spaces_bucket_object.object2"
	dataSourceName1 := "data.abrha_spaces_bucket_object.obj1"
	dataSourceName2 := "data.abrha_spaces_bucket_object.obj2"
	dataSourceName3 := "data.abrha_spaces_bucket_object.obj3"
	name := acceptance.RandomTestName()
	resourceOnlyConf, conf := testAccDataSourceAbrhaSpacesObjectConfig_multipleSlashes(name)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories:         acceptance.TestAccProviderFactories,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: resourceOnlyConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketObjectExists(resourceName1, &rObj1),
					testAccCheckAbrhaSpacesBucketObjectExists(resourceName2, &rObj2),
				),
			},
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesObjectDataSourceExists(dataSourceName1, &dsObj1),
					resource.TestCheckResourceAttr(dataSourceName1, "content_length", "3"),
					resource.TestCheckResourceAttr(dataSourceName1, "content_type", "text/plain"),
					resource.TestCheckResourceAttr(dataSourceName1, "body", "yes"),
					testAccCheckAbrhaSpacesObjectDataSourceExists(dataSourceName2, &dsObj2),
					resource.TestCheckResourceAttr(dataSourceName2, "content_length", "3"),
					resource.TestCheckResourceAttr(dataSourceName2, "content_type", "text/plain"),
					resource.TestCheckResourceAttr(dataSourceName2, "body", "yes"),
					testAccCheckAbrhaSpacesObjectDataSourceExists(dataSourceName3, &dsObj3),
					resource.TestCheckResourceAttr(dataSourceName3, "content_length", "2"),
					resource.TestCheckResourceAttr(dataSourceName3, "content_type", "text/plain"),
					resource.TestCheckResourceAttr(dataSourceName3, "body", "no"),
				),
			},
		},
	})
}

func TestAccDataSourceAbrhaSpacesBucketObject_RegionError(t *testing.T) {
	badRegion := "ny2"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`data "abrha_spaces_bucket_object" "object" {
  region = "%s"
  bucket = "foo.parspackspaces.com"
  key    = "test-key"
}`, badRegion),
				ExpectError: regexp.MustCompile(`expected region to be one of`),
			},
		},
	})
}

func testAccCheckAbrhaSpacesObjectDataSourceExists(n string, obj *s3.GetObjectOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find S3 object data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("S3 object data source ID not set")
		}

		s3conn, err := testAccGetS3ConnForSpacesBucket(rs)
		if err != nil {
			return err
		}

		out, err := s3conn.GetObject(
			&s3.GetObjectInput{
				Bucket: aws.String(rs.Primary.Attributes["bucket"]),
				Key:    aws.String(rs.Primary.Attributes["key"]),
			})
		if err != nil {
			return fmt.Errorf("Failed getting S3 Object from %s: %s",
				rs.Primary.Attributes["bucket"]+"/"+rs.Primary.Attributes["key"], err)
		}

		*obj = *out

		return nil
	}
}

func testAccDataSourceAbrhaSpacesObjectConfig_basic(name string) (string, string) {
	resources := fmt.Sprintf(`
resource "abrha_spaces_bucket" "object_bucket" {
  name   = "%s"
  region = "nyc3"
}
resource "abrha_spaces_bucket_object" "object" {
  bucket  = abrha_spaces_bucket.object_bucket.name
  region  = abrha_spaces_bucket.object_bucket.region
  key     = "%s-object"
  content = "Hello World"
}
`, name, name)

	both := fmt.Sprintf(`%s
data "abrha_spaces_bucket_object" "obj" {
  bucket = "%s"
  region = "nyc3"
  key    = "%s-object"
}
`, resources, name, name)

	return resources, both
}

func testAccDataSourceAbrhaSpacesObjectConfig_readableBody(name string) (string, string) {
	resources := fmt.Sprintf(`
resource "abrha_spaces_bucket" "object_bucket" {
  name   = "%s"
  region = "nyc3"
}
resource "abrha_spaces_bucket_object" "object" {
  bucket       = abrha_spaces_bucket.object_bucket.name
  region       = abrha_spaces_bucket.object_bucket.region
  key          = "%s-readable"
  content      = "yes"
  content_type = "text/plain"
}
`, name, name)

	both := fmt.Sprintf(`%s
data "abrha_spaces_bucket_object" "obj" {
  bucket = "%s"
  region = "nyc3"
  key    = "%s-readable"
}
`, resources, name, name)

	return resources, both
}

func testAccDataSourceAbrhaSpacesObjectConfig_allParams(name string) (string, string) {
	resources := fmt.Sprintf(`
resource "abrha_spaces_bucket" "object_bucket" {
  name   = "%s"
  region = "nyc3"
  versioning {
    enabled = true
  }
}

resource "abrha_spaces_bucket_object" "object" {
  bucket              = abrha_spaces_bucket.object_bucket.name
  region              = abrha_spaces_bucket.object_bucket.region
  key                 = "%s-all-params"
  content             = <<CONTENT
{"msg": "Hi there!"}
CONTENT
  content_type        = "application/unknown"
  cache_control       = "no-cache"
  content_disposition = "attachment"
  content_encoding    = "identity"
  content_language    = "en-GB"
}
`, name, name)

	both := fmt.Sprintf(`%s
data "abrha_spaces_bucket_object" "obj" {
  bucket = "%s"
  region = "nyc3"
  key    = "%s-all-params"
}
`, resources, name, name)

	return resources, both
}

func testAccDataSourceAbrhaSpacesObjectConfig_leadingSlash(name string) (string, string) {
	resources := fmt.Sprintf(`
resource "abrha_spaces_bucket" "object_bucket" {
  name   = "%s"
  region = "nyc3"
}
resource "abrha_spaces_bucket_object" "object" {
  bucket       = abrha_spaces_bucket.object_bucket.name
  region       = abrha_spaces_bucket.object_bucket.region
  key          = "//%s-readable"
  content      = "yes"
  content_type = "text/plain"
}
`, name, name)

	both := fmt.Sprintf(`%s
data "abrha_spaces_bucket_object" "obj1" {
  bucket = "%s"
  region = "nyc3"
  key    = "%s-readable"
}
data "abrha_spaces_bucket_object" "obj2" {
  bucket = "%s"
  region = "nyc3"
  key    = "/%s-readable"
}
data "abrha_spaces_bucket_object" "obj3" {
  bucket = "%s"
  region = "nyc3"
  key    = "//%s-readable"
}
`, resources, name, name, name, name, name, name)

	return resources, both
}

func testAccDataSourceAbrhaSpacesObjectConfig_multipleSlashes(name string) (string, string) {
	resources := fmt.Sprintf(`
resource "abrha_spaces_bucket" "object_bucket" {
  name   = "%s"
  region = "nyc3"
}
resource "abrha_spaces_bucket_object" "object1" {
  bucket       = abrha_spaces_bucket.object_bucket.name
  region       = abrha_spaces_bucket.object_bucket.region
  key          = "first//second///third//"
  content      = "yes"
  content_type = "text/plain"
}
# Without a trailing slash.
resource "abrha_spaces_bucket_object" "object2" {
  bucket       = abrha_spaces_bucket.object_bucket.name
  region       = abrha_spaces_bucket.object_bucket.region
  key          = "/first////second/third"
  content      = "no"
  content_type = "text/plain"
}
`, name)

	both := fmt.Sprintf(`%s
data "abrha_spaces_bucket_object" "obj1" {
  bucket = "%s"
  region = "nyc3"
  key    = "first/second/third/"
}
data "abrha_spaces_bucket_object" "obj2" {
  bucket = "%s"
  region = "nyc3"
  key    = "first//second///third//"
}
data "abrha_spaces_bucket_object" "obj3" {
  bucket = "%s"
  region = "nyc3"
  key    = "first/second/third"
}
`, resources, name, name, name)

	return resources, both
}
