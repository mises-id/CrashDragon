BEGIN TRANSACTION;

-- Add flag to indicate if a Crash is fixed (issue #19)
UPDATE "crashes" SET "fixed" = false;

-- Add flag to indicate if a Version is ignored (issue #21)
UPDATE "versions" SET "ignore" = false;

END TRANSACTION;
