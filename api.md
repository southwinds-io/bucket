<img src="https://github.com/southwinds-io/bucket/raw/main/bucket.png" width="100" align="right"/>

# Bucket API 

## Debian Packages

### Uploading a Package

curl -u user:pwd -H "Content-Type: multipart/form-data" -F "package=@./path/to/package.deb" "http(s)://**HOST:PORT**/debian/repo/**REPO-NAME**/dist/**DIST-NAME**/section/**SECTION-NAME**"

Example:

uploads package test_arm64.deb to the artisan repository, all metadata is read from the package

```bash
$ curl -u "admin:admin" \
    -H "Content-Type: multipart/form-data" \
    -F "package=@./internal/deb/test/test_arm64.deb" \
    "http://localhost:8085/debian/repository/artisan/dist/all/section/devel"
```
### Deleting a Package

curl -X DELETE -u user:pwd http://**HOSTNAME:PORT**/debian/repository/**REPO-NAME**/distro/**DISTRO-NAME**/section/**SECTION-NAME**/version/**VERSION-REGEX**

Example:

deletes from "artisan" repository, distribution "all", section "main" any packages with a version matching 0.4.*

```bash
curl -X DELETE http://localhost:8085/debian/repository/artisan/distro/all/section/main/version/0.4.*
```

### Retrieving the PGP public key for a repository

curl http://**HOSTNAME:PORT**/debian/repository/**REPO-NAME**/key

Example:

```bash
curl  http://localhost:8085/debian/repository/artisan/key
```

### Repository Browsing

Files for all debian repositories can be browsed under **http(s)://HOSTNAME:PORT/debian/repositories**

Example:

```bash
$ python -mwebbrowser http://localhost:8085/debian/repositories
```