#!/bin/bash

chmod u+x ./scripts/fetch_migrations.sh
./scripts/fetch_migrations.sh
atlas schema clean -u "mysql://root:pass@:3306/test" --auto-approve
atlas migrate apply -u "mysql://root:pass@:3306/test"
