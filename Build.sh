#!/bin/sh
watcher -cmd="sh Update.sh" -recursive -pipe=true -list ./control &
canthefason_watcher -run pikabu-control/control/cmd -watch pikabu-control
