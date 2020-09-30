FROM golang:latest
RUN mkdir /rest_api
ADD . /rest_api/
WORKDIR /rest_api
RUN make
CMD ["/rest_api/cmd/apisever/main"]