package spaces_test

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"sort"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	testAccAbrhaSpacesBucketObject_TestRegion = "nyc3"
)

func TestAccAbrhaSpacesBucketObject_noNameNoKey(t *testing.T) {
	bucketError := regexp.MustCompile(`bucket must not be empty`)
	keyError := regexp.MustCompile(`key must not be empty`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaSpacesBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() {},
				Config:      testAccAbrhaSpacesBucketObjectConfigBasic("", "a key"),
				ExpectError: bucketError,
			},
			{
				PreConfig:   func() {},
				Config:      testAccAbrhaSpacesBucketObjectConfigBasic("a name", ""),
				ExpectError: keyError,
			},
		},
	})
}
func TestAccAbrhaSpacesBucketObject_empty(t *testing.T) {
	var obj s3.GetObjectOutput
	resourceName := "abrha_spaces_bucket_object.object"
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaSpacesBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {},
				Config:    testAccAbrhaSpacesBucketObjectConfigEmpty(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketObjectExists(resourceName, &obj),
					testAccCheckAbrhaSpacesBucketObjectBody(&obj, ""),
				),
			},
		},
	})
}

func TestAccAbrhaSpacesBucketObject_source(t *testing.T) {
	var obj s3.GetObjectOutput
	resourceName := "abrha_spaces_bucket_object.object"
	name := acceptance.RandomTestName()

	source := testAccAbrhaSpacesBucketObjectCreateTempFile(t, "{anything will do }")
	defer os.Remove(source)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaSpacesBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaSpacesBucketObjectConfigSource(name, source),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketObjectExists(resourceName, &obj),
					testAccCheckAbrhaSpacesBucketObjectBody(&obj, "{anything will do }"),
				),
			},
		},
	})
}

func TestAccAbrhaSpacesBucketObject_content(t *testing.T) {
	var obj s3.GetObjectOutput
	resourceName := "abrha_spaces_bucket_object.object"
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaSpacesBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {},
				Config:    testAccAbrhaSpacesBucketObjectConfigContent(name, "some_bucket_content"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketObjectExists(resourceName, &obj),
					testAccCheckAbrhaSpacesBucketObjectBody(&obj, "some_bucket_content"),
				),
			},
		},
	})
}

func TestAccAbrhaSpacesBucketObject_contentBase64(t *testing.T) {
	var obj s3.GetObjectOutput
	resourceName := "abrha_spaces_bucket_object.object"
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaSpacesBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {},
				Config:    testAccAbrhaSpacesBucketObjectConfigContentBase64(name, base64.StdEncoding.EncodeToString([]byte("some_bucket_content"))),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketObjectExists(resourceName, &obj),
					testAccCheckAbrhaSpacesBucketObjectBody(&obj, "some_bucket_content"),
				),
			},
		},
	})
}

func TestAccAbrhaSpacesBucketObject_withContentCharacteristics(t *testing.T) {
	var obj s3.GetObjectOutput
	resourceName := "abrha_spaces_bucket_object.object"
	name := acceptance.RandomTestName()

	source := testAccAbrhaSpacesBucketObjectCreateTempFile(t, "{anything will do }")
	defer os.Remove(source)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaSpacesBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaSpacesBucketObjectConfig_withContentCharacteristics(name, source),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketObjectExists(resourceName, &obj),
					testAccCheckAbrhaSpacesBucketObjectBody(&obj, "{anything will do }"),
					resource.TestCheckResourceAttr(resourceName, "content_type", "binary/octet-stream"),
					resource.TestCheckResourceAttr(resourceName, "website_redirect", "http://google.com"),
				),
			},
		},
	})
}

func TestAccAbrhaSpacesBucketObject_NonVersioned(t *testing.T) {
	sourceInitial := testAccAbrhaSpacesBucketObjectCreateTempFile(t, "initial object state")
	defer os.Remove(sourceInitial)

	var originalObj s3.GetObjectOutput
	resourceName := "abrha_spaces_bucket_object.object"
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaSpacesBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaSpacesBucketObjectConfig_NonVersioned(name, sourceInitial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketObjectExists(resourceName, &originalObj),
					testAccCheckAbrhaSpacesBucketObjectBody(&originalObj, "initial object state"),
					resource.TestCheckResourceAttr(resourceName, "version_id", ""),
				),
			},
		},
	})
}

func TestAccAbrhaSpacesBucketObject_updates(t *testing.T) {
	var originalObj, modifiedObj s3.GetObjectOutput
	resourceName := "abrha_spaces_bucket_object.object"
	name := acceptance.RandomTestName()

	sourceInitial := testAccAbrhaSpacesBucketObjectCreateTempFile(t, "initial object state")
	defer os.Remove(sourceInitial)
	sourceModified := testAccAbrhaSpacesBucketObjectCreateTempFile(t, "modified object")
	defer os.Remove(sourceInitial)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaSpacesBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaSpacesBucketObjectConfig_updateable(name, false, sourceInitial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketObjectExists(resourceName, &originalObj),
					testAccCheckAbrhaSpacesBucketObjectBody(&originalObj, "initial object state"),
					resource.TestCheckResourceAttr(resourceName, "etag", "647d1d58e1011c743ec67d5e8af87b53"),
				),
			},
			{
				Config: testAccAbrhaSpacesBucketObjectConfig_updateable(name, false, sourceModified),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketObjectExists(resourceName, &modifiedObj),
					testAccCheckAbrhaSpacesBucketObjectBody(&modifiedObj, "modified object"),
					resource.TestCheckResourceAttr(resourceName, "etag", "1c7fd13df1515c2a13ad9eb068931f09"),
				),
			},
		},
	})
}

func TestAccAbrhaSpacesBucketObject_updateSameFile(t *testing.T) {
	var originalObj, modifiedObj s3.GetObjectOutput
	resourceName := "abrha_spaces_bucket_object.object"
	name := acceptance.RandomTestName()

	startingData := "lane 8"
	changingData := "chicane"

	filename := testAccAbrhaSpacesBucketObjectCreateTempFile(t, startingData)
	defer os.Remove(filename)

	rewriteFile := func(*terraform.State) error {
		if err := os.WriteFile(filename, []byte(changingData), 0644); err != nil {
			os.Remove(filename)
			t.Fatal(err)
		}
		return nil
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaSpacesBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaSpacesBucketObjectConfig_updateable(name, false, filename),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketObjectExists(resourceName, &originalObj),
					testAccCheckAbrhaSpacesBucketObjectBody(&originalObj, startingData),
					resource.TestCheckResourceAttr(resourceName, "etag", "aa48b42f36a2652cbee40c30a5df7d25"),
					rewriteFile,
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccAbrhaSpacesBucketObjectConfig_updateable(name, false, filename),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketObjectExists(resourceName, &modifiedObj),
					testAccCheckAbrhaSpacesBucketObjectBody(&modifiedObj, changingData),
					resource.TestCheckResourceAttr(resourceName, "etag", "fafc05f8c4da0266a99154681ab86e8c"),
				),
			},
		},
	})
}

func TestAccAbrhaSpacesBucketObject_updatesWithVersioning(t *testing.T) {
	var originalObj, modifiedObj s3.GetObjectOutput
	resourceName := "abrha_spaces_bucket_object.object"
	name := acceptance.RandomTestName()

	sourceInitial := testAccAbrhaSpacesBucketObjectCreateTempFile(t, "initial versioned object state")
	defer os.Remove(sourceInitial)
	sourceModified := testAccAbrhaSpacesBucketObjectCreateTempFile(t, "modified versioned object")
	defer os.Remove(sourceInitial)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaSpacesBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaSpacesBucketObjectConfig_updateable(name, true, sourceInitial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketObjectExists(resourceName, &originalObj),
					testAccCheckAbrhaSpacesBucketObjectBody(&originalObj, "initial versioned object state"),
					resource.TestCheckResourceAttr(resourceName, "etag", "cee4407fa91906284e2a5e5e03e86b1b"),
				),
			},
			{
				Config: testAccAbrhaSpacesBucketObjectConfig_updateable(name, true, sourceModified),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketObjectExists(resourceName, &modifiedObj),
					testAccCheckAbrhaSpacesBucketObjectBody(&modifiedObj, "modified versioned object"),
					resource.TestCheckResourceAttr(resourceName, "etag", "00b8c73b1b50e7cc932362c7225b8e29"),
					testAccCheckAbrhaSpacesBucketObjectVersionIdDiffers(&modifiedObj, &originalObj),
				),
			},
		},
	})
}

func TestAccAbrhaSpacesBucketObject_acl(t *testing.T) {
	var obj1, obj2, obj3 s3.GetObjectOutput
	resourceName := "abrha_spaces_bucket_object.object"
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaSpacesBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaSpacesBucketObjectConfig_acl(name, "some_bucket_content", "private"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketObjectExists(resourceName, &obj1),
					testAccCheckAbrhaSpacesBucketObjectBody(&obj1, "some_bucket_content"),
					resource.TestCheckResourceAttr(resourceName, "acl", "private"),
					testAccCheckAbrhaSpacesBucketObjectAcl(resourceName, []string{"FULL_CONTROL"}),
				),
			},
			{
				Config: testAccAbrhaSpacesBucketObjectConfig_acl(name, "some_bucket_content", "public-read"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketObjectExists(resourceName, &obj2),
					testAccCheckAbrhaSpacesBucketObjectVersionIdEquals(&obj2, &obj1),
					testAccCheckAbrhaSpacesBucketObjectBody(&obj2, "some_bucket_content"),
					resource.TestCheckResourceAttr(resourceName, "acl", "public-read"),
					testAccCheckAbrhaSpacesBucketObjectAcl(resourceName, []string{"FULL_CONTROL", "READ"}),
				),
			},
			{
				Config: testAccAbrhaSpacesBucketObjectConfig_acl(name, "changed_some_bucket_content", "private"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketObjectExists(resourceName, &obj3),
					testAccCheckAbrhaSpacesBucketObjectVersionIdDiffers(&obj3, &obj2),
					testAccCheckAbrhaSpacesBucketObjectBody(&obj3, "changed_some_bucket_content"),
					resource.TestCheckResourceAttr(resourceName, "acl", "private"),
					testAccCheckAbrhaSpacesBucketObjectAcl(resourceName, []string{"FULL_CONTROL"}),
				),
			},
		},
	})
}

func TestAccAbrhaSpacesBucketObject_metadata(t *testing.T) {
	name := acceptance.RandomTestName()
	var obj s3.GetObjectOutput
	resourceName := "abrha_spaces_bucket_object.object"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaSpacesBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbrhaSpacesBucketObjectConfig_withMetadata(name, "key1", "value1", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketObjectExists(resourceName, &obj),
					resource.TestCheckResourceAttr(resourceName, "metadata.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key2", "value2"),
				),
			},
			{
				Config: testAccAbrhaSpacesBucketObjectConfig_withMetadata(name, "key1", "value1updated", "key3", "value3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketObjectExists(resourceName, &obj),
					resource.TestCheckResourceAttr(resourceName, "metadata.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key3", "value3"),
				),
			},
			{
				Config: testAccAbrhaSpacesBucketObjectConfigEmpty(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbrhaSpacesBucketObjectExists(resourceName, &obj),
					resource.TestCheckResourceAttr(resourceName, "metadata.%", "0"),
				),
			},
		},
	})
}

func TestAccAbrhaSpacesBucketObject_RegionError(t *testing.T) {
	badRegion := "ny2"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAbrhaBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`resource "abrha_spaces_bucket_object" "object" {
  region = "%s"
  bucket = "foo.parspackspaces.com"
  key    = "test-key"
}`, badRegion),
				ExpectError: regexp.MustCompile(`expected region to be one of`),
			},
		},
	})
}

func testAccGetS3Conn() (*s3.S3, error) {
	client, err := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).SpacesClient(testAccAbrhaSpacesBucketObject_TestRegion)
	if err != nil {
		return nil, err
	}

	s3conn := s3.New(client)

	return s3conn, nil
}

func testAccCheckAbrhaSpacesBucketObjectVersionIdDiffers(first, second *s3.GetObjectOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if first.VersionId == nil {
			return fmt.Errorf("Expected first object to have VersionId: %s", first)
		}
		if second.VersionId == nil {
			return fmt.Errorf("Expected second object to have VersionId: %s", second)
		}

		if *first.VersionId == *second.VersionId {
			return fmt.Errorf("Expected Version IDs to differ, but they are equal (%s)", *first.VersionId)
		}

		return nil
	}
}

func testAccCheckAbrhaSpacesBucketObjectVersionIdEquals(first, second *s3.GetObjectOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if first.VersionId == nil {
			return fmt.Errorf("Expected first object to have VersionId: %s", first)
		}
		if second.VersionId == nil {
			return fmt.Errorf("Expected second object to have VersionId: %s", second)
		}

		if *first.VersionId != *second.VersionId {
			return fmt.Errorf("Expected Version IDs to be equal, but they differ (%s, %s)", *first.VersionId, *second.VersionId)
		}

		return nil
	}
}

func testAccCheckAbrhaSpacesBucketObjectDestroy(s *terraform.State) error {
	s3conn, err := testAccGetS3Conn()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "abrha_spaces_bucket_object":
			_, err := s3conn.HeadObject(
				&s3.HeadObjectInput{
					Bucket:  aws.String(rs.Primary.Attributes["bucket"]),
					Key:     aws.String(rs.Primary.Attributes["key"]),
					IfMatch: aws.String(rs.Primary.Attributes["etag"]),
				})
			if err == nil {
				return fmt.Errorf("Spaces Bucket Object still exists: %s", rs.Primary.ID)
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

func testAccCheckAbrhaSpacesBucketObjectExists(n string, obj *s3.GetObjectOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No S3 Bucket Object ID is set")
		}

		s3conn, err := testAccGetS3Conn()
		if err != nil {
			return err
		}

		out, err := s3conn.GetObject(
			&s3.GetObjectInput{
				Bucket:  aws.String(rs.Primary.Attributes["bucket"]),
				Key:     aws.String(rs.Primary.Attributes["key"]),
				IfMatch: aws.String(rs.Primary.Attributes["etag"]),
			})
		if err != nil {
			return fmt.Errorf("S3Bucket Object error: %s", err)
		}

		*obj = *out

		return nil
	}
}

func testAccCheckAbrhaSpacesBucketObjectBody(obj *s3.GetObjectOutput, want string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		body, err := io.ReadAll(obj.Body)
		if err != nil {
			return fmt.Errorf("failed to read body: %s", err)
		}
		obj.Body.Close()

		if got := string(body); got != want {
			return fmt.Errorf("wrong result body %q; want %q", got, want)
		}

		return nil
	}
}

func testAccCheckAbrhaSpacesBucketObjectAcl(n string, expectedPerms []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]

		s3conn, err := testAccGetS3Conn()
		if err != nil {
			return err
		}

		out, err := s3conn.GetObjectAcl(&s3.GetObjectAclInput{
			Bucket: aws.String(rs.Primary.Attributes["bucket"]),
			Key:    aws.String(rs.Primary.Attributes["key"]),
		})

		if err != nil {
			return fmt.Errorf("GetObjectAcl error: %v", err)
		}

		var perms []string
		for _, v := range out.Grants {
			perms = append(perms, *v.Permission)
		}
		sort.Strings(perms)

		if !reflect.DeepEqual(perms, expectedPerms) {
			return fmt.Errorf("Expected ACL permissions to be %v, got %v", expectedPerms, perms)
		}

		return nil
	}
}

func testAccAbrhaSpacesBucketObjectCreateTempFile(t *testing.T, data string) string {
	tmpFile, err := os.CreateTemp("", "tf-acc-s3-obj")
	if err != nil {
		t.Fatal(err)
	}
	filename := tmpFile.Name()

	err = os.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		os.Remove(filename)
		t.Fatal(err)
	}

	return filename
}

func testAccAbrhaSpacesBucketObjectConfigBasic(bucket, key string) string {
	return fmt.Sprintf(`
resource "abrha_spaces_bucket_object" "object" {
  region = "%s"
  bucket = "%s"
  key    = "%s"
}
`, testAccAbrhaSpacesBucketObject_TestRegion, bucket, key)
}

func testAccAbrhaSpacesBucketObjectConfigEmpty(name string) string {
	return fmt.Sprintf(`
resource "abrha_spaces_bucket" "object_bucket" {
  region        = "%s"
  name          = "%s"
  force_destroy = true
}

resource "abrha_spaces_bucket_object" "object" {
  region = abrha_spaces_bucket.object_bucket.region
  bucket = abrha_spaces_bucket.object_bucket.name
  key    = "test-key"
}
`, testAccAbrhaSpacesBucketObject_TestRegion, name)
}

func testAccAbrhaSpacesBucketObjectConfigSource(name string, source string) string {
	return fmt.Sprintf(`
resource "abrha_spaces_bucket" "object_bucket" {
  region        = "%s"
  name          = "%s"
  force_destroy = true
}

resource "abrha_spaces_bucket_object" "object" {
  region       = abrha_spaces_bucket.object_bucket.region
  bucket       = abrha_spaces_bucket.object_bucket.name
  key          = "test-key"
  source       = "%s"
  content_type = "binary/octet-stream"
}
`, testAccAbrhaSpacesBucketObject_TestRegion, name, source)
}

func testAccAbrhaSpacesBucketObjectConfig_withContentCharacteristics(name string, source string) string {
	return fmt.Sprintf(`
resource "abrha_spaces_bucket" "object_bucket" {
  region        = "%s"
  name          = "%s"
  force_destroy = true
}

resource "abrha_spaces_bucket_object" "object" {
  region           = abrha_spaces_bucket.object_bucket.region
  bucket           = abrha_spaces_bucket.object_bucket.name
  key              = "test-key"
  source           = "%s"
  content_language = "en"
  content_type     = "binary/octet-stream"
  website_redirect = "http://google.com"
}
`, testAccAbrhaSpacesBucketObject_TestRegion, name, source)
}

func testAccAbrhaSpacesBucketObjectConfigContent(name string, content string) string {
	return fmt.Sprintf(`
resource "abrha_spaces_bucket" "object_bucket" {
  region        = "%s"
  name          = "%s"
  force_destroy = true
}

resource "abrha_spaces_bucket_object" "object" {
  region  = abrha_spaces_bucket.object_bucket.region
  bucket  = abrha_spaces_bucket.object_bucket.name
  key     = "test-key"
  content = "%s"
}
`, testAccAbrhaSpacesBucketObject_TestRegion, name, content)
}

func testAccAbrhaSpacesBucketObjectConfigContentBase64(name string, contentBase64 string) string {
	return fmt.Sprintf(`
resource "abrha_spaces_bucket" "object_bucket" {
  region        = "%s"
  name          = "%s"
  force_destroy = true
}

resource "abrha_spaces_bucket_object" "object" {
  region         = abrha_spaces_bucket.object_bucket.region
  bucket         = abrha_spaces_bucket.object_bucket.name
  key            = "test-key"
  content_base64 = "%s"
}
`, testAccAbrhaSpacesBucketObject_TestRegion, name, contentBase64)
}

func testAccAbrhaSpacesBucketObjectConfig_updateable(name string, bucketVersioning bool, source string) string {
	return fmt.Sprintf(`
resource "abrha_spaces_bucket" "object_bucket_3" {
  region        = "%s"
  name          = "%s"
  force_destroy = true

  versioning {
    enabled = %t
  }
}

resource "abrha_spaces_bucket_object" "object" {
  region = abrha_spaces_bucket.object_bucket_3.region
  bucket = abrha_spaces_bucket.object_bucket_3.name
  key    = "updateable-key"
  source = "%s"
  etag   = "${filemd5("%s")}"
}
`, testAccAbrhaSpacesBucketObject_TestRegion, name, bucketVersioning, source, source)
}

func testAccAbrhaSpacesBucketObjectConfig_acl(name string, content, acl string) string {
	return fmt.Sprintf(`
resource "abrha_spaces_bucket" "object_bucket" {
  region        = "%s"
  name          = "%s"
  force_destroy = true

  versioning {
    enabled = true
  }
}

resource "abrha_spaces_bucket_object" "object" {
  region  = abrha_spaces_bucket.object_bucket.region
  bucket  = abrha_spaces_bucket.object_bucket.name
  key     = "test-key"
  content = "%s"
  acl     = "%s"
}
`, testAccAbrhaSpacesBucketObject_TestRegion, name, content, acl)
}

func testAccAbrhaSpacesBucketObjectConfig_withMetadata(name string, metadataKey1, metadataValue1, metadataKey2, metadataValue2 string) string {
	return fmt.Sprintf(`
resource "abrha_spaces_bucket" "object_bucket" {
  region        = "%s"
  name          = "%s"
  force_destroy = true
}

resource "abrha_spaces_bucket_object" "object" {
  region = abrha_spaces_bucket.object_bucket.region
  bucket = abrha_spaces_bucket.object_bucket.name
  key    = "test-key"

  metadata = {
    %[3]s = %[4]q
    %[5]s = %[6]q
  }
}
`, testAccAbrhaSpacesBucketObject_TestRegion, name, metadataKey1, metadataValue1, metadataKey2, metadataValue2)
}

func testAccAbrhaSpacesBucketObjectConfig_NonVersioned(name string, source string) string {
	return fmt.Sprintf(`
resource "abrha_spaces_bucket" "object_bucket_3" {
  region        = "%s"
  name          = "%s"
  force_destroy = true
}

resource "abrha_spaces_bucket_object" "object" {
  region = abrha_spaces_bucket.object_bucket_3.region
  bucket = abrha_spaces_bucket.object_bucket_3.name
  key    = "updateable-key"
  source = "%s"
  etag   = "${filemd5("%s")}"
}
`, testAccAbrhaSpacesBucketObject_TestRegion, name, source, source)
}
