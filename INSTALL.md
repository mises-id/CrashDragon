# Installing CrashDragon
The following guide will help you to install CrashDragon.

## Dependencies
First, make sure the OS is up to date, execute the following commands:
```
sudo apt-get update && apt-get upgrade
```

After that, install the actual dependencies:
```
sudo apt-get install golang git libcurl4-gnutls-dev rsync sassc autotools-dev autoconf make libjsoncpp-dev
```

_Note:_ You require a PostgreSQL database to run CrashDragon!

## Configuring environment
The following environment variables  specify the paths Go uses and the mode of Gin (the web framework). Depending on your setup put these commands into your bashrc file, or execute them:
```
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
export GIN_MODE=release
```

## Get code
The following command downloads CrashDragon and it's submodules, it will output a warning that there are no buildable Go files, which can be ignored. The second one changes the current directory to the project's source directory:
```
go get -d code.videolan.org/videolan/CrashDragon
cd $GOPATH/src/code.videolan.org/videolan/CrashDragon/
```

## Go dependencies
Install/Update govendor and use it to download the go libraries and project CrashDragon depends on:
```
go get -u github.com/kardianos/govendor
govendor sync
```

## Building
Run make to build the executable:
```
make
```

## Configuration
Now you have to configure CrashDragon, basically the SQL connection string and the IP/socket the program will listen to, the configuration file is generated on the first run of `crashdragon` or can be copy-pasted from the example below:
```
DatabaseConnection = "host=localhost user=crashdragon password=crashdragon dbname=crashdragon sslmode=disable"
UseSocket = false
BindAddress = "0.0.0.0:8080"
BindSocket = "/var/run/crashdragon/crashdragon.sock"
ContentDirectory = "/root/CrashDragon/Files"
TemplatesDirectory = "./templates"
AssetsDirectory = "./assets"
```

The configuration file can be found in `$HOME/CrashDragon/config.toml` and can be edited with whatever editor you like.

## Running
The server can be started by running `./bin/crashdragon`, this will keep the server process in the foreground. You can write init.d or systemctl scripts based on this command.
