REPOSITORY=artisan
DISTRIBUTION=all
SECTION=stable

curl -u "admin:admin" \
  -H 'accept: plain/text' \
  -H "Content-Type: multipart/form-data" \
  -F "package=@artisan_0.4.9-1_amd64.deb" \
  "http://localhost:8085/apt/repository/${REPOSITORY}/dist/${DISTRIBUTION}/section/${SECTION}"

curl -u "admin:admin" \
  -H 'accept: plain/text' \
  -H "Content-Type: multipart/form-data" \
  -F "package=@artisan_0.4.9-1_arm64.deb" \
  "http://localhost:8085/apt/repository/${REPOSITORY}/dist/${DISTRIBUTION}/section/${SECTION}"

curl -u "admin:admin" \
  -H 'accept: plain/text' \
  -H "Content-Type: multipart/form-data" \
  -F "package=@artisan_0.4.10-1_amd64.deb" \
  "http://localhost:8085/apt/repository/${REPOSITORY}/dist/${DISTRIBUTION}/section/${SECTION}"

curl -u "admin:admin" \
  -H 'accept: plain/text' \
  -H "Content-Type: multipart/form-data" \
  -F "package=@artisan_0.4.10-1_arm64.deb" \
  "http://localhost:8085/apt/repository/${REPOSITORY}/dist/${DISTRIBUTION}/section/${SECTION}"