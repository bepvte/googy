FROM golang:1.10-alpine3.7

RUN apk add --no-cache tesseract-ocr tesseract-ocr-dev leptonica-dev build-base figlet

WORKDIR /go/src/googy
COPY . .

RUN go install -v ./...

CMD ["googy"]
