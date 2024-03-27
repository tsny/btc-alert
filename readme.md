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
        "token": "token.WgvCew.ObvImI8RA0HAn-kf2u8omFQGINI",
        "channelId": "715135158184218",
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

### API

Checkout `routes.go` for a few of the endpoints, port defaults to 8080

`GET /crypto`
Returns the publishers associated with all tracked crypto coins
![Example API](https://i.imgur.com/eMDdj3S.png)


![Example BTC Graph](https://i.imgur.com/qME6WLJ.png)

### Discord

If you have a Discord bot, you can redirect to there!

![Example Screenshot](https://i.imgur.com/zbDfI9B.png)
