FROM postgres:latest

RUN apt-get update
RUN apt-get install -y curl less libpq-dev python-dev python-pip
RUN pip install pgcli

WORKDIR /root
RUN curl https://raw.githubusercontent.com/rubygems/rubygems.org/master/script/load-pg-dump --output load-pg-dump
RUN chmod +x ./load-pg-dump

ENV PAGER="less -S"
ENTRYPOINT /bin/bash
# /docker-entrypoint.sh postgres &
# ./load-pg-dump -c
# pgcli -u postgres -d rubygems