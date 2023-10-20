-- Create "features" table
CREATE TABLE `features` (`id` bigint NOT NULL AUTO_INCREMENT, `flag` varchar(255) NOT NULL, PRIMARY KEY (`id`), UNIQUE INDEX `features_flag` (`flag`)) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "players" table
CREATE TABLE `players` (`id` bigint NOT NULL AUTO_INCREMENT, `username` varchar(16) NOT NULL, `pw_hash` varchar(255) NOT NULL, PRIMARY KEY (`id`), UNIQUE INDEX `players_username` (`username`)) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "requests" table
CREATE TABLE `requests` (`id` bigint NOT NULL AUTO_INCREMENT, PRIMARY KEY (`id`)) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
