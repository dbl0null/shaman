package commands

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/cobra"

	"github.com/nanopack/shaman/config"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update entry in shaman database",
	Long:  ``,

	Run: update,
}

type updateBody struct {
	value string
}

func update(ccmd *cobra.Command, args []string) {
	if len(args) != 3 {
		fmt.Fprintln(os.Stderr, "Missing arguments: Needs record type, domain, and value")
		os.Exit(1)
	}
	var client *http.Client
	if config.Insecure {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	} else {
		client = http.DefaultClient
	}
	rtype := args[0]
	domain := args[1]
	value := args[2]
	fmt.Println("rtype:", rtype, "domain:", domain, "value:", value)
	data := url.Values{}
	data.Set("value", value)

	uri := fmt.Sprintf("https://%s:%s/records/%s/%s?%s", config.ApiHost, config.ApiPort, rtype, domain, data.Encode())
	fmt.Println(uri)
	req, err := http.NewRequest("PUT", uri, bytes.NewBufferString(data.Encode()))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	req.Header.Add("X-NANOBOX-TOKEN", config.ApiToken)
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	fmt.Println(string(b))
}