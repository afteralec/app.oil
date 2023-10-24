#!/bin/bash

curl --location --remote-header-name --remote-name https://raw.githubusercontent.com/petrichormud/db/main/schema.hcl \
  && atlas schema clean -u "mysql://root:pass@:3306/test" --auto-approve \
  && atlas migrate diff --dev-url "mysql://root:pass@:3306/test" --to "file://schema.hcl"
