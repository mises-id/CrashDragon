BEGIN TRANSACTION;

-- Add field to indicate Report processing time
ALTER TABLE "reports" ADD "processing_time" numeric;

END TRANSACTION;
