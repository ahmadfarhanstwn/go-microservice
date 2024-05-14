FROM alpine:latest

RUN mkdir /app

COPY mailerServiceApp /app
COPY template /templates

CMD [ "/app/mailerServiceApp" ]