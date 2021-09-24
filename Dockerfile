FROM golang:1.17

ARG VERSION=1.0.7

RUN apt-get update && \
    apt-get -y --no-install-recommends install ca-certificates curl unzip && \
    apt-get clean && apt-get autoclean && apt-get -y autoremove --purge && \
    rm -rf /var/lib/apt/lists/* /usr/share/doc /root/.cache/ && \
    # Install Terraform
    curl -s https://releases.hashicorp.com/terraform/${VERSION}/terraform_${VERSION}_linux_$(go env GOARCH).zip -o terraform.zip && \
    unzip -q terraform.zip && \
    mv terraform /usr/local/bin/terraform && \
    rm terraform.zip

COPY --from=goreleaser/goreleaser /usr/local/bin/goreleaser /usr/local/bin/goreleaser
COPY --from=golangci/golangci-lint /usr/bin/golangci-lint /usr/local/bin/golangci-lint
