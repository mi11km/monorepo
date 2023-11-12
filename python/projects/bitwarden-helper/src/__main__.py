#!/usr/bin/env python
import logging
import os

from slack_bolt import App
from slack_bolt.adapter.socket_mode import SocketModeHandler

logging.basicConfig(level=logging.DEBUG)

slack_bot_token = os.environ["SLACK_BOT_TOKEN"]
slack_app_token = os.environ["SLACK_APP_TOKEN"]

app = App(token=slack_bot_token)


@app.message("hi")
def message_hello(message, say):
    say(blocks=[{
        "type": "section",
        "text": {
            "type": "mrkdwn",
            "text": f"Hey there <@{message['user']}>!"
        },
        "accessory": {
            "type": "button",
            "text": {
                "type": "plain_text",
                "text": "Click Me"
            },
            "action_id": "button_click"
        }
    }],
        text=f"Hey there <@{message['user']}>!")


@app.action("button_click")
def action_button_click(body, ack, say):
    # アクションを確認したことを即時で応答します
    ack()
    # チャンネルにメッセージを投稿します
    say(f"<@{body['user']['id']}> clicked the button")


if __name__ == "__main__":
    SocketModeHandler(app, slack_app_token).start()
