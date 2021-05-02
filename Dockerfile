FROM golang

RUN mkdir /bot
WORKDIR /bot
COPY ./ ./
RUN go build .
RUN chmod +x discord-bot
ENTRYPOINT ["./discord-bot"] 