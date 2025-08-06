CREATE TABLE IF NOT EXISTS "members" (
  "id" BIGSERIAL PRIMARY KEY,
  "email" TEXT,
  "phone" TEXT NOT NULL,
  "name" TEXT NOT NULL,
  "group_id" BIGINT NOT NULL,
  "user_id" UUID NOT NULL,
  "created_at" TIMESTAMP DEFAULT (now()),
  "updated_at" TIMESTAMP,
  "deleted_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "groups" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" TEXT NOT NULL,
  "user_id" UUID, -- identifies the fellow who created the group
  description TEXT,
  "created_at" TIMESTAMP DEFAULT (now()),
  "updated_at" TIMESTAMP,
  "deleted_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "rsvps" (
  "id" BIGSERIAL PRIMARY KEY,
  "member_id" BIGINT NOT NULL,
  "event_id" BIGINT NOT NULL,
  "has_paid" BOOL DEFAULT false,
  "payment_data" JSONB,
  "payment_reference_id" BIGINT,
  "created_at" TIMESTAMP DEFAULT (now()),
  "updated_at" TIMESTAMP,
  "deleted_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "events" (
  "id" BIGSERIAL PRIMARY KEY,
  "title" TEXT,
  "image" TEXT,
  "description" TEXT,
  "group_id" BIGINT NOT NULL,
  "status" TEXT,
  "is_paid" bool DEFAULT false,
  "amount" BIGINT,
  "created_at" TIMESTAMP DEFAULT (now()),
  "updated_at" TIMESTAMP,
  "deleted_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "group_admins" (
  "id" BIGSERIAL PRIMARY KEY,
  "group_id" BIGINT NOT NULL,
  "member_id" BIGINT NOT NULL,
  "created_at" TIMESTAMP DEFAULT (now()),
  "updated_at" TIMESTAMP,
  "deleted_at" TIMESTAMP
);

CREATE UNIQUE INDEX ON "members" ("user_id", "group_id");

CREATE INDEX ON "rsvps" ("member_id");

CREATE INDEX ON "rsvps" ("event_id");

CREATE INDEX ON "events" ("group_id");

CREATE UNIQUE INDEX ON "group_admins" ("member_id", "group_id");

CREATE UNIQUE INDEX ON "groups" ("user_id", "name");

ALTER TABLE "members" ADD FOREIGN KEY ("user_id") REFERENCES "auth"."users" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "members" ADD FOREIGN KEY ("group_id") REFERENCES "groups" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "events" ADD FOREIGN KEY ("group_id") REFERENCES "groups" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "rsvps" ADD FOREIGN KEY ("event_id") REFERENCES "events" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "rsvps" ADD FOREIGN KEY ("member_id") REFERENCES "members" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "group_admins" ADD FOREIGN KEY ("member_id") REFERENCES "members" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "group_admins" ADD FOREIGN KEY ("group_id") REFERENCES "groups" ("id") ON DELETE CASCADE ON UPDATE CASCADE;