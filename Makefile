NAME=lfs-s3-transfer-agent
VERSION=0.0.1-SNAPSHOT-$(CIRCLE_BUILD_NUM)

all: package deploy

clean:
	rm *.deb *.rpm || true

.PHONY: package
package:
	fpm -s dir -t deb -n $(NAME) -v $(VERSION) -C ~/go bin/lfs-s3-transfer-agent=/usr/bin/lfs-s3-transfer-agent
	fpm -s dir -t rpm -n $(NAME) -v $(VERSION) -C ~/go bin/lfs-s3-transfer-agent=/usr/bin/lfs-s3-transfer-agent

deploy: package
	~/go/bin/lfs-s3-transfer-agent --version
	curl -T *.deb -ualdrinleal:$(BINTRAY_APIKEY) 'https://api.bintray.com/content/ingenieux/lfs-deb/$(NAME)/$(VERSION)/$(shell basename *.deb);deb_distribution=any;deb_component=main;deb_architecture=amd64;publish=1'
	curl -T *.rpm -ualdrinleal:$(BINTRAY_APIKEY) 'https://api.bintray.com/content/ingenieux/lfs-rpm/$(NAME)/$(VERSION)/$(shell basename *.rpm);publish=1'
