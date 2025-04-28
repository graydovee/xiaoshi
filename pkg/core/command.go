package core

import (
	"bytes"
	"fmt"
	"github.com/graydovee/xiaoshi/pkg/chat"
	"github.com/graydovee/xiaoshi/pkg/util"
	"github.com/spf13/cobra"
	"log/slog"
	"strconv"
)

type CommandFactory func(bot *OneBotCore, chatSession *chat.Session, root *cobra.Command)

func (s *OneBotCore) BuildCommand(chatSession *chat.Session) *cobra.Command {
	root := cobra.Command{}
	for _, cmdFactory := range []CommandFactory{
		//CharacterCommand,
		ChatCommand,
	} {
		cmdFactory(s, chatSession, &root)
	}
	return &root
}

func RunCmd(command *cobra.Command, arguments []string) (string, error) {
	out := bytes.NewBuffer(nil)
	command.SetArgs(arguments)
	command.SetOut(out)

	if err := command.Execute(); err != nil {
		slog.Error("execute sub command error: ", err)
		return "", err
	}
	return out.String(), nil
}

//func CharacterCommand(bot *OneBotCore, chatSession *chat.Session, root *cobra.Command) {
//	characterCmd := &cobra.Command{
//		Use: "character [subCommand]",
//	}
//	characterCmd.AddCommand(&cobra.Command{
//		Use:   "list",
//		Short: "查询预设人格",
//		Long:  "查询预设人格",
//		RunE: func(cmd *cobra.Command, args []string) error {
//			p := util.NewPrinter(cmd.OutOrStdout())
//			p.Println("预设人格列表：")
//			c := 1
//			for role := range bot.prompt.Characters {
//				p.Println(c, ". ", role)
//				c++
//			}
//			return nil
//		},
//	})
//	characterCmd.AddCommand(&cobra.Command{
//		Use:   "use [name]",
//		Short: "切换至预设人格",
//		Long:  "切换至预设人格",
//		RunE: func(cmd *cobra.Command, args []string) error {
//			p := util.NewPrinter(cmd.OutOrStdout())
//			if len(args) == 0 {
//				return fmt.Errorf("角色名为空")
//			}
//			rolePrompts, ok := bot.prompt.GetRolePrompt(args[0])
//			if ok {
//				chatSession.SetPrompt(rolePrompts...)
//				p.Println("角色切换至：", args[0])
//			} else {
//				p.Printf("角色: %s 不存在\n", args[0])
//			}
//			return nil
//		},
//	})
//	characterCmd.AddCommand(&cobra.Command{
//		Use:   "add [name] [detail]",
//		Short: "新增设定人格至预设",
//		Long:  "新增设定人格至预设",
//		RunE: func(cmd *cobra.Command, args []string) error {
//			p := util.NewPrinter(cmd.OutOrStdout())
//			if len(args) < 2 {
//				return fmt.Errorf("设定为空")
//			}
//			roleName := args[0]
//			roleDetail := strings.Join(args[1:], "\n")
//			bot.prompt.SetRolePrompt(roleName, roleDetail)
//			p.Println("新增人格完成")
//			return nil
//		},
//	})
//	characterCmd.AddCommand(&cobra.Command{
//		Use:   "del [name]",
//		Short: "删除设定人格",
//		Long:  "删除设定人格",
//		RunE: func(cmd *cobra.Command, args []string) error {
//			p := util.NewPrinter(cmd.OutOrStdout())
//			if len(args) == 0 {
//				return fmt.Errorf("设定为空")
//			}
//			roleName := args[0]
//			bot.prompt.DeleteRolePrompt(roleName)
//			p.Println("删除角色完成")
//			return nil
//		},
//	})
//
//	root.AddCommand(characterCmd)
//}

func ChatCommand(bot *OneBotCore, chatSession *chat.Session, root *cobra.Command) {
	chatCmd := cobra.Command{
		Use: "chat [subCommand]",
	}

	chatCmd.AddCommand(&cobra.Command{
		Use:   "limit [chatLength]",
		Short: "设置对话长度限制",
		Long:  "设置对话长度限制",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := util.NewPrinter(cmd.OutOrStdout())
			if len(args) != 1 {
				return fmt.Errorf("参数数量错误")
			}
			limit, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("参数错误: %w", err)
			}
			chatSession.History().SetLimit(limit)
			p.Println("设置对话长度限制为: ", limit)
			return nil
		},
	})

	chatCmd.AddCommand(&cobra.Command{
		Use:   "clear",
		Short: "清空对话历史",
		Long:  "清空对话历史",
		Run: func(cmd *cobra.Command, args []string) {
			p := util.NewPrinter(cmd.OutOrStdout())
			chatSession.History().Clear()
			p.Println("对话历史已清空")
		},
	})

	chatCmd.AddCommand(&cobra.Command{
		Use:   "mode [modeName]",
		Short: "设置语言模型",
		Long:  "设置语言模型",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				slog.Error("参数错误")
				return
			}
			p := util.NewPrinter(cmd.OutOrStdout())
			chatSession.ChatBot().SetModel(args[0])
			p.Println("设置语言模型为: ", args[0])
		},
	})

	root.AddCommand(&chatCmd)
}
