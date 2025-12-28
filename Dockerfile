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

FROM golang:1.22-alpine AS build

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/web-app .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /app/web-app ./web-app
COPY --from=build /app/templates ./templates

# Create a user with a known UID/GID (10014) within the range 10000-20000 for Choreo.
RUN adduser -D -H -s /sbin/nologin -u 10014 choreo \
  && chown -R 10014:10014 /app

USER 10014
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app/web-app"]
