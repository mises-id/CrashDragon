BEGIN TRANSACTION;

-- Create table for many2many relationship
CREATE TABLE IF NOT EXISTS "crash_versions" (
    "crash_id" uuid REFERENCES crashes(id) ON DELETE RESTRICT,
    "version_id" uuid REFERENCES versions(id) ON DELETE RESTRICT,
    PRIMARY KEY ("crash_id","version_id")
);

-- Populate table with existing crash/version relationships based on reports
INSERT INTO crash_versions SELECT crash_id, version_id FROM reports GROUP BY crash_id, version_id;

-- Drop no longer needed version_id from crashes
ALTER TABLE crashes DROP COLUMN IF EXISTS version_id;

END TRANSACTION;
