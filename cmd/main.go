package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"webhook-proxy/internal/helper"
	"webhook-proxy/internal/k8sclient"
	"webhook-proxy/internal/plugins"

	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
)

type Webhook struct {
	Name      string
	Namespace string
	Phase     string
	Metadata  map[string]interface{}
}

type WebHandler struct {
	context   context.Context
	clientSet *kubernetes.Clientset
}

func (wh WebHandler) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	detectedPlugin := vars["destination"]

	if r.Body != nil && r.ContentLength != 0 {

		var webhook Webhook
		err := json.NewDecoder(r.Body).Decode(&webhook)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var status int
		detectedDeployment := webhook.Name
		// Setup Deployment Name
		if v, found := webhook.Metadata["name"]; found {
			detectedDeployment = v.(string)
		}
		// Lookup Deployment
		sha, err := k8sclient.GetTag(webhook.Namespace, detectedDeployment, wh.clientSet, wh.context)
		if err != nil {
			fmt.Printf("Could not get tag for %s\n", detectedDeployment)
			status = http.StatusBadRequest
		} else {

			// Pick the correct Plugin
			switch detectedPlugin {
			case "sleuth":
				sleuth_environment := helper.GetEnv("SLEUTH_ENVIRONMENT", "")

				sleuth_api_key := helper.GetEnv("SLEUTH_API_KEY", "")
				if len(sleuth_environment) != 0 && len(sleuth_api_key) != 0 {

					plugins.SleuthWebhook(detectedDeployment, sha, sleuth_api_key, sleuth_environment)
				} else {
					fmt.Println("Missing SLEUTH_API_KEY & SLEUTH_ENVIRONMENT ")
				}

			default:
				fmt.Printf("Webhook Plugin %s not supported\n", detectedPlugin)
			}
			status = http.StatusOK
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
	} else {
		fmt.Println("Not able to process webhook")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	wh := WebHandler{}
	wh.clientSet, wh.context = k8sclient.Inconfig()
	r := mux.NewRouter()
	r.HandleFunc("/{destination}", wh.WebhookHandler).Methods("POST")
	log.Println("Listing for requests at http://localhost:8080/{destination}")
	log.Fatal(http.ListenAndServe(":8080", r))
}
