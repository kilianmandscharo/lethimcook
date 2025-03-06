#!/bin/bash

while true; 
    do inotifywait --event modify --exclude \.test\.ts$ ./src/ \
        && npm run build;
done
