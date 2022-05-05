CREATE TABLE IF NOT EXISTS rooms(
   id serial PRIMARY KEY,
   room_name VARCHAR (255) DEFAULT (''),
   created_at timestamptz NOT NULL DEFAULT (now()),
   updated_at timestamptz NOT NULL DEFAULT (now())
);