#include "CurlHttpTask.h"
#include <stdio.h>

pthread_mutex_t CurlHttpTask::m_mutex = PTHREAD_MUTEX_INITIALIZER;


TaskBase::~TaskBase()
{
}

size_t write_data(void* ptr, size_t size, size_t nmemb, void* userdata){
    CurlHttpTask* pTask = (CurlHttpTask*) userdata;
    bool ranged = pTask->m_ranged;
    size_t written;
    if(ranged)		            //分段下载
    {	
        //多线程写同一个文件, 需要加锁
        pthread_mutex_lock (&CurlHttpTask::m_mutex);
        if(pTask->m_startPos + size * nmemb <= pTask->m_stopPos)
        {
            fseek(pTask->m_fp, pTask->m_startPos, SEEK_SET);
            written = fwrite(ptr, size, nmemb, pTask->m_fp);
            pTask->m_startPos += size * nmemb;
        }
        else
        {
            fseek(pTask->m_fp, pTask->m_startPos, SEEK_SET);
            written = fwrite(ptr, 1, pTask->m_stopPos - pTask->m_startPos + 1, pTask->m_fp);
            pTask->m_startPos = pTask->m_stopPos;
        }
        pthread_mutex_unlock (&CurlHttpTask::m_mutex);
    }else
    {
        written = fwrite(ptr, size, nmemb, pTask->m_fp);
    }
    return written;
}

CurlHttpTask::~CurlHttpTask()
{
}


long CurlHttpTask::GetDownloadFileLength(const char *url)
{
    double dFile_Len = 0;
    CURL *handle = curl_easy_init ();

    curl_easy_setopt (handle, CURLOPT_URL, url);
    curl_easy_setopt (handle, CURLOPT_FOLLOWLOCATION, 1L);
    curl_easy_setopt (handle, CURLOPT_MAXREDIRS, 3L);
    curl_easy_setopt (handle, CURLOPT_AUTOREFERER, 1L);
    curl_easy_setopt (handle, CURLOPT_HEADER, 0L);  
    curl_easy_setopt (handle, CURLOPT_NOBODY, 1L);      //只获取响应头部信息
    curl_easy_setopt (handle, CURLOPT_FORBID_REUSE, 1);
	
    if (CURLE_OK == curl_easy_perform (handle) )
    {
        curl_easy_getinfo (handle, CURLINFO_CONTENT_LENGTH_DOWNLOAD, &dFile_Len);
    }
    else
    {
        dFile_Len = -1;
    }

    curl_easy_cleanup(handle);

    return dFile_Len;
}

int CurlHttpTask::Run(){
    CURL* curl;
    CURLcode res;
    int nResult = 0;

    char szRange[64] = { 0 };
    if(this->m_ranged) 
    {
        snprintf (szRange, sizeof(szRange), "%ld-%ld", this->m_startPos, this->m_stopPos);
        printf("CurlHttpTask range: [%s]\n", szRange);
    }
	
    curl = curl_easy_init();
    curl_easy_setopt(curl, CURLOPT_URL, this->m_url.c_str());
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_data);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void*)this);
    curl_easy_setopt(curl, CURLOPT_FOLLOWLOCATION, 1L);     //是否跟随重定向
    curl_easy_setopt(curl, CURLOPT_MAXREDIRS, 3L);          //HTTP 请求的最大重定向次数
    curl_easy_setopt(curl, CURLOPT_AUTOREFERER, 1L);
    curl_easy_setopt(curl, CURLOPT_NOSIGNAL, 1L);
    curl_easy_setopt(curl, CURLOPT_LOW_SPEED_LIMIT, 1L);
    curl_easy_setopt(curl, CURLOPT_LOW_SPEED_TIME, 30L);    //控制30秒传送CURLOPT_LOW_SPEED_LIMIT规定的字节数。
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, 300L); 	        //大文件容易超时，所以这里设置超时时间5分钟
    curl_easy_setopt(curl, CURLOPT_CONNECTTIMEOUT, 30L);    //连接超时时间设置30秒
    curl_easy_setopt(curl, CURLOPT_FORBID_REUSE, 1); 
    
    if(this->m_ranged)
    {
        curl_easy_setopt(curl, CURLOPT_RANGE, szRange);
    }
    res = curl_easy_perform(curl);
    if(CURLE_OK != res)
    {
        //写文件时当最后读取到的一段的大小超过endpos会只写部分，所以会有写错误产生不用处理
        if(res != 23){
            printf("Error Reason : %s, errorCode = %d\n", curl_easy_strerror(res),res);
            nResult =  -1;
        }
    }
    curl_easy_cleanup(curl);
    
    if (nResult < 0)
    {
        printf("down load error: %s\n",this->GetTaskName().c_str());
        return nResult;
    }
    
    printf("down load succeed: %s\n",this->GetTaskName().c_str());
    return 0;
}




