# Contributing

We love contributions! You are welcome to open a pull request, but it's a good idea to
open an issue and discuss your idea with us first.

Once you are ready to open a PR, please keep the following guidelines in mind:

1. Code should be `go fmt` compliant.
1. Types, structs and funcs should be documented.
1. Tests pass.

## Getting set up

`goApiAbrha` uses go modules. Just fork this repo, clone your fork and off you go!

## Running tests

When working on code in this repository, tests can be run via:

```sh
go test -mod=vendor .
```

## Versioning

GoApiAbrha follows [semver](https://www.semver.org) versioning semantics.
New functionality should be accompanied by increment to the minor
version number. Any code merged to main is subject to release.

## Releasing

Releasing a new version of goApiAbrha is currently a manual process.

Submit a separate pull request for the version change from the pull
request with your changes.

1. Update the `CHANGELOG.md` with your changes. If a version header
   for the next (unreleased) version does not exist, create one.
   Include one bullet point for each piece of new functionality in the
   release, including the pull request ID, description, and author(s).
   For example:

```
## [v1.8.0] - 2019-03-13

- #210 - @jcodybaker - Expose tags on storage volume create/list/get.
- #123 - @parspack - Update test dependencies
```

   To generate a list of changes since the previous release in the correct
   format, you can use [github-changelog-generator](https://github.com/parspack/github-changelog-generator).
   It can be installed from source by running:

```
go get -u github.com/parspack/github-changelog-generator
```

   Next, list the changes by running:

```
github-changelog-generator -org parspack -repo goApiAbrha
```

2. Update the `libraryVersion` number in `goApiAbrha.go`.
3. Make a pull request with these changes.  This PR should be separate from the PR containing the goApiAbrha changes.
4. Once the pull request has been merged, [draft a new release](https://github.com/abrhacom/go-api-abrha/releases/new).
5. Update the `Tag version` and `Release title` field with the new goApiAbrha version.  Be sure the version has a `v` prefixed in both places. Ex `v1.8.0`.
6. Copy the changelog bullet points to the description field.
7. Publish the release.

## Go Version Support

This project follows the support [policy of Go](https://go.dev/doc/devel/release#policy)
as its support policy. The two latest major releases of Go are supported by the project.
[CI workflows](.github/workflows/ci.yml) should test against both supported versions.
[go.mod](./go.mod) should specify the oldest of the supported versions to give
downstream users of goApiAbrha flexibility.
