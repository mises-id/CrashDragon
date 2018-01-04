-- Fix database change introduced in commit c564789 to fix issue #6
UPDATE versions 
    SET git_repo = (
        SELECT git_repo
        FROM products
        WHERE versions.product_id = products.id
    );

