BEGIN TRANSACTION;

-- Drop no longer needed xxx_crash_count from crashes
ALTER TABLE "crashes" DROP COLUMN IF EXISTS all_crash_count, DROP COLUMN IF EXISTS win_crash_count, DROP COLUMN IF EXISTS mac_crash_count, DROP COLUMN IF EXISTS lin_crash_count;

-- Add new module field which is with signature a unique index
ALTER TABLE "crashes" ADD "module" text;
CREATE UNIQUE INDEX idx_crash_signature_module ON "crashes"("signature", "module");
DROP INDEX IF EXISTS idx_crash_signature;

-- Also add module as field on Reports
ALTER TABLE "reports" ADD "module" text;

-- Add migrations Table
CREATE TABLE IF NOT EXISTS "migrations" ("id" uuid NOT NULL DEFAULT NULL,"created_at" timestamp with time zone,"updated_at" timestamp with time zone,"component" text,"version" text , PRIMARY KEY ("id"));
INSERT INTO "migrations" ("id","created_at","updated_at","component","version") VALUES ('8badb4c8-9b9e-47a7-b753-b21b9785254b',NOW(),NOW(),'database','1.2.0');
INSERT INTO "migrations" ("id","created_at","updated_at","component","version") VALUES ('5f2bc875-4d09-4562-b5f8-29857c746153',NOW(),NOW(),'crashdragon','');

-- Migrate fixed column from boolean to timestamp
ALTER TABLE crashes RENAME fixed TO fixed_old;
ALTER TABLE crashes ADD fixed timestamp with time zone;
UPDATE crashes SET fixed = NOW() WHERE fixed_old = true;
UPDATE crashes SET fixed = NULL WHERE fixed_old = false;
ALTER TABLE crashes DROP fixed_old;

END TRANSACTION;
