package my1Project

import (
	"github.com/gin-gonic/gin"
	//"fmt"
	"crypto/sha1"
	"fmt"
	"github.com/akkuman/parseConfig"
	"io"
	"sort"
	"strings"
)

func getArgTest(c *gin.Context) {
	name := c.DefaultQuery("name", "default")
	c.JSON(200, gin.H{
		"name":   name,
		"method": "GET",
	})
}

func makeSignature(timestamp, nonce, echostr string) string {
	strSlice := []string{timestamp, nonce, token}
	sort.Sort(sort.StringSlice(strSlice))
	srcStr := strings.Join(strSlice, "")
	fmt.Println(srcStr)
	sha1obj := sha1.New()
	io.WriteString(sha1obj, srcStr)
	sha1rst := sha1obj.Sum(nil)
	return fmt.Sprintf("%x", sha1rst)
}

func getSignature(c *gin.Context) {
	signature := c.DefaultQuery("signature", "")
	timestamp := c.DefaultQuery("timestamp", "")
	nonce := c.DefaultQuery("nonce", "")
	echostr := c.DefaultQuery("echostr", "")
	rightSignature := makeSignature(timestamp, nonce, echostr)

	data := make(map[string]string)
	data["in_signature"] = signature
	data["in_timestamp"] = timestamp
	data["in_nonce"] = nonce
	data["echostr"] = echostr

	data["out_signature"] = rightSignature

	if signature == rightSignature {
		/*
			c.JSON(200, gin.H{
				"code": "0",
				"msg": "success",
				"data": data,
			})
		*/
		c.String(200, echostr)
	} else {
		c.String(404, "error")
		/*
			c.JSON(200, gin.H{
				"code": "-1",
				"msg": "error",
				"data": data,
			})
		*/
	}

}

func postArgTest(c *gin.Context) {
	name := c.DefaultPostForm("name", "default")
	c.JSON(200, gin.H{
		"name":   name,
		"method": "POST",
	})
}

func getIP(c *gin.Context) {
	ip := c.ClientIP()
	c.JSON(200, gin.H{
		"ip": ip,
	})
}

func getHeader(c *gin.Context) {
	c.JSON(200, gin.H{
		"header": c.Request.Header,
	})
}

func abort(c *gin.Context) {
	c.AbortWithStatusJSON(403, "test")
}

var token string

func Run(configFile string) {
	var config = parseConfig.New(configFile)
	// 此为interface{}格式数据
	token = (config.Get("token")).(string)

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
			"code":   0,
		})
	})

	r.GET("/getArg", getArgTest)
	r.GET("/ip", getIP)
	r.GET("/header", getHeader)
	r.GET("/abort", abort)
	r.GET("/signature", getSignature)

	r.Run(":8003")
}
