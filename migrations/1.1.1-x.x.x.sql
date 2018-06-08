BEGIN TRANSACTION;

-- Drop no longer needed xxx_crash_count from crashes
ALTER TABLE crashes DROP COLUMN IF EXISTS all_crash_count, DROP COLUMN IF EXISTS win_crash_count, DROP COLUMN IF EXISTS mac_crash_count, DROP COLUMN IF EXISTS lin_crash_count;

END TRANSACTION;
