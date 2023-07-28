package gitlab

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	webhook "github.com/go-playground/webhooks/v6/gitlab"
	"github.com/gorilla/mux"
	"github.com/xanzy/go-gitlab"
	"quillbot.com/gitlab-webhook-parser/src/slack"
)

func initGilab() (cl *gitlab.Client, err error) {
	gitlabToken := os.Getenv("GITLAB_TOKEN")
	if gitlabToken == "" {
		log.Fatalf("Please provide valid GITLAB_TOKEN")
	}
	client, err := gitlab.NewClient(gitlabToken, gitlab.WithBaseURL("https://gitlab.com/api/v4"))
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return nil, err
	}
	return client, nil
}

func updateProtectedBranchAccess(projectID int, branchName string, perm gitlab.AccessLevelValue) error {
	client, err := initGilab()
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return err
	}
	pb, _, err := client.ProtectedBranches.GetProtectedBranch(projectID, branchName)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return err
	}
	pushAccess := []*gitlab.BranchPermissionOptions{&gitlab.BranchPermissionOptions{
		ID:          &pb.PushAccessLevels[0].ID,
		AccessLevel: gitlab.AccessLevel(perm),
	}}
	mergeAccess := []*gitlab.BranchPermissionOptions{&gitlab.BranchPermissionOptions{
		ID:          &pb.MergeAccessLevels[0].ID,
		AccessLevel: gitlab.AccessLevel(perm),
	}}
	updatedOptions := &gitlab.UpdateProtectedBranchOptions{
		CodeOwnerApprovalRequired: gitlab.Bool(true),
		AllowedToPush:             &pushAccess,
		AllowedToMerge:            &mergeAccess,
	}
	_, _, err = client.ProtectedBranches.UpdateProtectedBranch(projectID, branchName, updatedOptions)
	if err != nil {
		fmt.Printf("Error while updating protected branch - %v", err)
		return err
	}
	return nil
}

func sendError(w http.ResponseWriter, statusCode int, err string) {
	errorMessage := map[string]string{"error": err}
	sendResponse(w, http.StatusInternalServerError, errorMessage)
}

func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	resp, _ := json.Marshal(payload)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(resp)
}

func updateMergePermissions(w http.ResponseWriter, r *http.Request, perm gitlab.AccessLevelValue) {
	returnPl := map[string]string{}
	ipId := mux.Vars(r)["pId"]
	pId, err := strconv.Atoi(ipId)
	if err != nil {
		sendError(w, http.StatusNotFound, "Invalid Payload")
		return
	}
	branchName := mux.Vars(r)["branchName"]
	if branchName == "" {
		sendError(w, http.StatusNotFound, "Invalid Payload")
		return
	}
	err = updateProtectedBranchAccess(pId, branchName, perm)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	returnPl["status"] = "success"
	sendResponse(w, http.StatusOK, returnPl)
}

func EnableMerge(w http.ResponseWriter, r *http.Request) {
	updateMergePermissions(w, r, gitlab.AccessLevelValue(40))
}

func DisableMerge(w http.ResponseWriter, r *http.Request) {
	updateMergePermissions(w, r, gitlab.AccessLevelValue(0))
}

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	secretToken := os.Getenv("SECRET_TOKEN")
	if secretToken == "" {
		sendError(w, http.StatusBadRequest, "Please provide valid SECRET_TOKEN")
		return
	}
	hook, err := webhook.New(webhook.Options.Secret(secretToken))
	if err != nil {
		err := fmt.Sprintf("Error initialization in gitlab webhook %v", err.Error())
		sendError(w, http.StatusBadRequest, err)
		return
	}
	returnPl := map[string]string{}
	payload, err := hook.Parse(r, webhook.MergeRequestEvents)
	if err != nil {
		err := fmt.Sprintf("Error initialization in gitlab webhook %v", err.Error())
		sendError(w, http.StatusBadRequest, err)
		return
	}
	switch payload.(type) {
	case webhook.MergeRequestEventPayload:
		mrpl := payload.(webhook.MergeRequestEventPayload)
		switch mrpl.ObjectAttributes.State {
		case "merged":
			if mrpl.ObjectAttributes.TargetBranch == "main" || mrpl.ObjectAttributes.TargetBranch == "master" {
				sleepSecsString := os.Getenv("SLEEP_DURATION")
				if sleepSecsString == "" {
					log.Fatalf("Please provide valid SLEEP_DURATION")
				}
				sleepSecs, err := strconv.Atoi(sleepSecsString)
				if err != nil {
					sendError(w, http.StatusBadRequest, err.Error())
					return
				}
				time.Sleep(time.Duration(sleepSecs) * time.Second)
				err = updateProtectedBranchAccess(int(mrpl.Project.ID), mrpl.ObjectAttributes.TargetBranch, 0)
				if err != nil {
					sendError(w, http.StatusBadRequest, err.Error())
					return
				}
				returnPl["status"] = "success"
				slack.SendSlackNotification(mrpl)
				sendResponse(w, http.StatusOK, returnPl)
			}
		}
	}
}
