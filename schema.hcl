table "atlas_schema_revisions" {
  schema  = schema.test
  collate = "utf8mb4_bin"
  column "version" {
    null = false
    type = varchar(255)
  }
  column "description" {
    null = false
    type = varchar(255)
  }
  column "type" {
    null     = false
    type     = bigint
    default  = 2
    unsigned = true
  }
  column "applied" {
    null    = false
    type    = bigint
    default = 0
  }
  column "total" {
    null    = false
    type    = bigint
    default = 0
  }
  column "executed_at" {
    null = false
    type = timestamp
  }
  column "execution_time" {
    null = false
    type = bigint
  }
  column "error" {
    null = true
    type = longtext
  }
  column "error_stmt" {
    null = true
    type = longtext
  }
  column "hash" {
    null = false
    type = varchar(255)
  }
  column "partial_hashes" {
    null = true
    type = json
  }
  column "operator_version" {
    null = false
    type = varchar(255)
  }
  primary_key {
    columns = [column.version]
  }
}
table "features" {
  schema = schema.test
  column "id" {
    null           = false
    type           = bigint
    auto_increment = true
  }
  column "flag" {
    null = false
    type = varchar(255)
  }
  primary_key {
    columns = [column.id]
  }
  index "features_flag" {
    unique  = true
    columns = [column.flag]
  }
}
table "players" {
  schema = schema.test
  column "id" {
    null           = false
    type           = bigint
    auto_increment = true
  }
  column "username" {
    null = false
    type = varchar(16)
  }
  column "pw_hash" {
    null = false
    type = varchar(255)
  }
  primary_key {
    columns = [column.id]
  }
  index "players_username" {
    unique  = true
    columns = [column.username]
  }
}
table "requests" {
  schema = schema.test
  column "id" {
    null           = false
    type           = bigint
    auto_increment = true
  }
  primary_key {
    columns = [column.id]
  }
}
schema "test" {
  charset = "utf8mb4"
  collate = "utf8mb4_0900_ai_ci"
}
