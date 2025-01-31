FROM golang:1.22.2-alpine

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN go build .


FROM alpine

ENV SUBTITLES_EXTRACTOR_LIBRARIES /media
ENV SUBTITLES_EXTRACTOR_SLEEP 1
ENV SUBTITLES_EXTRACTOR_FORCED_ONLY 0
ENV SUBTITLES_EXTRACTOR_SKIP_SRT 0
ENV SUBTITLES_EXTRACTOR_STRIP_FORMATTING 0
ENV SUBTITLES_EXTRACTOR_LANGUAGES *
ENV SUBTITLES_EXTRACTOR_DATA_DIR /data

VOLUME /data
VOLUME /media

RUN apk update
RUN apk upgrade
RUN apk add --no-cache ffmpeg

COPY --from=0 /src/go-subtitles-extractor /usr/bin

CMD ["go-subtitles-extractor"]
