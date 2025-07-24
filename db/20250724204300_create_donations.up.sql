CREATE TABLE IF NOT EXISTS donations (
  id SERIAL,
  sender VARCHAR(255) NOT NULL,
  recipient_id INT NOT NULL,
  currency VARCHAR(255),
  amount VARCHAR(255),
  message VARCHAR(255)
);
