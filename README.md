# go-simple-spider

#### 介绍
每个网站都有自己的特点，所有Download和Resolve策略都不相同。针对这个问题，这个框架将提供一套并发流程，同时支持定制每一层的Download和Resolve. 用户只需要根据实际的网站规则，定制好Download和Resolve然后设置到框架中，其他的全部可以交给框架去自动处理。最终爬下来的结果，会缓存到一个buffer中，buffer大小可以自定义，buffer满了就会自动调用Flush()函数，自己可以实现Flush()函数从而自定义数据持久化路径。
###### 已经完成:
1. 整体框架结构开发基本结束。
2. 支持定制网站的每一层Download和Resolve。
3. 支持多个站点同时并发爬虫，互不干扰。
4. 日志支持，各种参数配置支持。
5. 下载数据异常，例如返回404,302，超时等等情况，重试。
6. 普通静态网页使用golang自带http下载，简单封装。
7. 需要渲染、动态网页将使用selenium, java selenium作为backend简单封装。
8. 包含htmlquery项目，根据需要，在自定义解析时使用。
9. src/parts/baidu和src/parts/sina是两个模板例子，各自自定义了三层下载规则和解析规则，解析的时候会解析出整个页面包含的href超链接。可以根据自己要爬的站点的需要，拷贝这其中任意一个模板，修改其中代码，自定义每一层的下载和解析。

###### TODO:
1. 增量爬虫
2. 网络ip代理
3. 防止回环，去重。
4. 已经爬的内容，缓存信息，加速爬取。
5. 下载的内容解析后不符合预期，需要重新下载，此种重试将利用提供的Unit.Backward()函数实现。
6. persistence.go 是最终解析出来的文本信息，爬虫出口。在文本数量缓存超过设置的numCachedText之后，会自动调用这里面的方法将缓存数据刷出来。目前只是将结果打印，没有存储到DB。
7. selenium里面SeleniumPool.java路径配置,有空再改。

#### 软件架构
待整理。。。
#### 安装教程
##### 依赖：
1. htmlquery vlastest: https://github.com/antchfx/htmlquery
2. selenium v3.14: https://www.selenium.dev/downloads/
依赖都已经放在了项目中，不需要再额外下载安装。

#### 使用说明
1.  启动selenium:
```shell
cd src/backend/java
javac -d . *.java -cp res/selenium-server-standalone-3.14.0.jar
java -cp .:res/selenium-server-standalone-3.14.0.jar server.Main
```
等待所有selenium池加载完毕。SeleniumPool.java里面driver路径要改一下。

2. main.go里面是同时爬baidu和sina的样例，按自己需要自定好之后，执行：
  ```
    export GOPATH="src目录的父目录"
    go run main.go
  ```
注意：第一次运行的时候程序会立即结束，会在GOPATH下生成configs和logs配置文件。修改好这些文件之后，再一次执行即可。改下common.conf里面maxDepth应该就可以运行demo了。

3. snet/request.go提供了使用Selenium请求http、使用golang网络包请求http两种封装，自定义下载的时候按需使用即可。

4. common.conf
  - maxDepth表示最大爬取深度。
  - numUnit表示组合成一个Group所需要的Unit个数(Note: 多个Unit组成一个Group，多个Group组成一层Layer，每个Group会启动一个线程，当前层所有节点结束才会进行下一层。每个Unit会链接多个Part, Part为最小可执行单位。)。
  - numCachedText表示缓存多少解析完成的文本数据之后，才会调用persistence函数刷出去。

5. log.conf
  - logRootDir表示日志存储根目录
  - logFileSize表示日志文件超过这个大小之后，才会新建一个日志文件。

6. global.log
  - 此文件是全局日志文件，不区分网站。

###### 其他
- author: Vandy Yu
- email: 13122017572@163.com
- This project is copied from my gitee project.[https://gitee.com/vandyyu/go-simple-spider/tree/framework-dev.0.9.1/]
