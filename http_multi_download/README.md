# program: multi_download
# author:  xiongyanfu
# date:    2024-06-21

这两天用libcurl包实现多线程下载文件功能，经MD5和SHA校验测试下载功能通过，但是存在以下问题：
1、在网络时延大或者下载超大文件时仍然会有报错，还需要进一步分析原因；
2、如果有一个下载任务报错，没有立即终止所有线程，而是要等到其他进行中的线程任务结束才能全部结束。



test: 

// ./multi_download  https://releases.ubuntu.com/22.04/ubuntu-22.04.4-live-server-amd64.iso.zsync  ubuntu-22.04.4.iso.zsync

//https://releases.ubuntu.com/22.04/ubuntu-22.04.4-live-server-amd64.iso.zsync  	3.9M
//https://releases.ubuntu.com/24.04/ubuntu-24.04-netboot-amd64.tar.gz				81M
//https://releases.ubuntu.com/23.10.1/ubuntu-23.10-netboot-amd64.tar.gz             112M
//https://releases.ubuntu.com/22.04/ubuntu-22.04.4-live-server-amd64.iso			2G
//https://releases.ubuntu.com/23.10.1/ubuntu-23.10.1-desktop-amd64.iso              4.8G   


https://mirror.elice.io/centos/7.9.2009/isos/x86_64/CentOS-7-x86_64-NetInstall-2009.iso



