BEGIN TRANSACTION;

-- Add field to indicate Report processing time
ALTER TABLE "reports" ADD "processing_time" numeric;

-- Remove Report TXT content from database
ALTER TABLE "reports" DROP "report_content_txt";

END TRANSACTION;
