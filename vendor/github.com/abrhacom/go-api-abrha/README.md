# go-api-abrha

[![GitHub Actions CI](https://github.com/abrhacom/go-api-abrha/actions/workflows/ci.yml/badge.svg)](https://github.com/abrhacom/go-api-abrha/actions/workflows/ci.yml)
[![GoDoc](https://godoc.org/github.com/abrhacom/go-api-abrha?status.svg)](https://godoc.org/github.com/abrhacom/go-api-abrha)

go-api-abrha is a Go client library for accessing the Abrha V1 API.

You can view the client API docs here: [http://godoc.org/github.com/abrhacom/go-api-abrha](http://godoc.org/github.com/abrhacom/go-api-abrha)

You can view Abrha API docs here: [https://docs.parspack.com/api/](https://docs.parspack.com/api/)

## Install
```sh
go get github.com/abrhacom/go-api-abrha@vX.Y.Z
```

where X.Y.Z is the [version](https://github.com/abrhacom/go-api-abrha/releases) you need.

or
```sh
go get github.com/abrhacom/go-api-abrha
```
for non Go modules usage or latest version.

## Usage

```go
import "github.com/abrhacom/go-api-abrha"
```

Create a new Abrha client, then use the exposed services to
access different parts of the Abrha API.

### Authentication

Currently, Personal Access Token (PAT) is the only method of
authenticating with the API. You can manage your tokens
at the Abrha Control Panel.

You can then use your token to create a new client:

```go
package main

import (
    goApiAbrha "github.com/abrhacom/go-api-abrha"
)

func main() {
    client := goApiAbrha.NewFromToken("my-abrha-api-token")
}
```

If you need to provide a `context.Context` to your new client, you should use [`goApiAbrha.NewClient`](https://godoc.org/github.com/abrhacom/go-api-abrha#NewClient#NewClient) to manually construct a client instead.

## Examples


To create a new Vm:

```go
vmName := "super-cool-vm"

createRequest := &goApiAbrha.VmCreateRequest{
    Name:   vmName,
    Region: "nyc3",
    Size:   "s-1vcpu-1gb",
    Image: goApiAbrha.VmCreateImage{
        Slug: "ubuntu-20-04-x64",
    },
}

ctx := context.TODO()

newVm, _, err := client.Vms.Create(ctx, createRequest)

if err != nil {
    fmt.Printf("Something bad happened: %s\n\n", err)
    return err
}
```

### Pagination

If a list of items is paginated by the API, you must request pages individually. For example, to fetch all Vms:

```go
func VmList(ctx context.Context, client *goApiAbrha.Client) ([]goApiAbrha.Vm, error) {
    // create a list to hold our vms
    list := []goApiAbrha.Vm{}

    // create options. initially, these will be blank
    opt := &goApiAbrha.ListOptions{}
    for {
        vms, resp, err := client.Vms.List(ctx, opt)
        if err != nil {
            return nil, err
        }

        // append the current page's vms to our list
        list = append(list, vms...)

        // if we are at the last page, break out the for loop
        if resp.Links == nil || resp.Links.IsLastPage() {
            break
        }

        page, err := resp.Links.CurrentPage()
        if err != nil {
            return nil, err
        }

        // set the page we want for the next request
        opt.Page = page + 1
    }

    return list, nil
}
```

Some endpoints offer token based pagination. For example, to fetch all Registry Repositories:

```go
func ListRepositoriesV2(ctx context.Context, client *goApiAbrha.Client, registryName string) ([]*goApiAbrha.RepositoryV2, error) {
    // create a list to hold our registries
    list := []*goApiAbrha.RepositoryV2{}

    // create options. initially, these will be blank
    opt := &goApiAbrha.TokenListOptions{}
    for {
        repositories, resp, err := client.Registry.ListRepositoriesV2(ctx, registryName, opt)
        if err != nil {
            return nil, err
        }

        // append the current page's registries to our list
        list = append(list, repositories...)

        // if we are at the last page, break out the for loop
        if resp.Links == nil || resp.Links.IsLastPage() {
            break
        }

        // grab the next page token
        nextPageToken, err := resp.Links.NextPageToken()
        if err != nil {
            return nil, err
        }

        // provide the next page token for the next request
        opt.Token = nextPageToken
    }

    return list, nil
}
```

### Automatic Retries and Exponential Backoff

The go-api-abrha client can be configured to use automatic retries and exponentional backoff for requests that fail with 429 or 500-level response codes via [go-retryablehttp](https://github.com/hashicorp/go-retryablehttp). To configure go-api-abrha to enable usage of go-retryablehttp, the `RetryConfig.RetryMax` must be set.

```go
tokenSrc := oauth2.StaticTokenSource(&oauth2.Token{
    AccessToken: "dop_v1_xxxxxx",
})

oauth_client := oauth2.NewClient(oauth2.NoContext, tokenSrc)

waitMax := goApiAbrha.PtrTo(6.0)
waitMin := goApiAbrha.PtrTo(3.0)

retryConfig := goApiAbrha.RetryConfig{
    RetryMax:     3,
    RetryWaitMin: waitMin,
    RetryWaitMax: waitMax,
}

client, err := goApiAbrha.New(oauth_client, goApiAbrha.WithRetryAndBackoffs(retryConfig))
```

Please refer to the [RetryConfig go-api-abrha documentation](https://pkg.go.dev/github.com/abrhacom/go-api-abrha#RetryConfig) for more information.

## Versioning

Each version of the client is tagged and the version is updated accordingly.

To see the list of past versions, run `git tag`.


## Documentation

For a comprehensive list of examples, check out the [API documentation](https://docs.parspack.com/api/#tag/SSH-Keys).

For details on all the functionality in this library, see the [GoDoc](https://godoc.org/github.com/abrhacom/go-api-abrha) documentation.


## Contributing

We love pull requests! Please see the [contribution guidelines](CONTRIBUTING.md).
