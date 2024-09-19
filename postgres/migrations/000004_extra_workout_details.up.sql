CREATE TYPE "status_enum" AS ENUM('pending', 'completed');

CREATE TABLE IF NOT EXISTS "workout_exercise_status" (
    "id" SERIAL NOT NULL,
    "workout_exercise_id" INTEGER NOT NULL,
    "status" status_enum NOT NULL,
    "comments" TEXT,
    "completed_at" TIMESTAMPTZ,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL,
    CONSTRAINT "workout_exercise_status_pkey" PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "workout_report" (
    "id" SERIAL NOT NULL,
    "user_id" INTEGER NOT NULL,
    "start_date" DATE NOT NULL,
    "end_date" DATE NOT NULL,
    "total_workouts" INTEGER NOT NULL DEFAULT 0,
    "completed_workouts" INTEGER NOT NULL DEFAULT 0,
    "completion_percentage" REAL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL,
    CONSTRAINT "workout_report_pkey" PRIMARY KEY ("id")
);

ALTER TABLE "workout_exercise_status" ADD CONSTRAINT "workout_exercise_status_workout_exercise_id_fkey" FOREIGN KEY ("workout_exercise_id") REFERENCES "workout_exercise"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE "workout_report" ADD CONSTRAINT "workout_report_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
