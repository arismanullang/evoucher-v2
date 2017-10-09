# e-voucher
----------
### Requerment Summary:

 * GO v1.8 or later varsion
 * Vendoring Using [Govendor](github.com/kardianos/govendor)


### Instalation :
* Read Go installation steps from [hire](golang.org/doc/install).
* After Go installed and the environment variables are set, install govendor

```sh
go get -u github.com/kardianos/govendor
```


### Project
* Clone the project
```sh
git clone https://github.com/gilkor/evoucher.git
```
* Restore vendor source , (full documentation [hire](github.com/kardianos/govendor/blob/master/doc/dev-guide.md))
```sh
govendor sync
```
* Setup Environment Variable to config file ([./files/etc/voucher/config.yml](github.com/gilkor/evoucher/blob/master/files/etc/voucher/config.yml)) or move file to any directory as you will. for example :
```sh
EVOUCHER_CONFIG=/etc/evoucher/config.yml
```
* Type below command To run application
```sh
go build && ./evoucher
```
or you can direct output binary , for example :
```sh
go build -o /usr/local/bin/evoucher
```
----
### api documentation

* postman (coming soon)
* [swagger](swagger.io) (coming soon)

