FROM golang:latest

ENV GOPATH /go
ENV PATH $PATH:$GOPATH/bin

ADD . $GOPATH/src/git.1750studios.com/GSoC/CrashDragon
WORKDIR $GOPATH/src/git.1750studios.com/GSoC/CrashDragon

RUN apt-get update && apt-get -y install libcurl4-gnutls-dev rsync postgresql

RUN go get -u github.com/kardianos/govendor
RUN govendor sync

RUN make

RUN postgres -D /usr/local/pgsql/data >/var/log/postgres 2>&1 &
RUN psql -c "CREATE USER CrashDragon; CREATE DATABASE CrashDragon OWNER CrashDragon;"

EXPOSE 8080
CMD ["./bin/CrashDragon"]
