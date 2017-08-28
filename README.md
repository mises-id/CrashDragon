## What is CrashDragon?

CrashDragon is a simple Minidump server, inspired by
[simple-breakpad-server][]. It's meant to be used as backend for Apps that
either use Google's [Breakpad client][bp] or its successor, [Crashpad][].

## GSoC

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
