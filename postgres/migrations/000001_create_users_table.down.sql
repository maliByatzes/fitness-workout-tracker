ALTER TABLE "profile" DROP CONSTRAINT IF EXISTS "profile_user_id_fkey";

DROP INDEX IF EXISTS "user_username_key";
DROP INDEX IF EXISTS "user_email_key";

DROP TABLE IF EXISTS "profile";
DROP TABLE IF EXISTS "user";
