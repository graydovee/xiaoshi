# 构建

```shell
make docker-build
``` 

# 运行
```shell
docker run -itd -v ${PWD}:config.yaml:/etc/xiaoshi/config.yaml graydovee/xiaoshi:latest -c /etc/xiaoshi/config.yaml
```