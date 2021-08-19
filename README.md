# awsacc

***A trusty helper for working with AWS account IDs.***

![](https://img.shields.io/badge/Go-1.16%2B-blue)
![](https://img.shields.io/github/issues/cbrgm/awsacc)
![](https://img.shields.io/github/license/cbrgm/awsacc)


Working with AWS account IDs often involves more manual effort than necessary. Often account IDs in ARNs have to be manually looked up in different files and compared with existing IDs. With a few accounts this can be manageable, but with a large number of accounts to manage it is difficult to keep track.

`awsacc` allows you to:

* quickly turn account ids into their ***human-readable names in any file***.
* ***lookup account ids*** by their corresponding ***name and vice versa*** via CLI and thus make them ***processable***.

It's designed to ***speeds up*** reviewing changes and finding account ID errors (e.g. misspelled, or unknown IDs) in files.
It can also be used as ***part of a CI/CD pipeline*** to match used account IDs in files against a defined list of valid IDs.

## Installation

You can find the latest release [here](https://github.com/cbrgm/awsacc/releases) (Linux, Win, Mac 64-bit).

You can also build `awsacc` from source (Go 1.16+ required).

```
git clone git@github.com:cbrgm/awsacc.git && cd awsacc
go mod vendor && make
```

Another option is to use `awsacc` inside a container

```bash
docker run --rm -it \
   -v ~/.aws/accounts.json:/data/accounts.json \
   -v $(pwd):$(pwd) \
   cbrgm/awsacc:latest -f $(pwd)/file.json
```

where `~/.aws/accounts.json` is your config file and `$(pwd)/file.json` is the file to check.

### Configuration 

A configuration file must be created so that the tool can search and substitute AWS account names and IDs. 
The configuration file is a JSON file consisting of accounts (tuples of names and IDs). 

By default, it is stored under the path `$HOME/.aws/accounts.json`. 
You can also reference a config file via the environment variable `AWSACC_CONFIG` (Example: `export AWSACC_CONFIG=/path/to/config.json`).

***~/.aws/accounts.json***
```json
{
    "Accounts": [
        {
            "Id": "620791285726",
            "Name": "Foo",
        },
    ]
}
```
The schema can be created ***manually*** or easily obtained using the `aws` CLI.

```bash
aws organizations list-accounts > ~/.aws/accounts.json
```

> Note: This operation can be called only from the organization's management account or by a member account that is a delegated administrator for an AWS service. ([API Reference](https://docs.aws.amazon.com/organizations/latest/APIReference/API_ListAccounts.html))

## Usage

The use of the tool can be accessed through the `help` subcommand

```
Usage: awsacc [options] [subcommand] [options] <args>

Replaces AWS account ids with their human-readable names
	-c Colored output. Default: false
	-f Path to the input files. Default: Stdin
	-s Strict mode, return on error. Default: false
	-v Verbose output. Default: false

Subcommand: search, ls
Description: Searches and prints out account ids or account names
	-c Colored output. Default: false
	-s Strict mode, return on error. Default: false
	-v Verbose output. Default: false
```

## Examples

Here are some examples of what you can do with `awsacc`

### Replace and highlight AWS account ids with their names

We assume that a file `test.json` exists. (Does not have to be JSON only, ***any kind of text file is possible***).

```bash
$ cat test.json
```

Output:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Example`:",
            "Effect": "Allow",
            "Action": [
                "sts:AssumeRole"
            ],
            "Resource": [
                "arn:aws:iam::614912345005:role/somebody",
                "arn:aws:iam::866714215829:role/somebody",
                "arn:aws:iam::012073564612:role/somebody"
            ]
        }
    ]
}
```

We can now replace the account IDs with the human-readable names, color-code the changed lines and print them out

```bash
$ cat test.json | awsacc -c -v
```

Output:

```diff
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "Example`:",
      "Effect": "Allow",
      "Action": [
        "sts:AssumeRole"
      ],
      "Resource": [
+        "arn:aws:iam::production:role/somebody",
-        "arn:aws:iam::614912345005:role/somebody",
+        "arn:aws:iam::integration:role/somebody",
-        "arn:aws:iam::866714215829:role/somebody",
+        "arn:aws:iam::development:role/somebody",
-        "arn:aws:iam::012073564612:role/somebody"
      ]
    }
  ]
}
```

Or we let only the changed lines be colored, without showing removed lines

```bash
$ cat test.json | awsacc -c
```

Output:

```diff
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "Example`:",
      "Effect": "Allow",
      "Action": [
        "sts:AssumeRole"
      ],
      "Resource": [
+        "arn:aws:iam::production:role/somebody",
+        "arn:aws:iam::integration:role/somebody",
+        "arn:aws:iam::development:role/somebody",
      ]
    }
  ]
}
```

### Lookup AWS account ids or names

List all accounts ids of accounts containing `dev` in their name
```bash
awsacc ls dev
```

List all accounts names of accounts containing `00124` in their account id
```bash
awsacc ls 00124
```

List all accounts names of accounts `312345643213`, `612345343211` and `822345643215`
```bash
awsacc ls 312345643213 612345343211 822345643215
```

List all accounts ids of accounts containing `doesntexist` in their name, use strict mode `-s` to exit with `err 1` when no results have been found
```bash
awsacc ls -s doesntexist
```

Count all accounts ids of accounts containing `dev` or `int` or `prd` in their name
```bash
awsacc ls dev int prod | wc -l
```

List all accounts and search for an account or id using `fzf` and copy the result to your clipboard
```bash
awsacc ls -v '*' | fzf | pbcopy
```

### Usage in your CI/CD pipeline

It can also be used as part of a CI/CD pipeline to match used account IDs in files against a defined list of valid IDs from your configuration.

`awsacc` searches below the path recursively in all directories for matching files and returns an error code `1` if an ID was found that is not stored in the configuration.

```bash
awsacc -s -f './path/to/*.json'
```

is equivalent to

```bash
cat './path/to/*.json' | awsacc -s
```

## Contributing & License

Feel free to submit changes! See the [Contributing Guide](https://github.com/cbrgm/contributing/blob/master/CONTRIBUTING.md). This project is open-source and is developed under the terms of the [Apache 2.0 License](https://github.com/cbrgm/awsacc/blob/master/LICENSE).
