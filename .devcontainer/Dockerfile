FROM mcr.microsoft.com/devcontainers/go:0-1.20-bullseye

WORKDIR /root
RUN curl https://dl.min.io/client/mc/release/linux-amd64/mc -o /usr/bin/mc
RUN chmod +x /usr/bin/mc

RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
RUN unzip awscliv2.zip
RUN ./aws/install -i /usr/local/aws-cli -b /usr/local/bin