USE auth_db;

DROP TABLE IF EXISTS users;

CREATE TABLE users
(
  id          INT(10) PRIMARY KEY AUTO_INCREMENT,
  user_id     VARCHAR(255),
	password    VARCHAR(255)
);
