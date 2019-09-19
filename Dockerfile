FROM ubuntu
RUN apt-get update && apt-get install pwgen
RUN mkdir /app
WORKDIR /app
COPY --from=compile /go/src/github.com/ESCah/go-against-humanity/server /app/server
COPY --from=compile /go/src/github.com/ESCah/go-against-humanity/app/views /app/app/views
COPY --from=compile /go/src/github.com/ESCah/go-against-humanity/public /app/public
COPY ./start.sh /app
ENTRYPOINT ["/app/start.sh"]
