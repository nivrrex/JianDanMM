# About JianDan MM

Get jandan.net ooxx picture.

获取煎蛋网妹子图的图片，并保存在本地。


## 更新
2013/9/10   V1.0   直接使用net/http标准库进行抓取。

2015/6/17   V2.0   由于煎蛋妹妹改版后，单纯抓页面会出现403错误。所以使用watir，模拟浏览器进行抓取。

2015/6/17   V2.1   发现watir有时会Timeout error，重新抓包后发现，添加User-Agent的head标识即可正常返回200 OK，重新使用net/http实现。

参考https://github.com/augustl/net-http-cheat-sheet

2015/6/19   V2.2   增加记录图片的支持和反对记录数据log，同时完善对gif图片的提取

2015/6/19   V2.3   根据支持和反对的比例，对保存的文件名添加“very\high\good\normal\bad”标识符

2015/7/23   V3.0   增加Golang语言实现，同时完美支持jpg和gif文件，前期Ruby版本的gif文件下载的是缩略图

2015/7/23   V3.1   jandan.net对妹子图的抓取，每隔50个页面会重置，此时403错误。通过变更User-Agent可以解决该问题

2015/7/24   V3.2   根据Golang实现，修改Ruby实现，目前已经支持正常Gif抓取

2015/7/28   V3.3   Golang实现部分优化，同时可以通过config.json设置开始及结束的抓取页面,避免反复编译

2015/8/26   V3.4   对Golang实现的部分错误处理进行完善

2016/4/2    V4.0   增加Python语言实现，使用requests库，实现确实很简单。目前发现如果一个回复中，有多个图片的，可能存在遗漏，待更新。

2016/10/15  V5.0   Python语言实现，改用pyquery(lxml)库，同时实现一个回复中多个图片的下载。（需在脚本目录下创建test文件夹）

Windows系统可以在 http://www.lfd.uci.edu/~gohlke/pythonlibs/#lxml 下载lxml库

## 备注
~~脚本执行目录需手动新建test文件夹，~~目前文件下载到脚本当前目录
