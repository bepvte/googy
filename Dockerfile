FROM golang:1.9-alpine3.7

RUN apk add --no-cache tesseract-ocr tesseract-ocr-dev leptonica-dev build-base

WORKDIR /go/src/googy
COPY . .

RUN go install -v ./...

CMD ["googy"]
