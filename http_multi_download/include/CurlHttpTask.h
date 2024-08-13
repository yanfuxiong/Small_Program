/*
* =============================================================================
*       文件名:  	    CurlHttpTask.h
*    	描述:    	    实现Curl Http下载任务， 继承于TaskBase
*       版本:    	    V1.0
*       创建时间:  	    2024-06-18
*       作者:  	    xiongyanfu
*       修改日志:   	1. Create file.
*	
* =============================================================================
*/

#ifndef __CURL_HTTP_TASK_H
#define __CURL_HTTP_TASK_H


#include "ThreadPool.h"
#include <curl/curl.h>
#include <string>


class CurlHttpTask :public TaskBase
{
public:	
    CurlHttpTask(){}
    ~CurlHttpTask();

    /*
    *函数功能：Curl Http任务类的构造函数
    *
    *输入参数：文件指针，写文件开始的位置，结束的位置，下载的url，是否分段
    *
    *返回值：void
    */
    CurlHttpTask(FILE* fp, long startPos, long stopPos, string url,bool ranged):m_fp(fp), m_startPos(startPos), m_stopPos(stopPos), m_url(url),m_ranged(ranged){}
	
    /*
    *函数功能：实现分段下载功能的run方法
    *
    *输入参数：void
    *
    *返回值：int
    */
    int Run();


    /*
    *函数功能：获取源文件长度
    *
    *输入参数：下载的url链接
    *
    *返回值：long
    */
    long GetDownloadFileLength(const char *url);
    /*
    *函数功能：设置线程任务的名称
    *
    *输入参数：自定义线程任务名称
    *
    *返回值：void
    */
    void SetTaskName(const string &name)
    {
        this->m_sTaskName = name;
    }

    string GetTaskName()
    {
        return this->m_sTaskName;
    }
public:
    	
    static pthread_mutex_t m_mutex;
    
    FILE*    m_fp;  			//本地文件句柄
    long     m_startPos;
    long     m_stopPos;
    string   m_url;
    bool     m_ranged; 			//是否需要http 分段下载?
};

#endif