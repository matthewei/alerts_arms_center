package aac

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/matthewei/alerts_arms_center/utils/my_lib/aactl/templates"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/url"
	"os"
)

type aacCommend struct {
}

func NewAACAlertCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "send-alerts",
		Short:   "Send alerts to Weebhook.",
		Example: usage(),
		Args:    templates.ExactArgs("send-alerts", 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := aacCommend{}
			return c.Run(args)
		},
	}
	cmd.SetUsageTemplate(templates.UsageTemplate())
	return cmd
}
func usage() string {
	usage := `  # Send alerts to Weebhook URL!:
  {{ProgramName}} send-alerts webhookurl
`

	return usage
}

type Result struct {
	Data string `json:"data"`
}
type Repo struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (a *aacCommend) Run(args []string) error {
	webhook_url := args[0]

	//判断是否是合法的url
	parse, err := url.Parse(webhook_url)
	if err != nil {
		return err
	}
	if !parse.IsAbs() {
		return errors.New("url is not valid")
	}
	// 打开json文件
	jsonFile, err := os.Open("./alerts_firing.json")
	// 最好要处理以下错误
	if err != nil {
		fmt.Println(err)
	}
	// 要记得关闭
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var SendBody map[string]interface{}
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &SendBody)
	client := req.C().DevMode()
	resp, err := client.R().SetBody(&SendBody).SetSuccessResult(&result).Post(webhook_url)
	if err != nil {
		return err
	}
	if !resp.IsSuccessState() {
		fmt.Println("bad response status:", resp.Status)
		return errors.New("bad response")
	}
	fmt.Println("+++++++++++++++++++++++request successfully!++++++++++++++++++++")
	return nil
}
