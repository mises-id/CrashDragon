DROP FUNCTION IF EXISTS migrate_counts();
CREATE FUNCTION migrate_counts() RETURNS integer AS $$
DECLARE
	crash RECORD;
	ver RECORD;
BEGIN
	RAISE NOTICE 'Migrating all crash counts...';
	
	FOR crash IN (SELECT * FROM crashes) LOOP
		FOR ver IN (SELECT * FROM versions) LOOP
		
			INSERT INTO crash_counts (id, created_at, updated_at, crash_id, version_id, os, count)
				SELECT
					uuid_generate_v4(), now(), now(), crash.id, ver.id, 'Windows NT', 
					count(*) FROM "reports"  WHERE (crash_id = crash.id AND version_id = ver.id AND os = 'Windows NT');
			RAISE NOTICE 'Migrating crash % os Windows NT version %...', crash.id, ver.slug;
			
			INSERT INTO crash_counts (id, created_at, updated_at, crash_id, version_id, os, count)
				SELECT
					uuid_generate_v4(), now(), now(), crash.id, ver.id, 'Mac OS X', 
					count(*) FROM "reports"  WHERE (crash_id = crash.id AND version_id = ver.id AND os = 'Mac OS X');
			RAISE NOTICE 'Migrating crash % os Mac OS X version %...', crash.id, ver.slug;
		
		END LOOP;
	END LOOP;
	
	RAISE NOTICE 'Done migrating crash counts.';
    RETURN 1;
END;
$$ LANGUAGE plpgsql;


DO $$
	BEGIN
		IF (SELECT version FROM migrations WHERE component = 'database') = '1.2.1' THEN
			RAISE NOTICE 'Database migration version is 1.2.1, migrating...';
			migrate_counts();
			DROP FUNCTION IF EXISTS migrate_counts();
			UPDATE migrations SET version = '1.3.0' WHERE component = 'database';
			RAISE NOTICE 'Database migration version is now 1.3.0';
		ELSE
			RAISE NOTICE 'Database migration version is not 1.2.1';
		END IF;
	END;
$$ LANGUAGE plpgsql;