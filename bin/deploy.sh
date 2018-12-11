#!/bin/bash

set -eu

mackerelApiToken=${MACKEREL_API_TOKEN}
retireDecisionPeriodHour=${RETIRE_DECISION_PERIOD_HOUR:-24}
retireDryRun=${RETIRE_DRY_RUN:-"true"}
region=${REGION:-$(aws configure get region)}
profile=${AWS_PROFILE:-default}

# build
mkdir -p ./artifact && GOARCH=amd64 GOOS=linux go build -o artifact/mackerenai

# mkdir temporary s3 bucket for deploy
bucket_name="temp-mackerenai-sam-$(openssl rand -hex 8)"
aws --profile $profile s3 mb "s3://${bucket_name}" --region $region

# deploy
aws --profile $profile cloudformation package \
  --output-template-file artifact/output.yaml \
  --template-file mackerenai.yaml \
  --s3-bucket="${bucket_name}"

aws --profile $profile cloudformation deploy \
  --template-file artifact/output.yaml \
  --stack-name MackerenaiSAM \
  --capabilities CAPABILITY_NAMED_IAM \
  --parameter-overrides "MackerelApiToken=${mackerelApiToken}" "RetireDecisionPeriodHour=${retireDecisionPeriodHour}" "RetireDryRun=${retireDryRun}" \
  --region $region

# delete temporary s3 bucket
aws --profile $profile s3 rb --force "s3://${bucket_name}" --region $region
