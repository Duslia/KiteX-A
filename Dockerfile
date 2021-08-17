FROM iron/go
WORKDIR /opt/tiger/toutiao/app/service_a

ADD service_a /opt/tiger/toutiao/app/service_a

CMD [ "./service_a" ]
