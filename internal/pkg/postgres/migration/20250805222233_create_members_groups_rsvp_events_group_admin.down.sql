-- Drop foreign key constraints first
ALTER TABLE "rsvps" DROP CONSTRAINT IF EXISTS "rsvp_member_id_fkey";

ALTER TABLE "rsvps" DROP CONSTRAINT IF EXISTS "rsvp_event_id_fkey";

ALTER TABLE "events" DROP CONSTRAINT IF EXISTS "events_group_id_fkey";

ALTER TABLE "members" DROP CONSTRAINT IF EXISTS "members_group_id_fkey";

ALTER TABLE "members" DROP CONSTRAINT IF EXISTS "members_user_id_fkey";

-- Drop indexes
DROP INDEX IF EXISTS "group_admin_member_id_group_id_idx";

DROP INDEX IF EXISTS "events_group_id_idx";

DROP INDEX IF EXISTS "rsvp_event_id_idx";

DROP INDEX IF EXISTS "rsvp_member_id_idx";

DROP INDEX IF EXISTS "members_user_id_group_id_idx";

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS "group_admins";

DROP TABLE IF EXISTS "rsvps";

DROP TABLE IF EXISTS "events";

DROP TABLE IF EXISTS "members";

DROP TABLE IF EXISTS "groups";
