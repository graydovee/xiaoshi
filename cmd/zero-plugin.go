package cmd

import (
	"fmt"
	aichatplugin "git.graydove.cn/graydove/xiaoshi.git/pkg/aichat-plugin"
	"git.graydove.cn/graydove/xiaoshi.git/pkg/config"
	webctrl "github.com/FloatTech/zbputils/control/web"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	"strings"

	// ---------以下插件均可通过前面加 // 注释，注释后停用并不加载插件--------- //
	// ----------------------插件优先级按顺序从高到低---------------------- //
	//                                                                  //
	//                                                                  //
	//                                                                  //
	//                                                                  //
	//                                                                  //
	// ----------------------------高优先级区---------------------------- //
	// vvvvvvvvvvvvvvvvvvvvvvvvvvvv高优先级区vvvvvvvvvvvvvvvvvvvvvvvvvvvv //
	//               vvvvvvvvvvvvvv高优先级区vvvvvvvvvvvvvv               //
	//                      vvvvvvv高优先级区vvvvvvv                      //
	//                          vvvvvvvvvvvvvv                          //
	//                               vvvv                               //

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/antiabuse" // 违禁词

	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/chat" // 基础词库

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/chatcount" // 聊天时长统计

	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/sleepmanage" // 统计睡眠时间

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/atri" // ATRI词库

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/manager" // 群管

	_ "github.com/FloatTech/zbputils/job" // 定时指令触发器

	//                               ^^^^                               //
	//                          ^^^^^^^^^^^^^^                          //
	//                      ^^^^^^^高优先级区^^^^^^^                      //
	//               ^^^^^^^^^^^^^^高优先级区^^^^^^^^^^^^^^               //
	// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^高优先级区^^^^^^^^^^^^^^^^^^^^^^^^^^^^ //
	// ----------------------------高优先级区---------------------------- //
	//                                                                  //
	//                                                                  //
	//                                                                  //
	//                                                                  //
	//                                                                  //
	// ----------------------------中优先级区---------------------------- //
	// vvvvvvvvvvvvvvvvvvvvvvvvvvvv中优先级区vvvvvvvvvvvvvvvvvvvvvvvvvvvv //
	//               vvvvvvvvvvvvvv中优先级区vvvvvvvvvvvvvv               //
	//                      vvvvvvv中优先级区vvvvvvv                      //
	//                          vvvvvvvvvvvvvv                          //
	//                               vvvv                               //

	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/ahsai"            // ahsai tts
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/aifalse" // 服务器监控
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/aiwife"  // 随机老婆
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/alipayvoice"      // 支付宝到账语音
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/autowithdraw" // 触发者撤回时也自动撤回
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/baiduaudit"       // 百度内容审核
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/base16384"   // base16384加解密
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/base64gua"   // base64卦加解密
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/baseamasiro" // base天城文加解密
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/bilibili"         // b站相关
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/bookreview"       // 哀伤雪刃吧推书记录
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/chess"            // 国际象棋
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/choose"           // 选择困难症帮手
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/chouxianghua"     // 说抽象话
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/chrev"            // 英文字符翻转
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/coser"            // 三次元小姐姐
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/cpstory"          // cp短打
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/dailynews" // 今日早报
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/danbooru"         // DeepDanbooru二次元图标签识别
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/diana" // 嘉心糖发病
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/dish"  // 程序员做饭指南
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/drawlots"         // 多功能抽签
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/driftbottle"      // 漂流瓶
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/emojimix"   // 合成emoji
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/event"      // 好友申请群聊邀请事件处理
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/font"       // 渲染任意文字到图片
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/fortune"    // 运势
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/funny"      // 笑话
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/genshin"    // 原神抽卡
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/gif"        // 制图
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/github"     // 搜索GitHub仓库
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/guessmusic" // 猜歌
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/hitokoto"         // 一言
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/hs"               // 炉石
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/hyaku"            // 百人一首
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/inject" // 注入指令
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/jandan"           // 煎蛋网无聊图
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/jptingroom"       // 日语听力学习材料
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/kfccrazythursday" // 疯狂星期四
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/lolicon"          // lolicon 随机图片
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/lolimi"           // 桑帛云 API
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/magicprompt"  // magicprompt吟唱提示
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/mcfish"       // 钓鱼模拟器
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/midicreate"   // 简易midi音乐制作
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/moegoe"       // 日韩 VITS 模型拟声
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/moyu"         // 摸鱼
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/moyucalendar" // 摸鱼人日历
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/music"        // 点歌
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/nativesetu"   // 本地涩图
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/nbnhhsh" // 拼音首字母缩写释义工具
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/nihongo"     // 日语语法学习
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/novel"       // 铅笔小说网搜索
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/nsfw"        // nsfw图片识别
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/nwife"       // 本地老婆
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/omikuji" // 浅草寺求签
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/poker"       // 抽扑克
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/qqwife"      // 一群一天一夫一妻制群老婆
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/qzone"       // qq空间表白墙
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/realcugan"   // realcugan清晰术
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/reborn"      // 投胎
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/robbery"     // 打劫群友的ATRI币
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/runcode"  // 在线运行代码
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/saucenao" // 以图搜图
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/score"       // 分数
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/setutime"    // 来份涩图
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/shadiao"     // 沙雕app
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/shindan" // 测定
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/steam"       // steam相关
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/tarot"    // 抽塔罗牌
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/tiangou"  // 舔狗日记
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/tracemoe" // 搜番
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/translation" // 翻译
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/vitsnyaru"   // vits猫雷
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/wallet"      // 钱包
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/wantquotes"  // 据意查句
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/warframeapi" // warframeAPI插件
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/wenxinvilg"  // 百度文心AI画图
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/wife"        // 抽老婆
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/wordcount"   // 聊天热词
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/wordle" // 猜单词
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/ygo"         // 游戏王相关插件
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/ymgal"       // 月幕galgame
	//_ "github.com/FloatTech/ZeroBot-Plugin/plugin/yujn"        // 遇见API

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/process"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	file.SkipOriginal = true

}

func Run(cfg *config.Config) {
	aichatplugin.InitAIBot(cfg)

	go webctrl.RunGui(fmt.Sprintf("%s:%d", cfg.QQBot.WebGui.Host, cfg.QQBot.WebGui.Port))

	var drivers []zero.Driver
	if cfg.QQBot.Ws != nil {
		wsUrl := fmt.Sprintf("ws://%s:%d", strings.TrimPrefix(cfg.QQBot.Ws.Addr, "ws://"), cfg.QQBot.Ws.Port)
		drivers = append(drivers, driver.NewWebSocketClient(wsUrl, cfg.QQBot.Ws.Token))
	}

	zeroCfg := &zero.Config{
		NickName:      cfg.QQBot.Zero.NickNames,
		CommandPrefix: "/",
		SuperUsers:    cfg.QQBot.Zero.SuperUsers,
		Driver:        drivers,
	}
	zero.RunAndBlock(zeroCfg, process.GlobalInitMutex.Unlock)
}
