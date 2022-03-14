package plugins

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"webhook-proxy/internal/helper"
)

func SleuthWebhook(deploymentName string, sha string, sleuth_api_key string, sleuth_environment string) {

	//Encode the data
	postBody := strings.NewReader(fmt.Sprintf("api_key=%s&environment=%s&sha=%s", sleuth_api_key, sleuth_environment, sha))

	//Build Sleuth URL
	orgslug := helper.GetEnv("SLEUTH_ORG_SLUG", "DEFAULTORGSLUG")
	resp, err := http.Post(fmt.Sprintf("https://app.sleuth.io/api/1/deployments/%s/%s/register_deploy", orgslug, deploymentName), "application/json", postBody)

	//Handle Error
	if err != nil {
		fmt.Printf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	sb := string(body)
	fmt.Println(sb)

	fmt.Println(resp.Status)
	switch resp.StatusCode {
	case 200:
		fmt.Println("Successful Deployment Registration")
	case 400:
		fmt.Println("Input date problem, including if SHA doesn")
	case 401:
		fmt.Println("API key not valid or the deployment is not in the specific organization")
	default:
		fmt.Println("Unknown failure")
	}
}
