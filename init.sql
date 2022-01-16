create database nftmaker;
use nftmaker;

CREATE TABLE IF NOT EXISTS users (
uuid VARCHAR(36) NOT NULL,
name VARCHAR(255) NOT NULL,
surname VARCHAR(255) NOT NULL,
email VARCHAR(255) UNIQUE NOT NULL,
password VARCHAR(255) NOT NULL,
created_at INT NOT NULL,
modified_at INT NOT NULL,
PRIMARY KEY (uuid)
) ;