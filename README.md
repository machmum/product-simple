# Product-simple

- set env based on `env` file

- docker build
```sh
_IMAGE_NAME="product-simple"
_BUILD_FILE="rest"
docker build -t $_IMAGE_NAME \
	--build-arg _BUILD_FILE=$_BUILD_FILE \
	--no-cache -f Dockerfile .
```

- docker run
```sh
docker run -it \
--env-file ./env \
-p 8080:8080 \
--name product-simple $_IMAGE_NAME /bin/sh
```