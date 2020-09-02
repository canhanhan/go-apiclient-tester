你好！
很冒昧用这样的方式来和你沟通，如有打扰请忽略我的提交哈。我是光年实验室（gnlab.com）的HR，在招Golang开发工程师，我们是一个技术型团队，技术氛围非常好。全职和兼职都可以，不过最好是全职，工作地点杭州。
我们公司是做流量增长的，Golang负责开发SAAS平台的应用，我们做的很多应用是全新的，工作非常有挑战也很有意思，是国内很多大厂的顾问。
如果有兴趣的话加我微信：13515810775  ，也可以访问 https://gnlab.com/，联系客服转发给HR。
﻿# go-apiclient-tester
[![GoDoc-Postman](https://godoc.org/github.com/finarfin/go-apiclient-tester/postman?status.svg)](https://godoc.org/github.com/finarfin/go-apiclient-tester/postman)
[![GoDoc-Tester](https://godoc.org/github.com/finarfin/go-apiclient-tester/tester?status.svg)](https://godoc.org/github.com/finarfin/go-apiclient-tester/tester)

A simple, easy to use library to test HTTP based API clients.

Main use-case for this library is to be able to use [Postman Collections](https://learning.postman.com/docs/postman/collections/intro-to-collections/) for unit testing. Create requests/examples in [Postman](https://www.postman.com/); then export as "Collection v2.1" format. It can verify that requests from your API client match to the ones in the collection and respond with the recorded response.

## Example

```go
    package awesomeclient

    import (
        "testing"
        "github.com/finarfin/go-apiclient-tester/postman"
        "github.com/finarfin/go-apiclient-tester/tester"
        "github.com/stretchr/testify/assert"
    )

    func TestUserCreateSuccess(t *testing.T) {
        tester, err := postman.NewTester("testdata/collection.json")        
        if err != nil {
            t.Fatal(err)
        }
        defer tester.Close()
        tester.Setup("user", "create_success")

        _, err = c.CreateUser("user1")
        if err != nil {
            t.Fatal(err)
        }
    }

    func TestUserCreateFailure(t *testing.T) {
        tester, err := postman.NewTester("testdata/collection.json")
        if err != nil {
            t.Fatal(err)
        }
        defer tester.Close()
        tester.Setup("user", "create_failure")

        _, err = c.CreateUser("user1")
        assert.Error(t, err)
    }
```

## License

[MIT](LICENSE)
