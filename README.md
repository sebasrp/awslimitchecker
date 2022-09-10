

# awslimitchecker

[![codecov](https://codecov.io/gh/sebasrp/awslimitchecker/branch/main/graph/badge.svg?token=Y5AOHU08FU)](https://codecov.io/gh/sebasrp/awslimitchecker)

Simple module to programatically retrieve your AWS account limits (whether they are supporter by servicequotas or not). It also provide a lightweight cli program to access the data.

## cli

A utility `awslimitchecker` CLI is provided, that exposes the module through a simple interface.

## Usage

Make sure you are logged into your AWS account (`aws configure` or through environment variables). This account needs to have the required IAM permissions.

Check the help page with `awslimitchecker --help` to see all available commands.

### List required permissions

`awslimitchecker` requires a set of permissions in order to retrieve usage and quota information. To list the required AWS IAM policies, use the `iam` command line argument

```shell
awslimitchecker iam
```

### Run a check on a single service

```shell
awslimitchecker check s3
```

### Run all the available checks

```shell
awslimitchecker check all
```

### Export data to csv

```shell
awslimitchecker check all --csv
```

## Development

To run the latest:

```shell
cd awslimitchecker
go build ./... && go install ./...
awslimitechchecker --help
```

When making changes:

1. make sure you add relevant tests (there is a github action doing codecov validation)
2. make sure the existing tests pass `go test ./...` from root directory
