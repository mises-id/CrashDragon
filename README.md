## What is CrashDragon?

CrashDragon is a simple Minidump server, inspired by
[simple-breakpad-server][]. It's meant to be used as backend for Apps that
either use Google's [Breakpad client][bp] or its successor, [Crashpad][].

## GSoC 2018

CrashDragon was part of the 2018 GSoC program, where the software got
extended and various critical bugs were fixed. During the GSoC period
the commits 5c05624c5bbc696434bf971c354b3d5be0f7901b to
2acc4984e05ad990e067e215b98bb1f35499ec37 have been done. The overall goal
was to make CrashDragon easier to use and to extend it's functionality.

Most of the GSoC work this year was researching RESTful JSON API concepts,
doing research about testing and performance tweaks to CrashDragon, which has
a rather big database by now.

### Features implemented
* More statistics about OS/Version on Crash (#33)
* Add selection of version besides selection of product (#32)
* Button to remove one specific crash (#31)
* Allow giving the path to the symbolicator binary in the config (#13)
* Show some info about a Crash on a crash page (#11)
* Add install target to makefile (#38)
* Display when a crash was marked as fixed (#36)
* Put symbols in separate filesystem locations (#15)

### Bugfixes
The following bugs have been fixed in the GSoC period:
* Uploading symbols for non-existing version of a product should not fail (#30)
* Crash on wrong breakpad output (#37)
* Get correct OS crash counts when filtering for Version (#40)
* Support graceful shutdown (#25)
* Improve separation of the "empty stacks" based on the dll in crashed in (#39)

The main changes were extensions to the UI, a JSON API and also first
API tests have been implemented.

## GSoC 2017

CrashDragon was part of the 2017 GSoC program, where the software got
implemented from ground up. This includes all commits from
e97412632ed6a7261330015c052fda29e7d867da to
46d5d442923a75feb9e93d8c664c734cb20d00a4. The goal was to have an server
which is performing very well under heavy load and is still easy to use.
Not all milestones could be achieved during the GSoC time, but I will keep
working on the project after GSoC has finished to implement those missing
features and to keep the software maintained.

### Features
* Management of Users, Products and Versions in the admin interface
* Upload of Symbol files for Products/Versions
* Upload of minidump files which get processed by the server
* Automatic grouping of multiple reports into crashes
* Linking from source files/lines to the respective file in a Git web interface
* View and download of reports in stacktrace or JSON version
* User authentication based on the `Auth`-header
* Mobile friendly frontend based on Foundation
* Diagrams on how many reports there are for each version/product/platform

### Not yet implemented features
* Improve the way of matching reports into crashes
* Improve pagination and report-to-report navigation
* Linking between crashes and issues in a bugtracker
* Add way to add new issues to bugtracker based on crash
* Specify and implement JSON API

The work on these features will be continued after GSoC finishes.

### Work outside of this repository
I also did some work outside of this repository, for example tweking the
macOS integration of the breakpad minidump sender and trying to integrate
this sender into the Windows versions of VLC, which didn't work out as
expected as it is very hard to compile breakpad under MinGW. I will also have
a look at these problems after the GSoC period has finished.

[simple-breakpad-server]: https://github.com/acrisci/simple-breakpad-server
[bp]: https://chromium.googlesource.com/breakpad/breakpad
[crashpad]: https://chromium.googlesource.com/crashpad/crashpad/
