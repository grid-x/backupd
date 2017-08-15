FROM ubuntu:16.04


RUN apt-get update && apt-get install -y curl apt-transport-https && \
            curl -sL https://repos.influxdata.com/influxdb.key | apt-key add - && \
            echo "deb https://repos.influxdata.com/ubuntu xenial stable" | tee /etc/apt/sources.list.d/influxdb.list && \
            apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 0C49F3730359A14518585931BC711F9BA15703C6 && \
            echo "deb [ arch=amd64,arm64 ] http://repo.mongodb.org/apt/ubuntu xenial/mongodb-org/3.4 multiverse" | tee /etc/apt/sources.list.d/mongodb-org-3.4.list && \
            apt-get update && apt-get install -y influxdb mongodb-org && \
            rm -rf /var/lib/apt/lists/* 

ADD bin/server.linux /usr/local/bin/server

ENTRYPOINT ["/usr/local/bin/server"]
