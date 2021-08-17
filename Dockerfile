FROM iron/go
WORKDIR /opt/tiger/toutiao/app/service_c

ADD kitex.service.c /opt/tiger/toutiao/app/service_c

CMD [ "./kitex.service.c" ]
