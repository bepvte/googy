FROM golang:alpine

RUN apk add --no-cache tesseract-ocr tesseract-ocr-data-ara tesseract-ocr-data-jpn tesseract-ocr-data-rus tesseract-ocr-data-kor tesseract-ocr-data-ell figlet 

RUN

WORKDIR /go/src/googy
COPY . .

ENV GO111MODULE=on

RUN apk add --no-cache --virtual .build-deps upx build-base leptonica-dev tesseract-ocr-dev \
 && go install -v -mod=vendor ./... \
 && strip /go/bin/googy \
 && upx /go/bin/googy \
 && apk del .build-deps

CMD ["googy"]

LABEL IMPORTANT=yes
