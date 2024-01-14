@echo off
set PROJECT_NAME=shiro_killer_mod

REM 设置go语言的交叉编译参数，根据你的目标平台和架构进行修改
set GOOS=windows
set GOARCH=amd64
go build -o %PROJECT_NAME%_windows_amd64.exe

REM 重复上述步骤，修改GOOS和GOARCH的值，以适应不同的平台和架构
set GOOS=linux
set GOARCH=amd64
go build -o %PROJECT_NAME%_linux_amd64

set GOOS=darwin
set GOARCH=amd64
go build -o %PROJECT_NAME%_macos_amd64

set GOOS=windows
set GOARCH=386
go build -o %PROJECT_NAME%_windows_386.exe

set GOOS=linux
set GOARCH=386
go build -o %PROJECT_NAME%_linux_386

set GOOS=darwin
set GOARCH=386
go build -o %PROJECT_NAME%_macos_386

set GOOS=windows
set GOARCH=arm
go build -o %PROJECT_NAME%_windows_arm.exe

set GOOS=linux
set GOARCH=arm
go build -o %PROJECT_NAME%_linux_arm

set GOOS=darwin
set GOARCH=arm64
go build -o %PROJECT_NAME%_macos_arm

echo 编译完成
pause
