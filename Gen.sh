#!/usr/bin/bash

if [ $# -lt 2 ]; then 
	echo "$0 <prefix> <service-name> [<http-port>]"
	exit 1;
fi

PREFIX=$1
SERVICE=$2
PROJECT_NAME=flowork
PROJECT_HOME=PIKABU_HOME
DOCKER_HUB_PASSWORD=KFXB3UVFpB7tmnC

GO_PATH=/go/src
PACKAGE_BASE=${PREFIX}-${SERVICE}
APP=app_${PREFIX}_${SERVICE}
BUILD_SCRIPT=Build.sh
UPDATE_SCRIPT=Update.sh
RUN_SCRIPT=Run.sh
BUILD_DEV_IMAGE_SCRIPT=BuildDevImage.sh
BUILD_DEP_IMAGE_SCRIPT=BuildDepImage.sh
DEPLOY_DEV_IMAGE_SCRIPT=DeployDevImage.sh
DEPLOY_DEP_IMAGE_SCRIPT=DeployDepImage.sh
BUILD_IMAGE=${PROJECT_NAME}'/golang:1.11.5-dev'
HUB_BUILD_IMAGE=${PROJECT_NAME}/golang:1.11.5-dev
DEV_IMAGE=${PROJECT_NAME}'/'${PREFIX}-${SERVICE}-dev':latest'
DEPLOY_BASE_IMAGE=${DEV_IMAGE}
DEV_NETWORK=develop_network
DEP_NETWORK=deploy_network
DEP_VOLUME=${PREFIX}_${SERVICE}_volume
DEP_IMAGE=${PROJECT_NAME}'/'${PREFIX}-${SERVICE}':latest'
CPU_LIMIT=1.0
DEP_CPU_LIMIT=0.1
MEMORY_SIZE=1000M
DEP_MEMORY_SIZE=100M

# start

cat > KitGen.sh<< EOF
#!/bin/sh
kit generate service ${SERVICE} --gorilla
kit generate service ${SERVICE} --dmw
kit generate service ${SERVICE} --gorilla
kit generate service ${SERVICE} -t grpc
kit generate docker
EOF
chmod +x KitGen.sh
sh KitGen.sh

cat > KitGen.sh<< EOF
#!/bin/sh
kit generate service ${SERVICE} --gorilla
kit generate service ${SERVICE} --dmw
kit generate service ${SERVICE} --gorilla
kit generate service ${SERVICE} -t grpc
EOF

if [ $# -eq 3 ]; then
	PORT=$3
	sed -i "s/- 8800:8081/- ${PORT}:8081/g" docker-compose.yml
	sed -i "s/ - 8801:8081/#- ${PORT}:8081/g" docker-compose.yml
	sed -i "s/ - 8801:8082/#- 8801:8082/g" docker-compose.yml
	sed -i "s/ - 8802:8082/#- 8802:8082/g" docker-compose.yml
else
	sed -i "s/ ports:/#ports:/g" docker-compose.yml
	sed -i "s/ - 8/#- 8/g" docker-compose.yml
fi

sed -i "s/ image/#image/g" docker-compose.yml
sed -i "s/\"2\"/\"3\.2\"/g" docker-compose.yml
sed -i "s/\<${SERVICE}:/${PREFIX}_${SERVICE}:/g" docker-compose.yml
sed -i "s/container_name: ${SERVICE}/container_name: ${PREFIX}_${SERVICE}/g" docker-compose.yml
cat >> docker-compose.yml << EOF
    image: ${DEV_IMAGE}
EOF

cat > docker-stack.develop.yml<< EOF
version: "3.7"
services:
  ${PREFIX}_${SERVICE}:
    image: ${DEV_IMAGE}
    volumes:
      - type: bind
        source: .
        target: /go/src/${PACKAGE_BASE}
    deploy:
      replicas: 1
      restart_policy:
        condition: any
      resources:
        limits:
         cpus: "${CPU_LIMIT}"
         memory: ${MEMORY_SIZE}
EOF
if [ $# -eq 3 ]; then
	echo '    ports:' >> docker-stack.develop.yml
	echo '      - '${PORT}':8081' >> docker-stack.develop.yml
fi
cat >> docker-stack.develop.yml<< EOF
    networks:
      - ${DEV_NETWORK}
networks:
  ${DEV_NETWORK}:
EOF

cp docker-stack.develop.yml docker-stack.deploy.yml
sed -i "s/${DEV_NETWORK}/${DEP_NETWORK}/g" docker-stack.deploy.yml
sed -i "s/${PREFIX}-${SERVICE}-dev/${PREFIX}-${SERVICE}/g" docker-stack.deploy.yml
sed -i "s/- type: bind/- type: volume/g" docker-stack.deploy.yml
sed -i "s/ source: \./ source: ${DEP_VOLUME}/g" docker-stack.deploy.yml
sed -i "s/ target: \/go\/src\/${PACKAGE_BASE}/ target: \/var\/local/g" docker-stack.deploy.yml
sed -i "s/cpus: \"${CPU_LIMIT}\"/cpus: \"${DEP_CPU_LIMIT}\"/g" docker-stack.deploy.yml
sed -i "s/memory: "${MEMORY_SIZE}"/memory: ${DEP_MEMORY_SIZE}/g" docker-stack.deploy.yml

cat >> docker-stack.deploy.yml << EOF
volumes:
  ${DEP_VOLUME}:
    external: true
EOF

cat > ${BUILD_SCRIPT} << EOF
#!/bin/sh
watcher -cmd="sh ${UPDATE_SCRIPT}" -recursive -pipe=true -list ./${SERVICE} &
canthefason_watcher -run ${PACKAGE_BASE}/${SERVICE}/cmd -watch ${PACKAGE_BASE}
EOF
chmod +x ${BUILD_SCRIPT}

cat > ${UPDATE_SCRIPT} << EOF
#!/bin/bash
echo update > DontRemoveMe.go
EOF
chmod +x ${UPDATE_SCRIPT}

cat > ${RUN_SCRIPT} << EOF
#!/bin/bash
${PROJECT_HOME}=${PWD}
go run ${SERVICE}/cmd/main.go
EOF
chmod +x ${RUN_SCRIPT}

cat > Dockerfile << EOF
FROM ${HUB_BUILD_IMAGE} as build
WORKDIR ${GO_PATH}/${PACKAGE_BASE}
ADD . .
#RUN apk add --no-cache bash git openssh
#RUN dep init -v -no-examples
RUN go build -o ${APP} ${PACKAGE_BASE}/${SERVICE}/cmd

FROM alpine:3.9
ENV ${PROJECT_HOME} /var/local
WORKDIR /var/local/${PACKAGE_BASE}/config
COPY --from=build ${GO_PATH}/${PACKAGE_BASE}/config .
WORKDIR /var/local/${PACKAGE_BASE}/img
COPY --from=build ${GO_PATH}/${PACKAGE_BASE}/img .
WORKDIR /usr/local/bin
COPY --from=build ${GO_PATH}/${PACKAGE_BASE}/${APP} /usr/local/bin/${APP}

ENTRYPOINT ["${APP}"]
EOF

cat > ${SERVICE}/Dockerfile << EOF
FROM ${HUB_BUILD_IMAGE}
ENV ${PROJECT_HOME} ${GO_PATH}
WORKDIR ${GO_PATH}/${PACKAGE_BASE}
ADD . .
#RUN dep init -v -no-examples
#RUN dep ensure -v -vendor-only
#RUN go get -v ${PACKAGE_BASE}/${SERVICE}/cmd
RUN chmod +x ${BUILD_SCRIPT}
ENTRYPOINT sh ${BUILD_SCRIPT}
EOF

cat > ${BUILD_DEV_IMAGE_SCRIPT} << EOF
#!/bin/sh
EXE=docker-compose
unameOut="\$(uname -s)"
if [ `echo \${unameOut} | grep -ciP '(cygwin|mingw)'` -ne 0 ]; then
	EXE=\${EXE}.exe
fi

\${EXE} build --no-cache
docker push ${DEV_IMAGE}
EOF
chmod +x ${BUILD_DEV_IMAGE_SCRIPT}

cat > ${BUILD_DEP_IMAGE_SCRIPT} << EOF
#!/bin/sh
EXE=docker
unameOut="\$(uname -s)"
if [ `echo \${unameOut} | grep -ciP '(cygwin|mingw)'` -ne 0 ]; then
	EXE=\${EXE}.exe
fi

\${EXE} build --no-cache -t ${DEP_IMAGE} .
docker push ${DEP_IMAGE}
EOF
chmod +x ${BUILD_DEP_IMAGE_SCRIPT}

cat > ${DEPLOY_DEV_IMAGE_SCRIPT} << EOF
#!/bin/sh
docker service rm ${PROJECT_NAME}_${PREFIX}_${SERVICE} 2>/dev/null
docker stack deploy -c docker-stack.develop.yml ${PROJECT_NAME}
EOF
chmod +x ${DEPLOY_DEV_IMAGE_SCRIPT}

cat > ${DEPLOY_DEP_IMAGE_SCRIPT} << EOF
#!/bin/sh
docker service rm ${PROJECT_NAME}_${PREFIX}_${SERVICE} 2>/dev/null
docker stack deploy -c docker-stack.deploy.yml ${PROJECT_NAME}
EOF
chmod +x ${DEPLOY_DEP_IMAGE_SCRIPT}

cat > GoGet.sh << EOF
#!/bin/sh
if [ \$# -eq 1 ]; then
        GOPATH=\$1
fi

cd ${SERVICE}/cmd
go get -v 
EOF


cat > .gitlab-ci.yml << EOF
variables:
  PACKAGE_PATH: /go/src/${PACKAGE_BASE}
  DOCKER_HUB_REGISTRY_PATH: ${PROJECT_NAME}/${PACKAGE_BASE}

stages:
  - build
  - deploy

build:
  tags:
    - dev
  only:
    - master
  stage: build
  image: docker:18.09.2
  services:
    - docker:dind
  script:
    - docker login -u ${PROJECT_NAME} -p ${DOCKER_HUB_PASSWORD}
    - docker build -t \$DOCKER_HUB_REGISTRY_PATH:latest .
    - docker push \$DOCKER_HUB_REGISTRY_PATH:latest

deploy:
  tags:
    - dep
  only:
    - master
  dependencies:
    - build
  stage: deploy
  script:
    - docker login -u ${PROJECT_NAME} -p ${DOCKER_HUB_PASSWORD}
    - docker pull \$DOCKER_HUB_REGISTRY_PATH:latest
    - docker stack deploy -c docker-stack.deploy.yml ${PROJECT_NAME}
EOF

mkdir config
cat > config/README.md << EOF
This folder is for Configuration
EOF
mkdir img
cat > img/README.md << EOF
This folder is for Image
EOF
