FROM postgres:latest

RUN apt-get update
RUN apt-get install -y wget python-pip less
RUN pip install pgcli

WORKDIR /root
RUN wget https://s3-us-west-2.amazonaws.com/rubygems-dumps/production/public_postgresql/2019.03.11.21.21.01/public_postgresql.tar
RUN wget https://raw.githubusercontent.com/rubygems/rubygems.org/master/script/load-pg-dump
RUN chmod +x ./load-pg-dump

ENV PAGER="less -S"
ENTRYPOINT /bin/bash
# /docker-entrypoint.sh postgres &
# ./load-pg-dump ./public_postgresql.tar
# pgcli -u postgres -d rubygems