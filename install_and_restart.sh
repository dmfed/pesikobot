#!/bin/bash
systemctl stop pesikobot.service
cp pesikobot /usr/local/bin
systemctl start pesikobot.service

