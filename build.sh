#!/bin/bash

go mod vendor && go fmt ./... && go mod tidy -v && go mod vendor && make darwin && make linux && make windows
