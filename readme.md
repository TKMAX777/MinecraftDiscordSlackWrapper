# MinecraftDiscordWrapper
## 概要
DiscordにMinecraftの出力内容を垂れ流すのが目的なプログラム。ついでに許可したユーザに関してはsayを用いることが可能。

## 目次
<!-- TOC -->

- [MinecraftDiscordWrapper](#minecraftdiscordwrapper)
    - [概要](#概要)
    - [目次](#目次)
    - [設定](#設定)

<!-- /TOC -->

## 設定
このプログラムをserver.jarと同階層に配置。settings.jsonも次のように記述したものを配置する。

```json
{
    "Discord":{
        "GuildID": "用いるDiscordサーバのGuildID",
        "Token": "DiscordToken",
        "Default":{
            "AvaterURI": "BotのアイコンのURI。お好みで。",
            "UserName":"MinecraftWrapper",
            "HookURI": "https://discord.com/api/webhooks/***"
        },
        "Error":{
            "AvaterURI": "Error時のBotのアイコンのURI。お好みで。",
            "UserName":"MinecraftWrapper - Error -",
            "HookURI": "https://discord.com/api/webhooks/***(省略可。その場合Defaultが適用される)"
        }        
    }
}
```

お好みで次のように
`name_dict.json`
ファイルを同層に作成すると、Discordで
```
say Content
```
によってメッセージを飛ばすことが可能。


```json
[
    {
	"DiscordID":"DiscordUserID",
	"Name": "USER_NAME"
    }
]
```




