# malawi get status
[![Quality Gate Status](http://sonar.pgzaoffice.local:9000/api/project_badges/measure?project=malawi-get-status&metric=alert_status)](http://sonar.pgzaoffice.local:9000/dashboard?id=malawi-get-status)
[![Code Smells](http://sonar.pgzaoffice.local:9000/api/project_badges/measure?project=malawi-get-status&metric=code_smells)](http://sonar.pgzaoffice.local:9000/dashboard?id=malawi-get-status)
[![Maintainability Rating](http://sonar.pgzaoffice.local:9000/api/project_badges/measure?project=malawi-get-status&metric=sqale_rating)](http://sonar.pgzaoffice.local:9000/dashboard?id=malawi-get-status)
[![Security Rating](http://sonar.pgzaoffice.local:9000/api/project_badges/measure?project=malawi-get-status&metric=security_rating)](http://sonar.pgzaoffice.local:9000/dashboard?id=malawi-get-status)

This boilerplate code is setup mostly to illustrate the desired structure and design of future Lambda functions. It
provides you with a full test suite, relevant service layer, repository layer, producers, library files and entry
through the controller. Adjust the boilerplate as you may require to fulfill your project requirements.

## Prerequisites
Before getting started with this project, you should have the following installed:
- [Git Bash](https://git-scm.com/downloads)
- [GoLang](https://golang.org/doc/install)

## Getting Started
To get started with this boilerplate project, you can clone it from the GitLabs repository.
```bash
$ git clone some-git-url
$ cd lambda-boiler
```

## Module init
If you need to redo the modules
```bash
$ go mod init example.com/hello
```

You should immediately then be able to run the test suite to ensure that all is as it should be.
```bash
$ go test ./...
```

You can also run coverage.
```bash
$ go test -coverprofile cover.out ./...
$ go tool cover -html=cover.out
```

## Our Commonly Used External Libraries
Here is short list of libraries that are used when needing to solve a particular problem
 - [AWS Client v2](https://aws.github.io/aws-sdk-go-v2/)
 - [MSSQL Client](https://github.com/denisenkom/go-mssqldb)

## Guidelines
Below are some links that will assist with continuing with the development and structure of this boilerplate project.

General:
 - [How to write a README](https://www.makeareadme.com/)
 - [README template #1](https://gist.github.com/PurpleBooth/109311bb0361f32d87a2)
 - [README template #2](https://gist.github.com/fvcproductions/1bfc2d4aecb01a834b46)
 - [Cyclomatic Complexity](https://webuniverse.io/cyclomatic-complexity-refactoring-tips/) just to keep in mind.

Keeping a README up-to-date is quite useful when someone wants to get up and running with a project.
Keeping in mind cyclomatic complexity will save time in the long run. Try keep code as simple and elegant as possible.

Unit test related:
 - [SQL Mock](https://github.com/DATA-DOG/go-sqlmock)

Code style guidelines
 - [Golang](https://directpayonline.atlassian.net/wiki/spaces/PAYG/pages/1851392004/GoLang+Style+Guide)

## Environment
You can use a `.env` to add environment variables locally to run the sample `integration` test. It's more just executing
the code as opposed to being an integration test but it allows you to ensure that the end result is as it is expected.
```text
AWS_REGION=some-region
SECRET_NAME=some-secret
`LOG_LEVEL = WARN`
```
You can find the `.env.example` file as an example as well.
