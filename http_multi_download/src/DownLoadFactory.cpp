#include "DownLoadFactory.h"


#define RANGE_SIZE	 (1024*1024 *10)		//以10M 为单位创建一个下载任务，超过10M的文件需要分段


Factory::~Factory()
{
}

DownLoadFactory::DownLoadFactory(const string &url,   const eFileLoadType type)
{

    curl_global_init(CURL_GLOBAL_ALL); 	

    this->m_url = url;
    this->m_eFileLoadType = type;
}

DownLoadFactory::~DownLoadFactory()
{

    fclose(this->m_fp);
    this->m_fp = NULL;
	
}


//工厂模式获取具体下载类型对象
TaskBase* DownLoadFactory::CreateTaskInstance()
{
    TaskBase* pTask = NULL;

    switch(m_eFileLoadType)
    {
        case FileLoad_Http_Curl:
            pTask = new CurlHttpTask();
            break;
        case FileLoad_Ftp:
            printf("this file load type: %d is unrealized! \n", m_eFileLoadType);
            break;
        default:
            printf("unknow file load type!\n");
            break;
				
    }
    return pTask;
}

void DownLoadFactory::MakeTask(const string &userfilename)
{
    long lFileLen = 0L, start, stop,seg_num;

    TaskBase *pTask = this->CreateTaskInstance();
    if (NULL == pTask)
        return;

    lFileLen = pTask->GetDownloadFileLength(this->GetUrl().c_str());
    delete pTask;
    pTask = NULL;
	
    if(lFileLen <= 0){
        printf("this %s file get len error\n", this->GetUrl().c_str());
        exit(1);
    }
    printf("file:%s  get len is %ld\n",  this->GetUrl().c_str(), lFileLen);

    string TargetFilePath = "./" + userfilename; 
    printf("target file path is %s \n", TargetFilePath.c_str());

    this->m_fp = fopen(TargetFilePath.c_str(), "wb"); 		    // 在此统一打开本地目标文件
    if(this->m_fp == NULL)
    {
        printf("target file %s open error\n", TargetFilePath.c_str());
        return;
    }
	
    if(lFileLen <= RANGE_SIZE)
    {
        start = 0;
        stop = lFileLen - 1;
        char szTaskName[64] = {0};
        if (m_eFileLoadType == FileLoad_Http_Curl)
        {
            snprintf(szTaskName,sizeof(szTaskName),"Curl Http down load, index-1[%ld-%ld]", start, stop);
            pTask = new CurlHttpTask(this->m_fp, start, stop, this->GetUrl(), false);
            pTask->SetTaskName(string(szTaskName));
            this->PushData(pTask);
        }
   	
    }else           // 创建分段下载子任务
    {
        seg_num = (long)lFileLen / RANGE_SIZE;
        for(int i = 0; i <= seg_num; i++)
        {
            if(i < seg_num)
            {
                start = i * RANGE_SIZE;
                stop = (i + 1) * RANGE_SIZE - 1;
            }
            else
            {
                if(lFileLen % RANGE_SIZE != 0){
                    start = i * RANGE_SIZE;
                    stop = lFileLen - 1;
                }else
                    break;
            }
            char szTaskName[64] = {0};
            if (m_eFileLoadType == FileLoad_Http_Curl)
            {
                snprintf(szTaskName,sizeof(szTaskName), "Curl Http down load, index-%d[%ld-%ld]", i+1,start, stop);
                TaskBase *pTask = new CurlHttpTask(this->m_fp, start, stop, this->GetUrl(), true);
                pTask->SetTaskName(string(szTaskName));
                this->PushData(pTask);
            }

        }
    }
}

//启动线程池
void DownLoadFactory::init(int defaultNum)
{
    this->m_pool.SetThreadNum(defaultNum);
    StartParser();
}

void DownLoadFactory::PushData(TaskBase* task)
{
    m_pool.AddTask(task);
}

void DownLoadFactory::StartParser()
{
    m_pool.Start();
}

void DownLoadFactory::StopParser()
{
    m_pool.Stop();
}
void DownLoadFactory::JoinParser()
{
    m_pool.Join();
}
