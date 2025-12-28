# Copyright (c) 2023, WSO2 LLC. (https://www.wso2.com/)
# All Rights Reserved.
#
# WSO2 LLC. licenses this file to you under the Apache License,
# Version 2.0 (the "License"); you may not use this file except
# in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied. See the License for the
# specific language governing permissions and limitations
# under the License.

FROM golang:1.22.4-alpine AS build-env

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN addgroup -g 10014 choreo \
  && adduser --disabled-password --no-create-home --uid 10014 --ingroup choreo choreouser

COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/app -buildvcs=false

FROM alpine
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build-env /go/bin/app /go/bin/app
COPY --from=build-env /app/templates ./templates

USER 10014
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/go/bin/app"]
