CREATE TABLE users(
  user_id UUID PRIMARY KEY,
  email varchar(50) UNIQUE NOT NULL,
  password TEXT NOT NULL,
  first_name varchar(50),
  last_name varchar(50),
  user_type varchar(50) NOT NULL,
  token TEXT,
  refresh_token TEXT,
  updated_at TIMESTAMP
)

