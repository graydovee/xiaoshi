package aichat_plugin

import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
)

func Register() {
	engine := control.AutoRegister(&ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "AI聊天",
		Help:             "AI聊天",
	})

	engine.OnCommand("config").
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			GetBot().OnCommand(ctx)
		})

	engine.OnMessage(zero.OnlyToMe).
		SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			GetBot().OnMessage(ctx)
		})
}
