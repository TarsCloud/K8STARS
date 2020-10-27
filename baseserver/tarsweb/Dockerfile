FROM centos/nodejs-8-centos7:8

COPY TarsWeb start.sh default.env /web/
ENV PATH /opt/rh/rh-nodejs8/root/usr/bin:/opt/app-root/src/node_modules/.bin/:/opt/app-root/src/.npm-global/bin/:/opt/app-root/src/bin:/opt/app-root/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
USER root

RUN ln -nsf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN npm config set registry http://mirrors.cloud.tencent.com/npm/
RUN cd /web && npm install && mkdir -p /web/files

CMD cd /web && sh /web/start.sh
