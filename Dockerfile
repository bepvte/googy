FROM golang:alpine

RUN apk add --no-cache tesseract-ocr tesseract-ocr-dev tesseract-ocr-data-ara tesseract-ocr-data-jpn tesseract-ocr-data-rus tesseract-ocr-data-kor tesseract-ocr-data-ell leptonica-dev build-base figlet

WORKDIR /go/src/googy
COPY . .

ENV GO111MODULE=on

# RUN go get
RUN go install -v -mod=vendor ./...
RUN strip /go/bin/googy
RUN apk add --no-cache upx && upx /go/bin/googy && apk del upx

CMD ["googy"]
