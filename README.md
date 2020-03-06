# Voucher V2
----------
### Requerment Summary:

 * GO v1.12.0 or later varsion

### Instalation :
* Read Go installation steps from [hire](golang.org/doc/install).

* Go to your repository dir , then clone the project
```sh
go clone https://github.com/gilkor/evoucher-v2.git
```
* This project are using [GoModules](https://github.com/golang/go/wiki/Modules) , hit command bellow to get all dependenciess and create `./vendor` dir.
```sh
go mod vendor
```


* Environment Variable using [Godotenv](https://github.com/joho/godotenv) , Add your application configuration to your `.env` file in the root of your project, use file `.env.sample` as template.

### Run Project
* type command below to run the project
```sh
go build && ./evoucher
```
* or you can direct output binary , for example :
```sh
go build -o /usr/local/bin/evoucher
```
----
### api documentation
* [swagger](https://swaggerhub.com/apis/malfanmh/e-voucher/v1) (coming soon)
