
# awslimitchecker

[![codecov](https://codecov.io/gh/sebasrp/awslimitchecker/branch/main/graph/badge.svg?token=Y5AOHU08FU)](https://codecov.io/gh/sebasrp/awslimitchecker)
[![license](https://img.shields.io/github/license/sebasrp/awslimitchecker)](https://tldrlegal.com/license/mit-license)
[![CI](https://github.com/sebasrp/awslimitchecker/actions/workflows/workflow.yml/badge.svg)](https://github.com/sebasrp/awslimitchecker/actions/workflows/workflow.yml)
[![go Report Card](https://goreportcard.com/badge/github.com/sebasrp/awslimitchecker)](https://goreportcard.com/report/github.com/sebasrp/awslimitchecker)

Simple module to programatically retrieve your AWS account limits (whether they are supporter by servicequotas or not). It also provide a lightweight cli program to access the data.

## Status

The project is under active development. We are focusing primarly to put the basic foundations of the module/cli in order to make it useful.
Not many services are currently supported, but it's fairly simple to add them - priority at the moment is to put the foundations - adding services is done gradually.

## Features

* Check current AWS resource usage against AWS ServiceQuota limits (see [wiki - Supported Quotas](https://github.com/sebasrp/awslimitchecker/wiki/Supported-Quotas) for complete list)
* Retrieves current usage
* Compare current usage to limits
* When available, retrieves applied (different than default) values
* Supports explicitely setting the AWS region

## cli

A utility `awslimitchecker` CLI is provided, that exposes the module through a simple interface.

### Usage

Make sure you are logged into your AWS account (`aws configure` or through environment variables). This account needs to have the required IAM permissions.

Check the help page with `awslimitchecker --help` to see all available commands.

### List required permissions

`awslimitchecker` requires a set of permissions in order to retrieve usage and quota information. To list the required AWS IAM policies, use the `iam` command line argument

```shell
➜ awslimitchecker iam
Required IAM permissions to retrieve usage/limits:
* dynamodb:ListTables
* eks:ListClusters
* eks:ListNodegroups
* elasticache:DescribeCacheClusters
* elasticloadbalancing:DescribeLoadBalancers
* elasticloadbalancing:DescribeAccountLimits
* iam:GetAccountSummary
* kinesis:DescribeLimits
* rds:DescribeAccountAttributes
* s3:ListAllMyBuckets
* sns:ListTopics
* sns:ListSubscriptions
```

### Run a check on a single service

(note - all "usage" have been manufactured/are examples)

```shell
➜ awslimitchecker check rds --console
AWS profile: default | AWS region: ap-southeast-1 | service: rds
* [rds] DB instances  100/600
* [rds] DB clusters  100/300
* [rds] Reserved DB instances  0/600
```

### Run all the available checks

(note - all "actuals" have been manufactured/are examples)

```shell
➜ awslimitchecker check all
AWS profile: default | AWS region: ap-southeast-1 | service: all
* [rds] DB instances  100/600
* [rds] DB clusters  100/300
* [rds] Reserved DB instances  0/600
* [dynamodb] Maximum number of tables  100/2500
* [eks] Clusters  1/100
* [eks] Managed node groups per cluster (AWS::EKS::Cluster::cluster1) 0/30
* [eks] Managed node groups per cluster (AWS::EKS::Cluster::cluster2) 0/30
* [elasticache] Nodes per Region  10/300
* [s3] Buckets  20/100
* [sns] Topics per Account  300/100000
* [sns] Pending Subscriptions per Account  300/5000
* [elasticloadbalancing] Classic Load Balancers per Region  12/100
* [elasticloadbalancing] Application Load Balancers per Region  12/100
* [elasticloadbalancing] Network Load Balancers per Region  12/50
* [iam] Instance profiles per Account  100/1000
* [iam] Policies per Account  1000/3000
* [iam] Server Certificates per Account  10/25
* [iam] Roles per Account  1000/5000
* [iam] Users per Account  100/5000
* [iam] Groups per Account  100/300
* [kinesis] On-demand Data Streams per account  10/50
* [kinesis] Shards per Region  10/200
```

### Export data to csv

```shell
awslimitchecker check all --csv
```

### Configuration file

Tired of manually selecting the different parameters? You can save those in a file and provide it with the `--config flag` - or just place it under `$HOME/.awslimitchecker` to be automatically picked up. The format and options supported are (order does not matter)

```yaml
awsprofile: <name of profile>
region: <region to evaluate>
console: true /false
csv: true / false
verbose: true / false
```

## Development

To run the latest:

```shell
cd awslimitchecker
go build ./... && go install ./...
awslimitchecker --help
```

When making changes:

1. make sure you add relevant tests (there is a github action doing codecov validation)
2. make sure the existing tests pass `go test ./...` from root directory
3. make sure the changes passes [golangci-lint](https://golangci-lint.run/) `golangci-lint run` from root directory
