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

function usage {
  echo "Usage: $0 [-p project] [-r region] [-v version]" 1>&2
  exit 1
}

while getopts "p:r:v:h" OPT
do
  case $OPT in
    p) project=${OPTARG}
      ;;
    r) region=${OPTARG}
      ;;
    v) config_version=${OPTARG}
      ;;
    h) usage
      ;;
  esac
done

project=${project:-}
region=${region:-}

shift $(($OPTIND - 1))

if [ -z "${project}" ]; then
  echo -n "Google Cloud Project: "
  read project
fi

if [ -z "${region}" ]; then
  echo -n "Google Cloud Region: "
  read region
fi

docker_repository="gcr.io/${project}/apigateway-urldecode-backend"
docker_tag=${BACKEND_DOCKER_TAG:-$(git rev-parse --short HEAD)}

env_var_array=()
env_var_array+=( "GOOGLE_CLOUD_PROJECT=${project}" )
env_var_array+=( "REGION=${region}" )
env_vars="$(IFS=","; echo "${env_var_array[*]}")"

echo "Docker build"
docker build --file backend/Dockerfile --target release --tag ${docker_repository}:${docker_tag} .

echo "Docker push"
docker push ${docker_repository}:${docker_tag}

echo "Deploy backend to Cloud Run"
echo "- Project: ${project}"
echo "- Image: ${docker_repository}:${docker_tag}"
gcloud run deploy apigateway-urldecode-backend \
  --project ${project} \
  --platform managed \
  --allow-unauthenticated \
  --image ${docker_repository}:${docker_tag} \
  --min-instances 1 \
  --region ${region} \
  --set-env-vars ${env_vars}
