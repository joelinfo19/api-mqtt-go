#FROM surnet/alpine-wkhtmltopdf:3.16.2-0.12.6-full as wkhtmltopdf
FROM golang:1.17-alpine

#MAINTAINER Bengie Nick <bengie.serrano@smartc.pe>

# wkhtmltopdf install dependencies


# wkhtmltopdf copy bins from ext image
#COPY --from=wkhtmltopdf /bin/wkhtmltopdf /bin/libwkhtmltox.so /bin/

# Install required packages
#RUN apk update
#
#RUN apk add dmidecode
#RUN apk add ca-certificates
RUN apk add dos2unix
#RUN apk add openssl
#RUN apk --no-cache add tzdata
#
#RUN mkdir /home/storage
#RUN mkdir /home/docs
ENV URL_PATH "http://localhost:8080"
ENV BROKER "broker.mqttdashboard.com:1883"
COPY main /home/main


EXPOSE 80

WORKDIR /home

ADD cicd/run.sh /home/run.sh
RUN ["chmod", "+x", "/home/run.sh"]
RUN dos2unix /home/run.sh

ENTRYPOINT ["/home/run.sh"]
