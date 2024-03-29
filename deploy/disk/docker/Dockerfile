# +-------------------------------------------------------------------------
# | Copyright (C) 2018 Yunify, Inc.
# +-------------------------------------------------------------------------
# | Licensed under the Apache License, Version 2.0 (the "License");
# | you may not use this work except in compliance with the License.
# | You may obtain a copy of the License in the LICENSE file, or at:
# |
# | http://www.apache.org/licenses/LICENSE-2.0
# |
# | Unless required by applicable law or agreed to in writing, software
# | distributed under the License is distributed on an "AS IS" BASIS,
# | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# | See the License for the specific language governing permissions and
# | limitations under the License.
# +-------------------------------------------------------------------------

FROM golang:1.19-alpine as builder
WORKDIR /qingcloud-csi
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -mod=vendor  -ldflags "-s -w" -o  _output/qingcloud-disk-csi-driver ./cmd/disk

FROM debian:stretch
LABEL maintainers="Yunify"
LABEL description="QingCloud CSI plugin"

RUN apt-get update && \
    apt-get install -y util-linux e2fsprogs xfsprogs mount ca-certificates udev && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /qingcloud-csi/_output/qingcloud-disk-csi-driver /qingcloud-disk-csi-driver
RUN chmod +x /qingcloud-disk-csi-driver && \
    mkdir -p /var/log/qingcloud-disk-csi-driver
ENTRYPOINT ["/qingcloud-disk-csi-driver"]
