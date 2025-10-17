FROM golang:1.25-alpine

RUN apk add --no-cache ffmpeg python3 py3-pip bash

# virtualenv + yt-dlp
RUN python3 -m venv /opt/venv && \
    /opt/venv/bin/pip install --upgrade pip && \
    /opt/venv/bin/pip install yt-dlp && \
    ln -s /opt/venv/bin/yt-dlp /usr/bin/yt-dlp

WORKDIR /app
COPY . .

RUN go build -o server ./cmd/api

EXPOSE 8080

CMD ["./server"]
