# A backend system for review service using Golang.
The project is the modified and enriched version of Qimi's online courses. If there is any inferingement, please contact me for deletion.
Course Address： [course address](https://study.163.com/course/courseMain.htm?courseId=1212937804)
## QuickView For Codes
### review-service: providing Remote Process Calls for users, stores and audits.

supported methods:
- Detailed logics of all methods in **service for users/shops/audits**, http apis and remote process calls are provided
- select reviews from elasticsearch by storeID.
- select reviews from elasticsearch with not null comments.
 
### service for users: not inplemented serperately, http apis and grpc methods are written in **review-service**.

supported methods:
- create review
- update review
- delete review
- see review details
- see all reviews created by someone

### service for shops: review-b.

supported methods:(remote calls in **review-service**)
- create for reply
- update for reply
- appeal for reply

### service for audits: review-o.
supported methods:(remote calls in **review-service**)
- audit for reviews
- audit for appeals

### read messages from kafka into elasticsearch:review-job.

## Necessary Environments:
#### Go Environment
#### Protoc (add into your system PATH): [link](https://github.com/google/protobuf/releases)
#### Protoc-gen (add into your GOPATH):
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

```
#### Kratos:
  ```
  go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
  ```
#### MySQL(Local):v8.1.0
Suggest to setup data tables in your MySQL first (review_info,review_reply_info and review_appeal_info), see [.sql file](https://github.com/MysteriousX0214/Review-Service/blob/master/review-service/review.sql) for details. Denote the database named **review**. 
#### Redis(Local):v.3.2.100
(**unimplemented**) Add cache to redis when querying for reviews.
#### Docker(Local): 
Find a suitable version in [docker official website](https://www.docker.com/), you need to create a account first.
#### Consul(In Docker):
```
git clone https://github.com/hashicorp/learn-consul-docker.git
cd datacenter-deploy-service-discovery
docker-compose up -d
```
Or view [link](https://www.liwenzhou.com/posts/Go/consul/) for details
#### Canal(In Docker): 
```
  docker pull canal/canal-server:latest
  docker run -d --name canal-server -p 11111:11111 canal/canal-server
```
#### Kafka(In Docker):
```
  go get github.com/segmentio/kafka-go
```
See tutorial in [link](https://www.liwenzhou.com/posts/Go/kafka-go/) to setup kafka,zookeeper and kafka-ui in Docker
#### ElasticSearch(In Docker):
```
  go get github.com/elastic/go-elasticsearch/v8@latest
```
See tutorial in [link](https://www.liwenzhou.com/posts/Go/elasticsearch/) to setup elasticsearch and Kibana in Docker
#### Postman(Local,Optional):
To test if http apis and grpcs works well.
Find a suitable version in https://www.postman.com/, you may need to create an account for convinent use (like storing a certain http/grpc request route for multiplexing).

## How to run
### Configs
#### Extra account named 'canal'(or any name you like) in MySQL to store binlogs from database "review"(storing three data tables mentioned above)
```
CREATE USER 'canal'@'localhost' IDENTIFIED BY 'canal';
GRANT ALL PRIVILEGES ON my_database.* TO 'canal'@'localhost';
FLUSH PRIVILEGES;
```
#### Canal in Docker
Enter Canal Container, execute:
```
vi canal-server/conf/example/instance.properties
```
modify the following settings and save, then canal is able to catch changes in database "review":
```
canal.instance.master.address=host.docker.internal:3306 (when MySQL is deployed locally, if it is in Docker, use 127.0.0.1 or localhost instead.)
canal.instance.tsdb.dbUsername=canal (Extra account created just now)
canal.instance.tsdb.dbPassword=canal
canal.instance.dbUsername=root (your main account storing database "review")
canal.instance.dbPassword=root
```
you can also set the topic name of kafka, for examle, put the message from canal(listening "review") to kafka, and set the topic name such as **example**:
```
canal.mq.dynamicTopic=example:review\\..*
```
then execute:
```
vi canal-server/conf/canal.properties
```
modify the following settings and save, then canal is able to push changes to kafka:
```
serverMode:kafka
kafka.bootstrap.servers = host.docker.internal:29092 (your kafka is deployed in docker)
```
#### Configs in project
Check configs/config.yaml under each service folder (like review-service/configs/config.yaml):
- You can change the ports of http and grpc service. 
- Check the address of your MySQL.
- Keep consul and elasticsearch's address same with containers in Docker.

### Run
Execute the following command in the root directory of each folder:
```
kratos run
```
- **review-job** reads from kafka and writes into elasticsearch.
- **review-service** register all remote process calls to consul
- **review-b/review-o** discover remote process calls from consul
- changes in Mysql tables (caused by http or rpc calls) will be captured by canal and pushed to kafka (if your setting is right in above steps).

### Check
- check if the service is available with postman or [swagger editor](https://editor.swagger.io/)
- see routes and parameters in openapi.yaml under each root directory, like [link](https://github.com/MysteriousX0214/Review-Service/blob/master/review-service/openapi.yaml).
