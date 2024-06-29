package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"

	"github.com/gabriel-vasile/mimetype"
)

type Config struct {
	Stills  string `json:"stills"`
	HStills string `json:"hstills"`
	Mov     string `json:"mov"`
	HMov    string `json:"hmov"`
	Ai      string `json:"ai"`
}

var config Config

func main() {
	r := gin.Default()

	ip := flag.String("ip", "127.0.0.1", "IP to listen on")
	port := flag.String("port", "4000", "Port to listen on")
	configPath := flag.String("config", "./config.json", "Config file path")
	flag.Parse()

	address := fmt.Sprintf("%s:%s", *ip, *port)
	config = readConfig(*configPath)

	r.Use(cors.Default())
	r.Use(headers)

	r.GET("/", home)
	r.POST("/r/", real)
	r.POST("/v/", virtual)
	r.POST("/ai/", ai)

	// this should be a test
	//log.Println(getFinalDir(config, false, false, false), getFinalDir(config, true, false, false), getFinalDir(config, true, true, false), getFinalDir(config, true, true, true), getFinalDir(config, false, true, true), getFinalDir(config, false, false, true), getFinalDir(config, true, false, true))

	r.Run(address)
}

func readConfig(configPath string) Config {
	var config Config

	data, _ := os.ReadFile(configPath)
	_ = json.Unmarshal(data, &config)

	return config
}

func headers(c *gin.Context) {
	c.Header("Content-Security-Policy", "default-src *;")
	//c.Header("Access-Control-Allow-Origin", "*")
	//c.Header("Access-Control-Allow-Methods", "POST")

	c.Next()
}

func home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"error": "home",
	})
}

type Data struct {
	Url string `form:"url"`
}

type Err struct {
	Message string
}

func real(c *gin.Context) {
	var d Data
	c.Bind(&d)

	if d.Url == "" {
		c.JSON(http.StatusOK, gin.H{
			"error": "Url is empty",
		})
		return
	}

	res, err := reqFile(d.Url, true, false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": res,
	})
}

func virtual(c *gin.Context) {
	var d Data
	c.Bind(&d)

	if d.Url == "" {
		c.JSON(http.StatusOK, gin.H{
			"error": "Url is empty",
		})
		return
	}

	res, err := reqFile(d.Url, false, false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": res,
	})
}

func ai(c *gin.Context) {
	var d Data
	c.Bind(&d)

	if d.Url == "" {
		c.JSON(http.StatusOK, gin.H{
			"error": "Url is empty",
		})
		return
	}

	res, err := reqFile(d.Url, false, true)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": res,
	})
}

func reqFile(uri string, isReal bool, isAi bool) (string, error) {
	log.Printf("Req '%s'", uri)

	puri, _ := url.Parse(uri)
	res, err := http.Get(uri)
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	toVid := isVid(body)

	file := path.Base(puri.Path)
	finalPath := path.Join(getFinalDir(config, isReal, toVid, isAi), file)
	f, _ := os.Create(finalPath)

	io.Copy(f, bytes.NewReader(body))

	message := fmt.Sprintf("File '%s' | is vid => %t", file, toVid)

	return message, err
}

func isVid(file []byte) bool {
	tp := mimetype.Detect(file)

	return strings.HasPrefix(tp.String(), "video")
}

func getFinalDir(config Config, isReal bool, isVid bool, isAi bool) string {
	result := ""

	if isAi {
		result = config.Ai
	} else if isReal && isVid {
		result = config.Mov
	} else if isReal && !isVid {
		result = config.Stills
	} else if !isReal && isVid {
		result = config.HMov
	} else if !isReal && !isVid {
		result = config.HStills
	}

	return result
}
