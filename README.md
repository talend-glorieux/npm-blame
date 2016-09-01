# npm-blame

[![Build Status](https://travis-ci.org/talend-glorieux/npm-blame.svg?branch=master)](https://travis-ci.org/talend-glorieux/npm-blame)

Reports on npm packages common errors and useless files.

## Install

Get the latest release for you operating system architecture.

## Usage

Run `npm-blame` from inside your project's node_module folder.

## Build 
* Get the [latest Golang release](https://golang.org/dl/)
* Set up your workspace
* Run `go get github.com/talend-glorieux/npm-blame` 
* Go to the projects folder. `$GOPATH/src/github.com/talend-glorieux/npm-blame`
* Run `go install ./cmd/npm-blame`

## Status

npm-blame is still a work in progress. It will need to handle more errors, add a
report mode that alerts maintainers as well as a bit of refactoring.
