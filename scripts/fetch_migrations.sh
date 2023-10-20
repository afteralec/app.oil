#!/bin/bash

curl --location --remote-header-name --remote-name https://raw.githubusercontent.com/afteralec/db.oil/main/schema.hcl \
  && atlas migrate diff --dev-url "mysql://root:pass@:3306/clean" --to "file://schema.hcl"
