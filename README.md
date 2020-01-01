# Demo Service (ElasticSearch, gRPC)
## Prerequisite
1. ElasticSearch Services, with data and index table (please read the proto file for detail in ./protos)
2. Your client side PHP project (Optional)

## Generate Protocol Buffers
### (Example for EsApi service esapi.proto)

### For Golang

Prerequisite:
- Install the protocol compiler plugins https://grpc.io/docs/languages/go/quickstart/

1. go to "/protos" and create a directory called "esapi"

```bash
cd ./protos && mkdir esapi
```

2. execute the following command to generate protos files in "esapi" directory

```bash
protoc --go_out=./esapi --go_opt=paths=source_relative --go-grpc_out=./esapi --go-grpc_opt=paths=source_relative esapi.proto
```

3. move the "esapi" directory to "../search-service/protos"

```bash
mv esapi ../search-service/protos
```

<br>
<br>

### For PHP (Generate Client Side)

1. uncomment the "php-grpc-env" setting in docker-compose.yml
2. run the service in docker

```bash
docker-compose up --build --force-recreate --no-deps -d
```

3. execute the following command to run the shell script in "php-grpc-env" terminal
```bash
/var/www/protos/gen_esapi_php.sh
```
4. Copy the generated files to your PHP client side project. (Esapi folder, GPBMetadata folder)

# Start the Service
## env setup
1. Create a replica ".env.example" file and name it to ".env"
2. modify the var

## For Docker Development Environment

#set up debugger \
set up the Delve Debugger in your IDE and listen to port :40000

```bash
# for build and recreate
docker-compose up --build --force-recreate --no-deps -d
```
