# go-dict-server

## 功能
1. 通过```sqlite3```读取数据库文件；
2. 通过对外提供接口，返回```json```格式的数据；

## 使用
1. 进入 ```go-dict-server-prepare``` 文件夹
2. ```go build``` 然后会生成 ```go-dict-server-prepare``` 可执行文件
3. 如果需要指定特殊端口，使用如下命令<br>```nohup ./go-dict-server-linux -port 1234 -dbpath abcdef &```
4. 如果使用默认端口4000，则使用如下命令<br>```nohup ./go-dict-server-linux &```<br>```nohup xxx &```的意思是，后台执行，并且对ctrl+c免疫。如果不加 ```&``` 则使用ctrl+c后，xxx程序会关闭
   
