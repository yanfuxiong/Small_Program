/*
* =============================================================================
*       文件名:  	    TaskBase.h
*    	描述:    	    下载任务虚基类，所有的协议下载任务类由此派生
*       版本:    	    V1.0
*       创建时间:       2024-06-18
*       作者:         xiongyanfu
*       修改日志:  		1. Create file.
*
* =============================================================================
*/

#ifndef __TASK_BASE_H
#define __TASK_BASE_H

#include <string>


/*
*
*类说明：单个任务的抽象类
*/
class TaskBase
{
public:
    /*
    *函数功能：抽象类的构造函数
    *
    *输入参数：void
    *
    *返回值：void
    */

    TaskBase(){}

    virtual ~TaskBase() = 0;


    /*
    *函数功能：设置线程的名称
    *
    *输入参数：自定义线程名称
    *
    *返回值：void
    */
    virtual void SetTaskName(const string &name) = 0;

    /*
    *函数功能：获取文件长度
    *
    *输入参数：下载的url链接
    *
    *返回值：long
    */
    virtual long GetDownloadFileLength(const char *url) = 0;

    /*
    *函数功能：虚函数，执行下载动作
    *
    *输入参数：void
    *
    *返回值：int
    */
    virtual int Run() = 0;


public:
    //void* 	m_arg;
    string 		m_sTaskName;		//任务明细
};


#endif

