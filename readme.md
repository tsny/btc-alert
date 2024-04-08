`btc-alert` is a Discord bot (or desktop app) that can alert you when price movements occur for crypto. (mostly BTC)

## Usage

`go get github.com/tsny/btc-alert`

Follow the configuration section below, go build, and run it

## Configuration

Create a `config.json` file in the root of the project.

`config.json` houses each 'interval' which represents an interval that is checked every minute.
If an interval lapses, then a notification is sent if the absolute value of the overall percentage change in the asset is 
more than the `percentThreshold` field.

`occurence`s are minutes.

Config file will look like this: 

```json
{
    "bootNotification": false,
    "discord": {
        "enabled": false,
        "token": "token.WgvCew.ObvImI8faketoken-kf2u8omFQGINI",
        "usersToNotify": [
            "84090395092353024"
        ]
    },
    "thresholds": [
        {
            "threshold": 500
        }
    ],
    "intervals": [
        {
            "maxOccurences": 5,
            "percentThreshold": 1
        },
        {
            "maxOccurences": 20,
            "percentThreshold": 2
        },
        {
            "maxOccurences": 60,
            "percentThreshold": 3
        }
    ]
}
```

Example of using the bot via discord

![exmaple-discordbot](https://i.postimg.cc/x8VP962Y/image.png)