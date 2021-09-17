#!/bin/sh

ssh-keygen -t rsa -b 1024 -f rs256.key
# Don't add passphrase
openssl rsa -in rs256.key -pubout -outform PEM -out rs256.key.pub
cat rs256.key
cat rs256.key.pub
