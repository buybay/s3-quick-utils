# s3-quick-utils

The goal of these utilities is to be able to execute some basic S3
actions as fast as possible, compare with the standard AWS CLI.

```
S3 utilities focus on reduce execution times.

Usage:
  s3-quick-utils [command]

Available Commands:
  counter     Count objects in a S3 bucket
  delete      Delete all the objects in a S3 bucket
  help        Help about any command

Flags:
  -h, --help             help for s3-quick-utils
      --profile string   AWS local profile
      --region string    AWS region for the S3 bucket (default "eu-central-1")

Use "s3-quick-utils [command] --help" for more information about a command.

```

## counter

Count the objects in a bucket. The implementation is optimized for a
bucket with uniform distribution on keys prefixes.

```
Count objects in a S3 bucket

Usage:
  s3-quick-utils counter [bucket_name] [flags]

Examples:
counter xxx-eu-central-1-production

Flags:
  -h, --help   help for counter

Global Flags:
      --profile string   AWS local profile
      --region string    AWS region for the S3 bucket (default "eu-central-1")
```

## delete

Delete all the objects in a bucket. The implementation is optimized
for a bucket with uniform distribution on keys prefixes.

```
Delete all the objects in a S3 bucket

Usage:
  s3-quick-utils delete [bucket_name] [flags]

Examples:
delete xxx-eu-central-1-production

Flags:
  -h, --help   help for delete

Global Flags:
      --profile string   AWS local profile
      --region string    AWS region for the S3 bucket (default "eu-central-1")
```
