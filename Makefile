NAME=lfs-s3-agent
VERSION=0.0.1-SNAPSHOT-$(CIRCLE_BUILD_NUM)

all: package deploy

.PHONY: package
package:
	fpm -s dir -t deb -n $(NAME) -v $(VERSION) `which lfs-s3-agent`
	fpm -s dir -t rpm -n $(NAME) -v $(VERSION) `which lfs-s3-agent`

deploy:
	ls -la
	# curl -T deb -ualdrinleal:$API_KEY https://api.bintray.com/content/aldrinleal/deb/lfs-s3-agent/$VERSION/deb
	# curl -T pm -ualdrinleal:$API_KEY https://api.bintray.com/content/aldrinleal/deb/lfs-s3-agent/$VERSION/deb
