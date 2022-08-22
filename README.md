
# awslimitchecker

Simple module to programatically retrieve your AWS account limits (whether they are supporter by servicequotas or not). It also provide a lightweight cli program to access the data.

## cli

A utility `awslimitchecker` CLI is provided, that exposes the module through a simple interface.

## development

To run the latest:

```shell
cd awslimitchecker/cmd
go install
awslimitechchecker --help
```
