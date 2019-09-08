FROM golang as builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o /chess_ladder

FROM ubuntu

ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /
RUN chmod +x /wait-for-it.sh

COPY --from=builder /chess_ladder /
EXPOSE 8000
CMD /bin/bash -c "/wait-for-it.sh db:5432 && /chess_ladder"
