# About JianDan MM

Get jandan.net ooxx picture.
获取煎蛋网妹子图的图片，并保存在本地。


## 更新
2013/9/10 V 1.0，直接使用net/http标准库进行抓取
2015/6/17 V 2.0，由于煎蛋妹妹改版后，单纯抓页面会出现403错误。所以使用watir，模拟浏览器进行抓取
2015/6/17 V 3.0，发现watir有时会Timeout error，重新抓包后发现，添加User-Agent的head标识即可正常返回200 OK。参考http://augustl.com/blog/2010/ruby_net_http_cheat_sheet

## 备注
脚本执行目录需手动新建test文件夹