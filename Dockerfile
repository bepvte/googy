FROM golang:1.11-rc-alpine

RUN apk add --no-cache tesseract-ocr tesseract-ocr-dev leptonica-dev build-base figlet

WORKDIR /go/src/googy
COPY . .

ENV GO111MODULE=on

# RUN go get
RUN go install -v -mod=vendor ./...
RUN strip /go/bin/googy
RUN apk add --no-cache upx && upx /go/bin/googy && apk del upx

CMD ["googy"]
