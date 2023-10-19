-- Create "atlas_schema_revisions" table
CREATE TABLE `atlas_schema_revisions` (`version` varchar(255) NOT NULL, `description` varchar(255) NOT NULL, `type` bigint unsigned NOT NULL DEFAULT 2, `applied` bigint NOT NULL DEFAULT 0, `total` bigint NOT NULL DEFAULT 0, `executed_at` timestamp NOT NULL, `execution_time` bigint NOT NULL, `error` longtext NULL, `error_stmt` longtext NULL, `hash` varchar(255) NOT NULL, `partial_hashes` json NULL, `operator_version` varchar(255) NOT NULL, PRIMARY KEY (`version`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "features" table
CREATE TABLE `features` (`id` bigint NOT NULL AUTO_INCREMENT, `flag` varchar(255) NOT NULL, PRIMARY KEY (`id`), UNIQUE INDEX `features_flag` (`flag`)) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "players" table
CREATE TABLE `players` (`id` bigint NOT NULL AUTO_INCREMENT, `username` varchar(16) NOT NULL, `pw_hash` varchar(255) NOT NULL, PRIMARY KEY (`id`), UNIQUE INDEX `players_username` (`username`)) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "requests" table
CREATE TABLE `requests` (`id` bigint NOT NULL AUTO_INCREMENT, PRIMARY KEY (`id`)) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
