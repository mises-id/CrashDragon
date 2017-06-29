FROM golang:latest
ADD . /
WORKDIR /
RUN make
CMD ["./bin/CrashDragon"]
