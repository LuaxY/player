package main

import (
	"bytes"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
)

const assetPath = "./assets/"

func player(c echo.Context) error {
	path := c.Request().URL.Path

	if _, err := os.Stat(assetPath + path); !os.IsNotExist(err) {
		return c.File(assetPath + path)
	}

	path = strings.Replace(path, "/player/", "", 1)

	res, err := http.Get("https://ssl.p.jwpcdn.com/player/" + path)

	if err != nil {
		log.Error(errors.Wrap(err, "http get file"))
		return c.String(http.StatusInternalServerError, "something went wrong")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Error(errors.Wrap(err, "read body response"))
		return c.String(http.StatusInternalServerError, "something went wrong")
	}

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		var re *regexp.Regexp

		// Crack 8.2
		re = regexp.MustCompile(`(?m)ar s=Object\(r\.a\)\(e,Object\(o\.a\)\(c\)\),l=s\.split\("\/"\);`)
		body = re.ReplaceAll(body, []byte(`var l = e.split("/"");`))

		// Crack 8.3
		re = regexp.MustCompile(`(?m)var s=Object\(r\.a\)\(t,Object\(o\.a\)\(c\)\)\.split\("\/"\);`)
		body = re.ReplaceAll(body, []byte(`var s = t.split("/");`))

		// Crack 8.7+
		re = regexp.MustCompile(`(?m)var s=Object\(r\.a\)\(t\|\|"",Object\(o\.a\)\(".*"\)\)\.split\("\/"\);`)
		body = re.ReplaceAll(body, []byte(`var s = t.split("/");`))

		re = regexp.MustCompile(`(?m)u\+"//(.*)/ping\.gif`)
		body = re.ReplaceAll(body, []byte(`u+"//"+domain+"/ping.gif`))

		body = bytes.Replace(body, []byte(`ssl.p.jwpcdn.com`), []byte(domain), -1)
		body = bytes.Replace(body, []byte(`entitlements.jwplayer.com`), []byte(domain), -1)
		body = bytes.Replace(body, []byte(`prd.jwpltx.com`), []byte(domain), -1)
		body = bytes.Replace(body, []byte(`jwpltx.com`), []byte(domain), -1)
		body = bytes.Replace(body, []byte(`jwpsrv.com`), []byte(domain), -1)

		body = bytes.Replace(body, []byte(`this.children.audioTracks&&this.children.audioTracks.items[e].activate()`), []byte(`At(this.children.audioTracks,e)`), -1)

		dir := filepath.Dir(path)

		if err = os.MkdirAll(assetPath+"player/"+dir, os.ModePerm); err != nil {
			log.Warn(errors.Wrap(err, "mkdir cache path"))
		}

		if err = ioutil.WriteFile(assetPath+"player/"+path, body, os.ModePerm); err != nil {
			log.Warn(errors.Wrap(err, "write cache file"))
		}

		return c.Blob(http.StatusOK, mime.TypeByExtension(filepath.Ext(path)), body)
	}

	return c.Blob(res.StatusCode, res.Header.Get("Content-Type"), body)
}
