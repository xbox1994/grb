package repository

import (
	"errors"
	"fmt"
	"gopkg.in/AlecAivazis/survey.v1"
	"grb/common/project_type"
	"grb/repository/combiner"
	"grb/repository/creater"
	"grb/repository/loginer"
	"grb/repository/model"
	"os"
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
		Prompt:   &survey.Input{Message: "Repository namespace:"},
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

func Create(projectStructure string) string {
	answers := model.Answer{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	//// 测试数据
	//answers = model.Answer{
	//	GitHostAddress:   "http://wpsgit.kingsoft.net/",
	//	GitServerVersion: "GitLab 6.3.0 LDAP",
	//	RepoName:         "grbtest",
	//	RepoNamespace:    "galaxy",
	//	Username:         "wangtianyi1",
	//	Password:         "",
	//}

	// 选择creator, loginer
	var repoCreator creater.RepoCreator
	var gitLoginer loginer.GitWebInterfaceLoginer
	if "GitLab 6.3.0 LDAP" == answers.GitServerVersion {
		gitLoginer = &loginer.Gitlab630Ldap{}
		repoCreator = &creater.Gitlab630Ldap{}
	} else {
		panic(errors.New(answers.GitServerVersion + " no implement yet"))
	}

	// 登录获取到Cookie与RepoNamespaceId
	repoCreatePreInfo := gitLoginer.Login(model.LoginInfo{
		GitHostAddress: answers.GitHostAddress,
		Username:       answers.Username,
		Password:       answers.Password,
		RepoNamespace:  answers.RepoNamespace,
	})

	var repoCombiner combiner.RepoCombiner
	switch projectStructure {
	case project_type.OneIndependent:
		repoCombiner = combiner.SingleCombiner{}
	case project_type.TwoIndependent:
		repoCombiner = combiner.TwoIndependentCombiner{}
	case project_type.TwoIndependentWithParent:
		repoCombiner = combiner.TwoIndependentParentCombiner{}
	}

	// 在远端与本地创建父项目、子项目，合并子项目到父项目
	repoCombiner.CreateAndCombine(repoCreator, repoCreatePreInfo, answers)

	return answers.RepoName
}
