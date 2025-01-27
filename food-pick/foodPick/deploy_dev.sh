#!/bin/bash

sam build -t ./template.yaml && sam deploy --config-env apne2-dev --no-progressbar --no-confirm-changeset
