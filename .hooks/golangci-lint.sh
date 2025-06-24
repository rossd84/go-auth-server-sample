#!/bin/bash

export GO111MODULE=on
golangci-lint run ./...
