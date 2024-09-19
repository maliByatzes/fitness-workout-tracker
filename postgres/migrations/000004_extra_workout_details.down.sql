ALTER TABLE "workout_report" DROP CONSTRAINT IF EXISTS "workout_report_user_id_fkey";
ALTER TABLE "workout_exercise_status" DROP CONSTRAINT IF EXISTS "workout_exercise_status_workout_exercise_id_fkey";

DROP TABLE IF EXISTS "workout_report";
DROP TABLE IF EXISTS "wokrkout_exercise_status";
