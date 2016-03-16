FROM busybox:glibc
MAINTAINER Beldur

ADD _output/rsvgd_linux_amd64 /rsvgd

CMD ["/rsvgd"]
