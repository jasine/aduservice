# Set the base image to use to Ubuntu
FROM 192.168.5.46:5000/ubuntu:14.04

# Set the file maintainer (your name - the file's author)
MAINTAINER yadda Deepglint <zizhenyan@deepglint.com>

# Add files to the image
RUN mkdir root/aduservice

#add
ADD aduservice.linux root/aduservice/aduservice.linux

RUN mkdir /data
RUN mkdir /data/adu


ADD auth /data/adu/

WORKDIR root/aduservice
