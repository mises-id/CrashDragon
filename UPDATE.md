# Updating

## Version 1.0.x to 1.1.0

If you upgrade from version 1.0.x to 1.1.0 your config file will be moved from
`$HOME/CrashDragon/config.toml` to `PREFIX/etc/crashdragon.toml`! Also the
directories for assets and templates change to `PREFIX/share/crashdragon/...`,
this must be changed in the config file *manually*! An example can be found in
`INSTALL.md`. If you changed assets or templates you will also have to move
these. Uploaded files (symfiles, minidumps and generated textfiles) are now
located in `PREFIX/share/crashdragon/files` so you will have to either move
them there or adapt the configuration to point to your custom directory!

## All versions

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
or the command matching your database setup. The migrations must be run after
another, so if you use v0.0.1 and want to update to v1.0.0 you have to run
every migration on the way there (v0.0.1 -> v0.1.0 -> v0.1.1 -> v1.0.0) else
there may be stuff missing. If there is no migration file for your version you
only have to run the remaining migrations on the way to your version.

To compile the new version and also install it in your spcified `PREFIX` run
the following:

```
make prefix=PREFIX install
```
