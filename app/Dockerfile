FROM image-registry.openshift-image-registry.svc:5000/openshift/nodejs:20-ubi9

WORKDIR /opt/app-root/src

COPY package.json .
RUN npm install --production

COPY server.js .

EXPOSE 8080

CMD ["node", "server.js"]