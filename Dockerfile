FROM golang:1.12

# RUN  tesseract-ocr tesseract-ocr-data-ara tesseract-ocr-data-jpn tesseract-ocr-data-rus tesseract-ocr-data-kor tesseract-ocr-data-ell figlet

RUN apt-get update && apt-get install -y upx libtesseract-dev libleptonica-dev

WORKDIR /go/src/googy
# COPY . .

# ENV GO111MODULE=on
# RUN go build -ldflags="-s -w" -v -mod=vendor \
#  && upx googy

# CMD ["googy"]

# LABEL IMPORTANT=yes
