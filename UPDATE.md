# Updating

To update the project run 

```
git fetch
``` 

followed by

```
git checkout vx.x.x
```

(enter the most current version)

Then the submodule and breakpads submodules need to be updated:

```
git submodule update --init --recursive
```

After that you have to run the database migrations. This can be done by running

```
psql -U username -d dataBase -a -f migrations/file.sql
```
or the command matching your database set up. The migrations must be run after
each other, so if you use v0.0.1 and want to update to v1.0.0 you have to run
every migration on the way there (v0.0.1 -> v0.1.0 -> v0.1.1 -> v1.0.0) else
there may be stuff missing.
