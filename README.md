# SiteParser

Website scraping application written in golang.

Suitable for parsing sites with a structure similar to https://itproger.com or https://habr.com/ru/
a feature is the setting of parsing rules using the http interface. And searching scrapped articles by http handler.

# Launch
docker, docker-compose are required to run.

Run ./run.sh script.

http server will run on port :8080
mysql :3308
mysql_user: root
mysql_password: root

# Example using

An example of starting parsing is written in ./index.html.

Searching http://localhost:8080:/api/search/{searchValue}

will be given in json format 

