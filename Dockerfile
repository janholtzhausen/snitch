FROM alpine
COPY ./snitch /usr/bin/
CMD ["snitch"]
