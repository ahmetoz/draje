# draje

```shell
NAME:
   draje - docker registry api from jenkins - A new cli application

USAGE:
   draje.exe [global options] command [command options] [arguments...]

VERSION:
   0.0.1

DESCRIPTION:
   uses docker registy v2 api for deleting images from docker registry

AUTHOR:
   Ahmet Oz <bilmuhahmet@gmail.com>

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --host value           host address of docker registry
   --image value          name of the image
   --excluded_tags value  name of the image tag which will be excluded (default: "latest")
   --username value       registry user name
   --password value       registry user password
   --exclude_last value   exclude last n images (default: 1)
   --help, -h             show help
   --version, -v          print the version
```

## building the code 

``` shell
cd draje
go install
```
- not: go install- builds and installs go code; if code is package main, it creates a binary executable and drops it in the workspace’s **bin** folder;

- or you can download executables from [relases](https://github.com/ahmetoz/draje/releases)
