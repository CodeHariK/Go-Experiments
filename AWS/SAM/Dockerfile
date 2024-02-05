FROM amazonlinux:2023

# Update package manager metadata
RUN yum update -y

# Install dependencies
RUN yum install -y unzip wget tar gzip

# Install AWS CLI v2
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-aarch64.zip" -o "awscliv2.zip" && \
    unzip awscliv2.zip && \
    ./aws/install -i /usr/local/aws-cli -b /usr/local/bin && \
    rm -rf awscliv2.zip aws

# Install AWS SAM CLI
RUN curl -L "https://github.com/aws/aws-sam-cli/releases/latest/download/aws-sam-cli-linux-arm64.zip" -o "aws-sam-cli-linux-arm64.zip" && \
    unzip aws-sam-cli-linux-arm64.zip -d sam-installation && \
    ./sam-installation/install

# Download and install Go
RUN wget https://go.dev/dl/go1.21.6.linux-arm64.tar.gz && \
    tar -C /usr/local -xzf go1.21.6.linux-arm64.tar.gz && \
    rm go1.21.6.linux-arm64.tar.gz

# Add Go binaries to the PATH
ENV PATH=$PATH:/usr/local/go/bin

# Verify the Go installation
RUN go version

WORKDIR /app

COPY . .

# Verify installations
RUN aws --version && sam --version

CMD ["bash"]
