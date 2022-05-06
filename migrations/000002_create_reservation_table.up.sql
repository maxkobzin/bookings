CREATE TABLE IF NOT EXISTS reservations(
   id serial PRIMARY KEY,
   first_name VARCHAR (255) DEFAULT (''),
   last_name VARCHAR (255) DEFAULT (''),
   email VARCHAR (255),
   phone VARCHAR (255) DEFAULT (''),
   start_date DATE,
   end_date DATE,
   room_id int NOT NULL,
   created_at timestamptz NOT NULL DEFAULT (now()),
   updated_at timestamptz NOT NULL DEFAULT (now())
);