### Install dependencies dep

    go get -u github.com/golang/dep/cmd/dep
    dep ensure


### Alternative way

    go get github.com/gin-gonic/gin
    go get github.com/gin-contrib/cors
    go get github.com/gin-contrib/location
    go get github.com/jinzhu/gorm/dialects/mysql
    go get github.com/unidoc/unidoc/pdf/creator
    go get gopkg.in/gographics/imagick.v3/imagick

### Log file will be rewritten each time, when you rerun go lang application.
By default, file will be created in the root of current module and named stdout.log

