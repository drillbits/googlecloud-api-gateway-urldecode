#!/bin/bash

# Copyright 2021 drillbits
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -euo pipefail
scriptdir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
[[ -n "${DEBUG:-}" ]] && set -x
cd ${scriptdir} && cd ../

destdir=$(pwd)/genproto/api

if [ -d ${destdir} ]; then
  rm -rf ${destdir}
fi
mkdir -p ${destdir}

protoc \
  -I/usr/include \
  -I$(pwd) \
  --descriptor_set_out=$(pwd)/api_descriptor.pb \
  --include_imports \
  --include_source_info \
  --go_out=${destdir} \
  --go_opt=module=github.com/drillbits/googlecloud-api-gateway-urldecode/genproto/api \
  --go-grpc_out=${destdir} \
  --go-grpc_opt=require_unimplemented_servers=false \
  --go-grpc_opt=module=github.com/drillbits/googlecloud-api-gateway-urldecode/genproto/api \
  $(pwd)/echo.proto
