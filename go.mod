module github.com/lbryio/dispendium

go 1.13

replace github.com/btcsuite/btcd => github.com/lbryio/lbrycrd.go v0.0.0-20200203050410-e1076f12bf19

require (
	github.com/btcsuite/btcd v0.0.0-20190213025234-306aecffea32
	github.com/btcsuite/btcutil v0.0.0-20190425235716-9e5f4b9a998d
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/fatih/color v1.7.0
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gofrs/uuid v3.2.0+incompatible // indirect
	github.com/jasonlvhit/gocron v0.0.0-20191228163020-98b59b546dee
	github.com/jmoiron/sqlx v1.2.0
	github.com/johntdyer/slack-go v0.0.0-20180213144715-95fac1160b22 // indirect
	github.com/johntdyer/slackrus v0.0.0-20180518184837-f7aae3243a07
	github.com/jteeuwen/go-bindata v3.0.7+incompatible
	github.com/lbryio/lbry.go v1.1.2
	github.com/lbryio/lbry.go/v2 v2.4.7-0.20200203053542-c4772e61c565
	github.com/lbryio/ozzo-validation v0.0.0-20170512160344-202201e212ec
	github.com/mitchellh/mapstructure v1.2.2 // indirect
	github.com/pelletier/go-toml v1.7.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/pkg/profile v1.4.0
	github.com/prometheus/client_golang v0.9.3
	github.com/rubenv/sql-migrate v0.0.0-20200212082348-64f95ea68aa3
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.6.3
	github.com/volatiletech/inflect v0.0.0-20170731032912-e7201282ae8d // indirect
	github.com/volatiletech/null v8.0.0+incompatible
	github.com/volatiletech/sqlboiler v3.4.0+incompatible
	golang.org/x/sys v0.0.0-20200420163511-1957bb5e6d1f // indirect
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0 // indirect
	gopkg.in/ini.v1 v1.55.0 // indirect
)
