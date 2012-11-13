Gorbled v3.0
============

A simple blog system written by Go, Running in Google App Engine
--------------------------------------------------------------------------------------------------------

### Demo
* [Click here](http://blog.specode.org)

### Feature
* Post article
* Custom widget, can sort.
* Markdown support, Use [blackfriday](https://github.com/russross/blackfriday)
* File upload support
* Custom ID (URL optimization)
* DISQUS support

### Install && Config
* Edit ./config.json Change blog config : disqus, google analytics, lang
* Edit ./app.yaml for your gae app config
* Using GAE SDK update app

### License
* [blackfriday](https://github.com/russross/blackfriday) BSD License
* [Gorbled](https://github.com/specode/Gorbled) MIT License

### todo list: v4
* 主题
  + 我要一个和blogger.com一样的侧边栏 = =
* 利用GAE的免费配额，优化数据库操作，Memcache?
* Page页
* 优化针对google的索引
* 数据导入导出(wordpress格式？)
* 搜索页
* 标签目录？
* 