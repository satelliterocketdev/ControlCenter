## How to start

To run I68Social API for local development you need Docker and Docker-Compose

For local development you need start next services:
- app (Golang app) from Dockerfile
- db (MySQL server)


Next command setup **i69social** (Golang App) container with hot reloading and debugger and MySQL DB server
```bash
$ docker-compose -f ./docker-compose-DEV.yml up -d --build app db
```
Now you can access Food API using URL: http://localhost:3003


If you need only setup I69social API for using in local environment (without development), run next command:
```bash
$ docker-compose -f ./docker-compose.yml up -d --build app db
```
Now you can access Food API using URL: http://localhost:3000

All local environment variables you can find in file **.env**

---

## How to create and apply SQL migrations

First you need enter inside container
```bash
$ docker exec -it food-api bash
```

#### Create new SQL migration
To create new migration, you can use following command:
```bash
# /app/bin/migrate create -ext sql -dir ./migrations NAME
```

Now you can find generated files in folder **migrations**
- {digits}_NAME.up.sql
- {digits}_NAME.down.sql

#### Apply migrations
To apply all existing migrations, you can use following command:
```bash
# /app/bin/migrate -source file:///app/migrations/ -database "mysql://$MYSQL_USER:$MYSQL_PASSWORD@tcp($MYSQL_HOST:$MYSQL_PORT)/$MYSQL_DATABASE" up
```

If you want to see current DB version, you can use following command:
```bash
# /app/bin/migrate -source file:///app/migrations/ -database "mysql://$MYSQL_USER:$MYSQL_PASSWORD@tcp($MYSQL_HOST:$MYSQL_PORT)/$MYSQL_DATABASE" version
```

## Get credentials
Creds for Admin SDK
1. Go to project settings
2. Select service accounts
3. Open link `Manage access for service accounts`
4. Create JSON key for `firebase-adminsdk` and save it on your computer
5. Copy and paste content of downloaded file to `src/api/config/firebase.json`
6. Open Database page on firebase console
7. Copy Database URL and FIREBASE_DB_URL in .env file


## Deploy
```bash
$ docker build --no-cache --target base -t i69social-builder .
// prod
$ docker build --cache-from i69social-builder --rm --target prod -t i69social .

// alternative build
$ docker-compose -f ./docker-compose.yml build --force-rm --no-cache --pull app

// run new version
$ docker-compose up --no-deps -d app

// run new version
$ docker-compose -f ./docker-compose.yml -f ./docker-compose-DEV.yml up --no-deps -d app
```

## Setup as service

$ sudo nano /etc/systemd/system/docker-compose-app.service

[Unit]
Description=Docker Compose Application Service
Requires=docker.service
After=docker.service

[Service]
WorkingDirectory=/usr/sites/i69social
ExecStart=/usr/bin/docker-compose up
ExecStop=/usr/bin/docker-compose stop
TimeoutStartSec=0
Restart=on-failure
StartLimitIntervalSec=60
StartLimitBurst=3

[Install]
WantedBy=multi-user.target

// enable service to auto start after server reboot
$ sudo systemctl enable docker-compose-app
$ sudo systemctl start docker-compose-app.service

Renew lets encrypt certificates
```bash
cd /etc/letsencrypt/ && ./certbot-auto renew --dry-run
cd /etc/letsencrypt/ && ./certbot-auto renew
wget https://dl.eff.org/certbot-auto && chmod a+x certbot-auto
sudo mv certbot-auto /etc/letsencrypt/
sudo crontab -e
```
