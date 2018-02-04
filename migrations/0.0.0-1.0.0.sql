-- Fix database change introduced in commit c564789 to fix issue #6
ALTER TABLE versions ADD git_repo text;
UPDATE versions 
    SET git_repo = (
        SELECT git_repo
        FROM products
        WHERE versions.product_id = products.id
    );
ALTER TABLE products DROP COLUMN git_repo;

-- Add flag to indicate a Crash is fixed
ALTER TABLE "crashes" ADD "fixed" boolean;
