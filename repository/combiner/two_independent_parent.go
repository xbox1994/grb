package combiner

import (
	"fmt"
	"grb/common/util"
	"grb/repository/creater"
	"grb/repository/loginer"
	"grb/repository/model"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
)

type TwoIndependentParentCombiner struct {
}

func (r TwoIndependentParentCombiner) CreateAndCombine(
	repoCreator creater.RepoCreator,
	repoCreatePreInfo loginer.RepoCreatePreInfo,
	answers model.Answer) {
	// 在远端与本地同时创建所有Repo
	mainRepoName := answers.RepoName
	os.Mkdir(mainRepoName, os.ModeDir)
	repoCreator.CreateRemoteRepo(answers, repoCreatePreInfo)
	util.Run(exec.Command("git", "init"), mainRepoName)

	answers.RepoName = mainRepoName + "-admin"
	createSubRepo(repoCreator, repoCreatePreInfo, answers, mainRepoName)
	answers.RepoName = mainRepoName + "-server"
	createSubRepo(repoCreator, repoCreatePreInfo, answers, mainRepoName)

	// 将子Repo加入到父Repo中
	util.Run(exec.Command("git", "add", "."), mainRepoName)
	util.Run(exec.Command("git", "commit", "-m", "\"init\""), mainRepoName)
	parse, _ := url.Parse(answers.GitHostAddress)
	gitRepoPath := "git@" + parse.Host + ":" + answers.RepoNamespace + "/" + mainRepoName + ".git"
	util.Run(exec.Command("git", "remote", "add", "origin", gitRepoPath), mainRepoName)
	util.Run(exec.Command("git", "push", "-u", "origin", "master"), mainRepoName)
}

func createSubRepo(
	repoCreator creater.RepoCreator,
	repoCreatePreInfo loginer.RepoCreatePreInfo,
	subRepoAnswers model.Answer,
	mainRepoName string) {
	var subRepoFolderName string
	subRepoFolderName = subRepoAnswers.RepoName
	subRepoFolderPath := mainRepoName + "/" + subRepoFolderName
	os.Mkdir(subRepoFolderPath, os.ModeDir)
	fmt.Println("create README file for " + subRepoAnswers.RepoName)
	repoCreator.CreateRemoteRepo(subRepoAnswers, repoCreatePreInfo)
	util.Run(exec.Command("git", "init"), subRepoFolderPath)
	ioutil.WriteFile(subRepoFolderPath+"/README", []byte(""), 0644)
	util.Run(exec.Command("git", "add", "."), subRepoFolderPath)
	util.Run(exec.Command("git", "commit", "-m", "\"init\""), subRepoFolderPath)
	parse, _ := url.Parse(subRepoAnswers.GitHostAddress)
	gitRepoPath := "git@" + parse.Host + ":" + subRepoAnswers.RepoNamespace + "/" + subRepoAnswers.RepoName + ".git"
	util.Run(exec.Command("git", "remote", "add", "origin", gitRepoPath), subRepoFolderPath)
	util.Run(exec.Command("git", "push", "-u", "origin", "master"), subRepoFolderPath)
	util.Run(exec.Command("git", "submodule", "add", gitRepoPath, subRepoFolderName), mainRepoName)
}
