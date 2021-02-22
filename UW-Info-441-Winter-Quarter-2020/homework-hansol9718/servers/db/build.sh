docker build -t hansol9718/usersdb .

docker push hansol9718/usersdb

docker rm -f usersdb

docker run -d \
-p 3306:3306 \
--name usersdb \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=db \
hansol9718/usersdb
