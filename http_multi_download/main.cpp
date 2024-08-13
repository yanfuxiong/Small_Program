#include "DownLoadFactory.h"
#include <stdio.h>


// ./multi_download  https://releases.ubuntu.com/22.04/ubuntu-22.04.4-live-server-amd64.iso.zsync  ubuntu-22.04.4.iso.zsync

int main(int argc, char *argv[])
{

    if(argc != 3 && argc != 4)
    {
        printf("arg num error \n");
        printf("arg 1: url \n");
        printf("arg 2: target file name \n");
        printf("arg 3: thread num,  default 4 \n");
        
		printf("please input like this : \"./multi_download  url  filename  threadnum\" \n");
		return -1;
	}

	DownLoadFactory dFact(argv[1]);
	dFact.MakeTask(argv[2]);
    
    if (argc == 3)
        dFact.init();
    else
        dFact.init(atoi(argv[3]));

	dFact.JoinParser();
	
	return 0;
}