Mackerenai
====

A tool to automatically delete unused hosts from [Mackerel](https://mackerel.io/).

Mackerel is a monitoring service of SaaS.
Automatically monitor AWS resources with [AWS integration](https://mackerel.io/ja/docs/entry/integrations/aws). 
However, there is a problem that resources that you don't want to monitor or already deleted are left on Mackerel.

Mackerenai is a useful tool to automatically delete such resources.

# Deploy

This tool is executed by AWS Lambda function, needs to be deployed with [AWS SAM](https://docs.aws.amazon.com/lambda/latest/dg/serverless_app.html).

```bash
# Install SAM CLI
$ pip install aws-sam-cli
$ sam --version
SAM CLI, version 0.8.1

# deploy using AWS SAM
$ ./bin/deploy

# Install in Tokyo, the period is 12 hours, dry run is false (actually delete), delete every day at 7:00pm UTC.
$ MACKEREL_API_KEY="xxxxx" \
  RETIRE_DECISION_PERIOD_HOUR=12 \
  RETIRE_DRY_RUN=false \
  SCHEDULE="cron(0 19 * * ? *)" \
  REGION=ap-northeast-1 \
  ./bin/deploy.sh
```

# Options

## MACKEREL_API_KEY

API key for operating Mackerel.

## RETIRE_DECISION_PERIOD_HOUR

Period to judge that it is not used in Mackerel. If metrics is not acquired during this period, delete that host.  
The default is `24` (hours), unit is `Hour`.

## RETIRE_DRY_RUN

Dry run. If it is true, it is not actually deleted.  
The default is `true`.

# Contribution

1. Fork ([https://github.com/ohsawa0515/mackerenai/fork](https://github.com/ohsawa0515/mackerenai/fork))
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Run test suite with the `go test ./...` command and confirm that it passes
6. Run `gofmt -s`
7. Create new Pull Request

# License

See [LICENSE](https://github.com/ohsawa0515/mackerenai/blob/master/LICENSE).

# Author

Shuichi Ohsawa (@ohsawa0515)

