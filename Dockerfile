FROM scratch
COPY ./snitch /
COPY ./snitch.conf /
CMD ["/snitch"]
