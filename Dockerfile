# Sử dụng một hình ảnh Ubuntu làm cơ sở
FROM ubuntu:latest

RUN apt update && apt install -y openssh-server sudo sed curl golang-go build-essential cmake pkg-config libblas-dev liblapack-dev libatlas-base-dev libjpeg-dev
RUN pkg-config --libs --cflags libjpeg

RUN sed -i "s/UsePAM yes/UsePAM no/g" /etc/ssh/sshd_config && sed -i "s/#PermitRootLogin prohibit-password/PermitRootLogin yes/g" /etc/ssh/sshd_config
RUN echo "root:123456" | chpasswd
RUN echo "root   ALL=(ALL)       ALL" >> /etc/sudoers

RUN mkdir /run/sshd
EXPOSE 22

ENV GOPATH=/go
WORKDIR /root
RUN mkdir /go
WORKDIR /root/go
RUN go mod init go
RUN go get -u gocv.io/x/gocv
WORKDIR /go/pkg/mod/gocv.io/x/gocv@v0.35.0
RUN make install && make clean && pkg-config --cflags --libs opencv4
RUN go install gocv.io/x/gocv
RUN go clean -cache

WORKDIR /root
RUN wget http://dlib.net/files/dlib-19.24.tar.bz2 && tar xvf dlib-19.24.tar.bz2
WORKDIR /root/dlib-19.24
RUN rm -rf build && mkdir build
WORKDIR /root/dlib-19.24/build
RUN cmake .. && cmake --build . --config Release && sudo make install && sudo ldconfig
RUN rm -rf /root/dlib-19.24.tar.bz2
RUN pkg-config --libs --cflags dlib-1

WORKDIR /root
RUN rm -rf /root/go
COPY ./app/go.mod .
RUN go mod download
COPY ./app .
RUN go mod tidy
RUN go build -o main .

RUN rm -rf /root/dlib-19.24

CMD ["/usr/sbin/sshd", "-D"]
