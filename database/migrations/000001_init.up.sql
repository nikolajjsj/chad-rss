CREATE TABLE users(
  id SERIAL PRIMARY KEY,
  email VARCHAR NOT NULL,
  password VARCHAR NOT NULL,
  created_at timestamp default now(),
  updated_at timestamp default now()
);
