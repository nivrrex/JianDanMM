# About JianDan MM

Get jandan.net ooxx picture.

获取煎蛋网妹子图的图片，并保存在本地。


## 更新
2013/9/10  V1.0  直接使用net/http标准库进行抓取。

2015/6/17  V2.0  由于煎蛋妹妹改版后，单纯抓页面会出现403错误。所以使用watir，模拟浏览器进行抓取。

2015/6/17  V2.1  发现watir有时会Timeout error，重新抓包后发现，添加User-Agent的head标识即可正常返回200 OK，重新使用net/http实现。

参考https://github.com/augustl/net-http-cheat-sheet


## 备注
脚本执行目录需手动新建test文件夹
