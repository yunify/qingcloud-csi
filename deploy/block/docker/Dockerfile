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


FROM ubuntu:16.04
MAINTAINER "Yunify"

RUN apt-get update && \
    apt-get install -y ca-certificates ca-certificates-java && \
    apt-get install -y xfsprogs
COPY blockplugin /blockplugin
RUN chmod +x /blockplugin && \
    mkdir -p /var/log/blockplugin

ENTRYPOINT ["/blockplugin"]
