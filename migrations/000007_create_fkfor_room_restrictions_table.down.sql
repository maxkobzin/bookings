ALTER TABLE "room_restrictions" DROP CONSTRAINT IF EXISTS "room_restrictions_rooms_id_fk";
ALTER TABLE "room_restrictions" DROP CONSTRAINT IF EXISTS "room_restrictions_reservations_id_fk";
ALTER TABLE "room_restrictions" DROP CONSTRAINT IF EXISTS "room_restrictions_restrictions_id_fk";