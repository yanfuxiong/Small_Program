/*
* ==================================================================================================================================================
*
*       文件名:         DownloadFactory.h
*
*    	描述:      	下载器，抽象为一个类工厂可以通过不同的协议创建下载任务，方便其他协议扩展，目前只实现了curl http下载的方式
*
*       版本:    	    V1.0
*       创建时间:  	    2024-06-18
*       作者:  	    xiongyanfu
*       修改日志:  		1. Create file.
*
* ==================================================================================================================================================
*/

#ifndef __DOWN_LOAD_FACT_H
#define __DOWN_LOAD_FACT_H

#include "ThreadPool.h"
#include "CurlHttpTask.h"
//#include <vector>


//下载方式枚举
enum eFileLoadType
{
    FileLoad_Http_Curl    = 0,
    FileLoad_Ftp	
};


//工厂抽象类
class Factory
{
public:
    Factory(){}
    virtual ~Factory() = 0;

    virtual TaskBase *CreateTaskInstance() = 0;

    virtual void MakeTask(const string &filename) = 0;
    
};

class DownLoadFactory: public Factory
{
public:
    DownLoadFactory(){}
    
    
    /*
    *函数功能：下载器的构造函数
    *
    *输入参数：下载的url链接， 下载方式：默认 curl http方式
    *
    *返回值：void
    */
    DownLoadFactory(const string &url, const eFileLoadType type= FileLoad_Http_Curl);

    ~DownLoadFactory();

    /*
    *函数功能：工厂模式创建具体下载方式实例对象
    *
    *输入参数: void
    *
    *返回值：void
    */
    TaskBase *CreateTaskInstance();


    /*
    *函数功能：创建下载子任务
    *
    *输入参数: 下载的源文件存到本地的文件名称
    *
    *返回值：void
    */
    void MakeTask(const string &filename);


    /*
    *函数功能：子任务的任务队列进行初始化,并启动线程池处理业务
    *
    *输入参数: 线程数，默认为4
    *
    *返回值：void
    */
    void init(int defaultNum = 4);


    /*
    *函数功能：结束任务解析器
    *
    *输入参数: void
    *
    *返回值：void
    */
    void StopParser();

    /*
    *函数功能：等待任务解析器执行结束
    *
    *输入参数: void
    *
    *返回值：void
    */
    void JoinParser();

    
    /*
    *函数功能：获取当前下载的url链接
    *
    *输入参数: void
    *
    *返回值：返回url链接
    */
    inline string GetUrl();
private:
    /*
    *函数功能：开启任务解析器
    *
    *输入参数: void
    *
    *返回值：void
    */
    void StartParser();

    /*
    *函数功能：将下载任务添加到任务队列
    *
    *输入参数: void
    *
    *返回值：void
    */
    void PushData(TaskBase* task);

private:

    eFileLoadType   m_eFileLoadType;		// 下载类型
    string 	        m_url;
    ThreadPool 	    m_pool;
    FILE*           m_fp;
};

inline string DownLoadFactory::GetUrl()
{
    return this->m_url;
}

#endif