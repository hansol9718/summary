docker rm -f gateway

docker pull hansol9718/gateway


docker run \
    -d \
    -v /etc/letsencrypt:/etc/letsencrypt:ro \
    -e ADDR:=443 \
    -p 443:443\
    --name gateway \
    -e TLSKEY=/etc/letsencrypt/live/api.hansol7.me/privkey.pem \
    -e TLSCERT=/etc/letsencrypt/live/api.hansol7.me/fullchain.pem \
    hansol9718/gateway
exit