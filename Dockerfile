FROM golang:1.9-stretch

ENV GOPATH /go
ENV PATH $PATH:$GOPATH/bin
ENV GIN_MODE release

ADD . $GOPATH/src/code.videolan.org/videolan/CrashDragon
WORKDIR $GOPATH/src/code.videolan.org/videolan/CrashDragon

RUN apt-get update && apt-get -y install libcurl4-gnutls-dev rsync postgresql sassc autotools-dev autoconf

RUN go get -u github.com/kardianos/govendor
RUN govendor sync

RUN make

RUN /etc/init.d/postgresql start && su postgres -c 'createuser -w crashdragon' && su postgres -c 'createdb -w -O crashdragon crashdragon'
RUN echo "local all all trust" > /etc/postgresql/9.6/main/pg_hba.conf
RUN echo "host all all all trust" >> /etc/postgresql/9.6/main/pg_hba.conf

EXPOSE 8080
CMD /etc/init.d/postgresql start && sleep 15 && ./bin/crashdragon
