FROM golang:1.10-stretch

RUN apt-get update && apt-get install -y libtesseract-dev libleptonica-dev tesseract-ocr-eng && rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/googy
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["googy"]