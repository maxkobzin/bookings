CREATE TABLE IF NOT EXISTS restrictions(
   id serial PRIMARY KEY,
   restriction_name VARCHAR (255) DEFAULT (''),
   created_at timestamptz NOT NULL DEFAULT (now()),
   updated_at timestamptz NOT NULL DEFAULT (now())
);