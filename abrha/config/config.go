package config

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"golang.org/x/oauth2"
)

type Config struct {
	Token             string
	APIEndpoint       string
	SpacesAPIEndpoint string
	AccessID          string
	SecretKey         string
	RequestsPerSecond float64
	TerraformVersion  string
	HTTPRetryMax      int
	HTTPRetryWaitMax  float64
	HTTPRetryWaitMin  float64
}

type CombinedConfig struct {
	client                 *goApiAbrha.Client
	spacesEndpointTemplate *template.Template
	accessID               string
	secretKey              string
}

func (c *CombinedConfig) GoApiAbrhaClient() *goApiAbrha.Client { return c.client }

func (c *CombinedConfig) SpacesClient(region string) (*session.Session, error) {
	if c.accessID == "" || c.secretKey == "" {
		err := fmt.Errorf("Spaces credentials not configured")
		return &session.Session{}, err
	}

	endpointWriter := strings.Builder{}
	err := c.spacesEndpointTemplate.Execute(&endpointWriter, map[string]string{
		"Region": strings.ToLower(region),
	})
	if err != nil {
		return &session.Session{}, err
	}
	endpoint := endpointWriter.String()

	client, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(c.accessID, c.secretKey, ""),
		Endpoint:    aws.String(endpoint)},
	)
	if err != nil {
		return &session.Session{}, err
	}

	return client, nil
}

// Client() returns a new client for accessing abrha.
func (c *Config) Client() (*CombinedConfig, error) {
	tokenSrc := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: c.Token,
	})

	userAgent := fmt.Sprintf("Terraform/%s", c.TerraformVersion)
	var client *http.Client
	var goApiAbrhaOpts []goApiAbrha.ClientOpt

	client = oauth2.NewClient(context.Background(), tokenSrc)

	if c.HTTPRetryMax > 0 {
		retryConfig := goApiAbrha.RetryConfig{
			RetryMax:     c.HTTPRetryMax,
			RetryWaitMin: goApiAbrha.PtrTo(c.HTTPRetryWaitMin),
			RetryWaitMax: goApiAbrha.PtrTo(c.HTTPRetryWaitMax),
			Logger:       log.Default(),
		}

		goApiAbrhaOpts = []goApiAbrha.ClientOpt{goApiAbrha.WithRetryAndBackoffs(retryConfig)}
	}

	goApiAbrhaOpts = append(goApiAbrhaOpts, goApiAbrha.SetUserAgent(userAgent))

	headers := map[string]string{
		"Accept-Language": "en",
	}
	goApiAbrhaOpts = append(goApiAbrhaOpts, goApiAbrha.SetRequestHeaders(headers))

	if c.RequestsPerSecond > 0.0 {
		goApiAbrhaOpts = append(goApiAbrhaOpts, goApiAbrha.SetStaticRateLimit(c.RequestsPerSecond))
	}

	goApiAbrhaClient, err := goApiAbrha.New(client, goApiAbrhaOpts...)

	// TODO: logging.NewTransport is deprecated and should be replaced with
	// logging.NewTransportWithRequestLogging.
	//
	//nolint:staticcheck
	clientTransport := logging.NewTransport("Abrha", goApiAbrhaClient.HTTPClient.Transport)

	goApiAbrhaClient.HTTPClient.Transport = clientTransport

	if err != nil {
		return nil, err
	}

	apiURL, err := url.Parse(c.APIEndpoint)
	if err != nil {
		return nil, err
	}
	goApiAbrhaClient.BaseURL = apiURL

	spacesEndpointTemplate, err := template.New("spaces").Parse(c.SpacesAPIEndpoint)
	if err != nil {
		return nil, fmt.Errorf("unable to parse spaces_endpoint '%s' as template: %s", c.SpacesAPIEndpoint, err)
	}

	log.Printf("[INFO] Abrha Client configured for URL: %s", goApiAbrhaClient.BaseURL.String())

	return &CombinedConfig{
		client:                 goApiAbrhaClient,
		spacesEndpointTemplate: spacesEndpointTemplate,
		accessID:               c.AccessID,
		secretKey:              c.SecretKey,
	}, nil
}
