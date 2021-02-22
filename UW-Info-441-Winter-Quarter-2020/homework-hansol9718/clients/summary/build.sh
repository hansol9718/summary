docker build -t hansol9718/summary .


docker push hansol9718/summary
ssh -i ~/.ssh/MyKeyPair.pem ec2-user@ec2-52-91-212-232.compute-1.amazonaws.com < deploy.sh