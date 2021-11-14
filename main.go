package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fvbock/endless"
	"github.com/joho/godotenv"
	echo "github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
)

// successResponse ...
type successResponse struct {
	Status bool `json:"status"`
}

const assetsPath = "assets/"

// Params for the server
type Params struct {
    Host    string `json:"host"`
    Port    string `json:"port"`
    PidDir  string `json:"pid"`
}

func main() {
	e := echo.New()

	loadEnv()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(touch())
	if os.Getenv("USERNAME") != "" && os.Getenv("PASSWORD") != "" {
		basicAuth(e)
	}

	e.GET("/health", health)
	e.Static("/assets", assetsPath)
	e.POST("/upload", upload)

    params := cliParams()

	serveGracefully(e, params.Host, params.Port, params.PidDir)

}

//named paramters
func cliParams() Params {
    params := Params{
        Host: "localhost",
        Port: "3000",
        PidDir:  "./",
    }
    //go run main.go -host=localhost -port=3000 -pidDir=./
    for _, arg := range os.Args {
        if strings.HasPrefix(arg, "-host=") || strings.HasPrefix(arg, "--host=") {
            params.Host = strings.ReplaceAll(arg, "-host=", "")
            params.Host = strings.ReplaceAll(arg, "--host=", "")
        }
        if strings.HasPrefix(arg, "-port=") || strings.HasPrefix(arg, "--port=") {
            params.Port = strings.ReplaceAll(arg, "--port=", "")
        }
        if strings.HasPrefix(arg, "-pidDir=") || strings.HasPrefix(arg, "--pidDir=") {
            params.PidDir = strings.ReplaceAll(arg, "--pidDir=", "")
        }
    }
    return params
}

func touch() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			uri := req.RequestURI
			if strings.HasPrefix(uri, "/assets/") {
				filename := strings.ReplaceAll(uri, "/assets/", "")
				currenttime := time.Now().Local()

				err := os.Chtimes(assetsPath+filename, currenttime, currenttime)
				if err != nil {
					log.Println(err)
				}
			}
			err := next(c)
			if err != nil {
				return err
			}

			return nil
		}
	}
}

func basicAuth(e *echo.Echo) {
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == os.Getenv("USERNAME") && password == os.Getenv("PASSWORD") {
			return true, nil
		}
		return false, nil
	}))
}

func serveGracefully(e *echo.Echo, host string, port, pidDir string) {
    e.Server.Addr = host + ":" + port
    server := endless.NewServer(e.Server.Addr, e)
    server.BeforeBegin = func(add string) {
        log.Print("info: actual pid is", syscall.Getpid())
        pidFile := filepath.Join(pidDir, port+".pid")
        err := os.Remove(pidFile)
        if err != nil {
            log.Print("error: pid file error: ", err)
        } else {
            log.Print("success: pid file success: ", pidFile)
        }
        err = ioutil.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0600)
        if err != nil {
            log.Print("error: write pid file error: ", err)
        } else {
            log.Print("success: write pid file success: ", pidFile)
        }
    }
    if err := server.ListenAndServe(); err != nil {
        log.Print("critical: graceful error: ", err)
    }
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func health(c echo.Context) error {
	return c.JSON(http.StatusOK, &successResponse{Status: true})
}

func upload(c echo.Context) error {
	//-----------
	// Read file
	//-----------

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		log.Print(err.Error())
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(assetsPath + file.Filename)
	if err != nil {
		return err
	}

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &successResponse{Status: true})
}
