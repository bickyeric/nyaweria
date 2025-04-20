CREATE TABLE IF NOT EXISTS users (
  id SERIAL,
  username VARCHAR(255) NOT NULL,
  name VARCHAR(255),
  profile_picture VARCHAR(255),
  description VARCHAR(255)
);
