package aichat_plugin

import (
	"git.graydove.cn/graydove/xiaoshi.git/pkg/chatgpt"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
)

const dataUrl = "https://git.graydove.cn/graydove/xiaoshi/raw/branch/master/"

const (
	fileCharacter = "character.yaml"
)

func Register(bot *AIBot) {
	engine := control.AutoRegister(&ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "XiaoShiAI",
		Help:             "(@)小诗[对话内容]",
		PublicDataFolder: "XiaoShi",
	})

	data, err := engine.GetCustomLazyData(dataUrl, fileCharacter)
	if err != nil {
		panic(err)
	}
	bot.prompt = chatgpt.MustLoadRole(engine.DataFolder()+"/"+fileCharacter, data)

	engine.OnCommand("config").
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			bot.OnCommand(ctx)
		})

	engine.OnMessage(zero.OnlyToMe).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			bot.OnMessage(ctx)
		})
}
