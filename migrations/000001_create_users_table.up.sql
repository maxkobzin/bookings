CREATE TABLE IF NOT EXISTS users(
   id serial PRIMARY KEY,
   first_name VARCHAR (255) DEFAULT (''),
   last_name VARCHAR (255) DEFAULT (''),
   email VARCHAR (255),
   password VARCHAR (60),
   access_level int DEFAULT (1),
   created_at timestamptz NOT NULL DEFAULT (now()),
   updated_at timestamptz NOT NULL DEFAULT (now())
);