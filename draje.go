package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type TagResponse struct {
	Name string   `json: "name"`
	Tags []string `json: "tags"`
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func getTagResponse(body []byte) (*TagResponse, error) {
	s := new(TagResponse)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println(err)
	}
	return s, err
}

func getTagListFromRegisty(host, image, token string) []string {
	url := host + "/v2/" + image + "/tags/list"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Basic "+token)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	s, _ := getTagResponse([]byte(body))
	return s.Tags
}

func removeTagFromRegistry(host, image, token, tag string) {
	url_tag := host + "/v2/" + image + "/manifests/" + tag
	req_tag, _ := http.NewRequest("GET", url_tag, nil)
	req_tag.Header.Add("Authorization", "Basic "+token)
	req_tag.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	res_tag, _ := http.DefaultClient.Do(req_tag)

	defer res_tag.Body.Close()
	if res_tag.StatusCode != 200 {
		fmt.Println("image: " + image + " tag: " + tag + " could not get content digest, failed!")
		body, _ := ioutil.ReadAll(res_tag.Body)
		fmt.Println(string(body))
		return
	}
	digest := res_tag.Header.Get("Docker-Content-Digest")

	url_remove := host + "/v2/" + image + "/manifests/" + digest
	req_remove, _ := http.NewRequest("DELETE", url_remove, nil)
	req_remove.Header.Add("Authorization", "Basic "+token)
	req_remove.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	res_remove, _ := http.DefaultClient.Do(req_remove)

	if res_remove.StatusCode != 202 {
		fmt.Println("image: " + image + " tag: " + tag + " could not removed")
		body, _ := ioutil.ReadAll(res_tag.Body)
		fmt.Println(string(body))
		return
	}
	fmt.Println("image: " + image + " tag: " + tag + " removed!")

	defer res_remove.Body.Close()
}

func main() {
	app := cli.NewApp()
	app.Name = "draje - docker registry api from jenkins"
	app.Description = "uses docker registy v2 api for deleting images from docker registry"
	app.Version = "0.0.1"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Ahmet Oz",
			Email: "bilmuhahmet@gmail.com",
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Value: "",
			Usage: "host address of docker registry",
		},
		cli.StringFlag{
			Name:  "image",
			Value: "",
			Usage: "name of the image",
		},
		cli.StringFlag{
			Name:  "excluded_tags",
			Value: "latest",
			Usage: "name of the image tag which will be excluded",
		},
		cli.StringFlag{
			Name:  "username",
			Value: "",
			Usage: "registry user name",
		},
		cli.StringFlag{
			Name:  "password",
			Value: "",
			Usage: "registry user password",
		},
		cli.IntFlag{
			Name:  "exclude_last",
			Value: 1,
			Usage: "exclude last n images",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.String("host") == "" {
			return cli.NewExitError("host is required", 86)
		}
		if c.String("image") == "" {
			return cli.NewExitError("image is required", 86)
		}
		if c.String("username") == "" {
			return cli.NewExitError("user name is required", 86)
		}
		if c.String("password") == "" {
			return cli.NewExitError("password is required", 86)
		}

		token := basicAuth(c.String("username"), c.String("password"))

		host := c.String("host")
		image := c.String("image")
		tags := getTagListFromRegisty(host, image, token)
		excluded_tags := strings.Split(c.String("excluded_tags"), ",")
		included_tags := tags[:len(tags)-c.Int("exclude_last")]

		for _, tag := range included_tags {
			exclude := false
			for _, excluded_tag := range excluded_tags {
				if string(tag) == string(excluded_tag) {
					exclude = true
					break
				}
			}
			if !exclude {
				removeTagFromRegistry(host, image, token, tag)
			}
		}

		fmt.Println("Completed.")
		return nil
	}
	app.Run(os.Args)
}
