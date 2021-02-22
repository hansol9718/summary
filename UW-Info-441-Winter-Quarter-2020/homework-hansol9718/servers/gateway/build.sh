GOOS=linux go build
docker build -t hansol9718/gateway .
go clean

docker push hansol9718/gateway
ssh -i ~/.ssh/MyKeyPair.pem ec2-user@ec2-3-82-136-220.compute-1.amazonaws.com < deploy.sh



