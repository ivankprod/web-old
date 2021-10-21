FROM ubuntu
RUN apt-get update
RUN apt-get install -y ca-certificates
RUN update-ca-certificates
