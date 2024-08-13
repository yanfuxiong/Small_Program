/*
 * =============================================================================
 *
 *       文件名:       ThreadPool.h
 * 
 *    	 描述:        线程池的基础类，可以创建指定数量的线程，批量处理队列任务
 *
 *       版本:        V1.0
 *       创建时间:  	2024-06-18
 *       作者:       	xiongyanfu
 *       修改日志:  	1. Create file.
 *
 * =============================================================================
 */


#ifndef __THREAD_BASE_H
#define __THREAD_BASE_H

#include <pthread.h>
#include <string>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <sched.h>
#include <iostream>


#include "kdeque.hpp"
#include "TaskBase.h"


#define DEFAUL_THREAD_NUM 	4           //默认线程数
#define MAX_THREAD_NUM 		1000



class ThreadPool
{
public:
    /*
    *函数功能：线程池的构造函数
    *
    *输入参数：线程池创建的线程数量
    *
    *返回值：void
    */
    ThreadPool(int defaultNum = DEFAUL_THREAD_NUM);



    /*
    *函数功能：禁用线程池的拷贝构造函数
    */
    ThreadPool(const ThreadPool&) = delete;	



    /*
    *函数功能：禁用线程池的拷贝赋值函数
    */
    ThreadPool& operator=(const ThreadPool) = delete;



    /*
    *函数功能：线程池的析构函数
    *
    *输入参数：void
    *
    *返回值：void
    */
    ~ThreadPool();



    /*
    *函数功能：线程池启动
    *
    *输入参数：void
    *
    *返回值：void
    */
    void Start();



    /*
    *函数功能：线程池停止
    *
    *输入参数：void
    *
    *返回值：void
    */
    void Stop();


    /*
    *函数功能：等待线程池处理完毕
    *
    *输入参数：void
    *
    *返回值：void
    */
    void Join();


    /*
    *函数功能：向任务队列添加任务
    *
    *输入参数：void
    *
    *返回值：void
    */
    void AddTask(TaskBase *task);

    
    /*
    *函数功能：获取线程池的当前线程数量
    *
    *输入参数：void
    *
    *返回值：线程数量
    */
    int GetThreadNum();



    /*
    *函数功能：获取线程池的当前的任务数量
    *
    *输入参数：void
    *
    *返回值：任务量
    */
    int GetTaskNum();

    /*
    *函数功能：设置线程数量
    *
    *输入参数：线程数量
    *
    *返回值：void
    */
    void SetThreadNum(const int num);
private:
    /*
    *函数功能：向任务队列取任务
    *
    *输入参数：void
    *
    *返回值：任务指针
    */
    TaskBase* Take();


    /*
    *函数功能：线程函数
    *
    *输入参数：void*
    *
    *返回值：void
    */
    static void* ThreadFunc(void * arg);

    pthread_cond_t      m_con;

    pthread_mutex_t     m_mutex;
    
private:
    bool                m_bRunning;
    KDeque<TaskBase*>   m_deque;
    int                 m_nThreadNum;
    int                 m_nThreadMaxNum;
	
    pthread_t*          m_threadBuf;
	
};


#endif