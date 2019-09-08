FROM golang:1.12.0-alpine3.9
ADD config config
ADD data data
RUN ls -la
ADD main .

EXPOSE 80

CMD ["./main"]