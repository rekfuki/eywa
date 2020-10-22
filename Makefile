.PHONY: setup-repo
setup-repo: guard-USERNAME guard-PASSWORD guard-TLS_CERT_DIR
	htpasswd -Bbn ${USERNAME} ${PASSWORD} > htpasswd
	docker run -d \
	  -p 5000:5000 \
	  --restart=always \
	  --name registry \
	  -v $(shell pwd)/htpasswd:/htpasswd \
	  -e "REGISTRY_AUTH=htpasswd" \
	  -e "REGISTRY_AUTH_HTPASSWD_REALM=Registry Realm" \
	  -e REGISTRY_AUTH_HTPASSWD_PATH=/htpasswd \
	  -v ${TLS_CERT_DIR}/fullchain.pem:${TLS_CERT_DIR}/fullchain.pem \
	  -e REGISTRY_HTTP_TLS_CERTIFICATE=${TLS_CERT_DIR}/fullchain.pem \
	  -v ${TLS_CERT_DIR}/privkey.pem:${TLS_CERT_DIR}/privkey.pem \
	  -e REGISTRY_HTTP_TLS_KEY=${TLS_CERT_DIR}/privkey.pem \
	  registry:2

.PHONY: delete-repo
delete-repo:
	docker container stop registry
	docker container rm registry

guard-%:
	if [ -z '${${*}}' ]; then echo 'Environment variable $* not set' && exit 1; fi
