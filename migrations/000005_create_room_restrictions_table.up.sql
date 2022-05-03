CREATE TABLE IF NOT EXISTS room_restrictions(
   id serial PRIMARY KEY,
   start_date DATE,
   end_date DATE,
   room_id int,
   reservation_id int,
   restriction_id int
);