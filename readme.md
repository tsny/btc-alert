## Usage

`go get github.com/tsny/btc-alert`

Follow the configuration section below, go build, and run it

## Configuration

`config.json` houses each 'interval' which represents an interval that is checked every minute.
If an interval lapses, then a notification is sent if the absolute value of the overall percentage change in the asset is 
more than the `percentThreshold` field.

`occurence`s are minutes.

Config file will look like this: 

```json
{
    "bootNotification": false,
    "discord": {
        "useDiscord": false,
        "token": "token.WgvCew.ObvImI8RA0HAn-kf2u8omFQGINI",
        "channelId": "715135158184218",
        "clearOnBoot": false
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


### Desktop Notifications

Example:

![Example Screenshot](https://i.imgur.com/VZ6bHZU.png)

Notifications work on both Mac and Windows

### Discord

If you have a Discord bot, you can redirect to there!

![Example Screenshot](https://i.imgur.com/zbDfI9B.png)

## Example

![Example Screenshot](https://i.imgur.com/lKS8kzG.png)

## Dependencies

- github.com/gen2brain/beeep     -  Notifications
- github.com/bwmarrin/discordgo  -  Discord API
- Yahoo Finance                  -  BTC Price Data
- CoindeskAPI                    -  BTC Price Data
