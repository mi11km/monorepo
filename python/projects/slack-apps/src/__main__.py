#!/usr/bin/env python
import logging

from slack_bolt import App

if __name__ == "__main__":
    logging.basicConfig(level=logging.DEBUG)
    app = App()
    app.start(3000)  # POST http://localhost:3000/slack/events
