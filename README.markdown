Gorbled
=======

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
* Edit ./config.json Change blog config
* Edit ./app.yaml for your gae app config
* Edit ./gorbled/static/html/article/view.html for change disqus config
* Edit ./gorbled/static/html/layouts/main.html for change google analytics config
* Using GAE SDK update app

### License
* [blackfriday](https://github.com/russross/blackfriday) BSD License
* [Gorbled](https://github.com/specode/Gorbled) MIT License
