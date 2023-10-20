# winecontainer 后台程序

## 介绍

用于维护 wine 容器的后台程序，负责与调度器之间交互，并执行关于容器的控制命令

## 软件架构

基于 go 语言开发的后台程序

## 安装教程

## 使用说明

需要设置 `SERVERD_API_URL` 与 `HANDLER_TIMEOUT` 两个环境变量，`SERVERD_API_URL` 默认为 `/api/serverd`

## 接口文档

`/api/serverd/createvm`

`/api/serverd/deletevm`

## 常见问题

## ToDo
