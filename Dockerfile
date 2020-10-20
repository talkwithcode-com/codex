FROM golang

RUN apt-get update -y
RUN apt-get install curl vim sudo build-essential -y
RUN curl -sL https://deb.nodesource.com/setup_12.x | bash
RUN apt-get install nodejs -y
RUN node -v

EXPOSE 7000

RUN python3 -V
RUN node -v
RUN go version

# run each time it builds
ENV CODEX_PORT=7000
ARG CACHE_DATE=2020-01-01 
# RUN go get github.com/talkwithcode-com/codex
# RUN go install github.com/talkwithcode-com/codex
# RUN which codex
COPY ./codex .

RUN ls

CMD ["./codex"]