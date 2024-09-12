ALTER TABLE "workout_exercise" DROP CONSTRAINT IF EXISTS "workout_exercise_exercise_id_fkey";
ALTER TABLE "workout_exercise" DROP CONSTRAINT IF EXISTS "workout_exercise_workout_id_fkey";
ALTER TABLE "workout" DROP CONSTRAINT IF EXISTS "workout_user_id_fkey";

DROP TABLE IF EXISTS "workout_exercise";
DROP TABLE IF EXISTS "exercise";
DROP TABLE IF EXISTS "workout";
