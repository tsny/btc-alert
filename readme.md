## Usage

Build and run main.go

## Configuration

`config.json` houses each 'interval' which represents an interval that is checked every minute.
If an interval lapses, then a notification is sent if the absolute value of the overall percentage change in the asset is 
more than the `percentThreshold` field.

`occurence`s are minutes.

## Example

![Example Screenshot](https://i.imgur.com/lKS8kzG.png)
