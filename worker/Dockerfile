FROM fedora:42

RUN dnf -y update && dnf -y install \
git \
make \
wget \
gcc \
gcc-c++ \
pkg-config \
libcap-devel \
systemd-devel \
java-25-openjdk-devel.x86_64 \
python3 && \
dnf clean all && rm -rf /var/cache/dnf

# installing isolate
RUN git clone https://github.com/ioi/isolate.git /opt/isolate && \
cd /opt/isolate && \
make install

# installing golang 1.25
RUN wget https://go.dev/dl/go1.25.3.linux-amd64.tar.gz && \
 rm -rf /usr/local/go && tar -C /usr/local -xzf go1.25.3.linux-amd64.tar.gz && \
 rm -rf go1.25.3.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

ENTRYPOINT [ "go" ]

CMD [ "run", "." ]