package main

import (
	"fmt"
	"gopkg.in/AlecAivazis/survey.v1"
	"grb/model"
	"grb/project/beegocli"
	"grb/repo_creator"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

var qs = []*survey.Question{
	{
		Name:     "gitHostAddress",
		Prompt:   &survey.Input{Message: "Git Host address:"},
		Validate: survey.Required,
	},
	{
		Name: "gitServerVersion",
		Prompt: &survey.Select{
			Message: "Git server version:",
			Options: []string{"GitLab 6.3.0 LDAP"},
			Default: "GitLab 6.3.0 LDAP",
		},
	},
	{
		Name:     "repoName",
		Prompt:   &survey.Input{Message: "Main repository name:"},
		Validate: survey.Required,
	},
	{
		Name:     "repoNamespace",
		Prompt:   &survey.Input{Message: "All repository namespace:"},
		Validate: survey.Required,
	},
	{
		Name:     "username",
		Prompt:   &survey.Input{Message: "Git login username"},
		Validate: survey.Required,
	},
	{
		Name:     "password",
		Prompt:   &survey.Password{Message: "Git login password"},
		Validate: survey.Required,
	},
}

func main() {
	beegocli.CreateProject("testcli")

	//beegocli.Creator()
	//answers := model.Answer{}
	//
	////err := survey.Ask(qs, &answers)
	////if err != nil {
	////	fmt.Println(err.Error())
	////	return
	////}
	//
	//// 测试数据
	//answers = model.Answer{
	//	GitHostAddress:   "http://wpsgit.kingsoft.net/",
	//	GitServerVersion: "GitLab 6.3.0 LDAP",
	//	RepoName:         "grb",
	//	RepoNamespace:    "galaxy",
	//	Username:         "",
	//	Password:         "",
	//}
	//
	//// 选择creator
	//var repoCreator repo_creator.RepoCreator
	//if "GitLab 6.3.0 LDAP" == answers.GitServerVersion {
	//	repoCreator = &repo_creator.Gitlab630Ldap{}
	//} else {
	//	fmt.Println(errors.New(answers.GitServerVersion + " no implement yet"))
	//	return
	//}
	//
	//// 创建所有Repo
	//// 执行Git操作，将子Repo加入到父Repo中
	//mainRepoName := answers.RepoName
	//os.Mkdir(mainRepoName, os.ModeDir)
	//repoCreator.CreateRepo(answers)
	//run(exec.Command("git", "init"), mainRepoName)
	//
	//answers.RepoName = mainRepoName + "-admin"
	//createSubRepo(answers, mainRepoName, repoCreator)
	//answers.RepoName = mainRepoName + "-server"
	//createSubRepo(answers, mainRepoName, repoCreator)
	//answers.RepoName = mainRepoName + "-common"
	//createSubRepo(answers, mainRepoName, repoCreator)
	//answers.RepoName = mainRepoName + "-vendor"
	//createSubRepo(answers, mainRepoName, repoCreator)
	//
	//run(exec.Command("git", "add", "."), mainRepoName)
	//run(exec.Command("git", "commit", "-m", "\"init\""), mainRepoName)
	//parse, _ := url.Parse(answers.GitHostAddress)
	//gitRepoPath := "git@" + parse.Host + ":" + answers.RepoNamespace + "/" + mainRepoName + ".git"
	//run(exec.Command("git", "remote", "add", "origin", gitRepoPath), mainRepoName)
	//run(exec.Command("git", "push", "-u", "origin", "master"), mainRepoName)
}

func createSubRepo(subRepoAnswer model.Answer, mainRepoName string, repoCreator repo_creator.RepoCreator) {
	os.Mkdir(mainRepoName+"/"+subRepoAnswer.RepoName, os.ModeDir)
	log.Println("create README file for " + subRepoAnswer.RepoName)
	repoCreator.CreateRepo(subRepoAnswer)
	run(exec.Command("git", "init"), mainRepoName+"/"+subRepoAnswer.RepoName)
	ioutil.WriteFile(mainRepoName+"/"+subRepoAnswer.RepoName+"/README", []byte(""), 0644)
	run(exec.Command("git", "add", "."), mainRepoName+"/"+subRepoAnswer.RepoName)
	run(exec.Command("git", "commit", "-m", "\"init\""), mainRepoName+"/"+subRepoAnswer.RepoName)
	parse, _ := url.Parse(subRepoAnswer.GitHostAddress)
	gitRepoPath := "git@" + parse.Host + ":" + subRepoAnswer.RepoNamespace + "/" + subRepoAnswer.RepoName + ".git"
	run(exec.Command("git", "remote", "add", "origin", gitRepoPath), mainRepoName+"/"+subRepoAnswer.RepoName)
	run(exec.Command("git", "push", "-u", "origin", "master"), mainRepoName+"/"+subRepoAnswer.RepoName)
	var subRepoFolderName string
	if strings.HasSuffix(subRepoAnswer.RepoName, "vendor") {
		subRepoFolderName = "vendor"
	} else {
		subRepoFolderName = subRepoAnswer.RepoName
	}
	run(exec.Command("git", "submodule", "add", gitRepoPath, subRepoFolderName), mainRepoName)
}

func run(cmd *exec.Cmd, dir string) {
	fmt.Println(cmd.Args)
	cmd.Dir = dir
	e := cmd.Run()
	if e != nil {
		fmt.Println(e)
	}
}