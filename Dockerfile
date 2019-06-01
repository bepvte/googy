FROM golang:alpine

RUN apk add --no-cache tesseract-ocr tesseract-ocr-data-ara tesseract-ocr-data-jpn tesseract-ocr-data-rus tesseract-ocr-data-kor tesseract-ocr-data-ell figlet

RUN apk add --no-cache --virtual .build-deps upx build-base leptonica-dev tesseract-ocr-dev

WORKDIR /go/src/googy
COPY . .

ENV GO111MODULE=on
RUN go install -v -mod=vendor ./... \
 && strip /go/bin/googy \
 && upx /go/bin/googy
RUN apk del .build-deps

CMD ["googy"]

LABEL IMPORTANT=yes
