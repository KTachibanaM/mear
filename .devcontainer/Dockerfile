FROM mcr.microsoft.com/devcontainers/go:0-1.20-bullseye

WORKDIR /root
RUN curl https://dl.min.io/client/mc/release/linux-amd64/mc.RELEASE.2023-06-06T13-48-56Z -o /usr/bin/mc
RUN chmod +x /usr/bin/mc

RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
RUN unzip awscliv2.zip
RUN ./aws/install -i /usr/local/aws-cli -b /usr/local/bin

RUN echo 'deb [trusted=yes] https://repo.goreleaser.com/apt/ /' | sudo tee /etc/apt/sources.list.d/goreleaser.list
RUN sudo apt update
RUN sudo apt install -y goreleaser ntp
RUN sudo service ntp restart
