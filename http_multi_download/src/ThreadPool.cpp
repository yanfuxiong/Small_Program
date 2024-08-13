#include "ThreadPool.h"
#include <iostream>
using namespace std;


ThreadPool::~ThreadPool()
{
    m_bRunning      = false;

    if (NULL != m_threadBuf)
    {
        delete[]        m_threadBuf;
        m_threadBuf     = NULL;
    }
	
    pthread_mutex_destroy(&m_mutex);
    pthread_cond_destroy(&m_con);

    while(!m_deque.empty())
    {
        TaskBase * pTask = m_deque.pop_front();
        delete pTask;
    }

   	m_deque.set_stop();
}

ThreadPool::ThreadPool(int defaultNum)
{
    m_nThreadNum 	= defaultNum;  
    m_bRunning 		= false;
    m_nThreadMaxNum = MAX_THREAD_NUM;
    m_threadBuf     = NULL;
	
}

void ThreadPool::Start()
{
    bool bRun = true;

    pthread_mutex_lock (&m_mutex);
    bRun = m_bRunning;
    if(!m_bRunning)
        m_bRunning = true;
    pthread_mutex_unlock (&m_mutex);

    if(!bRun)
    {
        pthread_cond_init(&m_con, NULL);
        pthread_mutex_init(&m_mutex, NULL);
        cout << "ThreadPool:: Start" << endl;

        auto cpuNum = static_cast<int>(sysconf(_SC_NPROCESSORS_ONLN)); 
        cpu_set_t cpu_info;

        m_deque.set_start();
				
        if (m_nThreadNum > this->GetTaskNum())			//线程数不能大于任务数
            this->SetThreadNum(this->GetTaskNum());

        if (m_nThreadNum > m_nThreadMaxNum)
            this->SetThreadNum(m_nThreadMaxNum);			

        std::cout<<" m_nThreadNum = "<<m_nThreadNum<<std::endl;

        m_threadBuf = new pthread_t[m_nThreadNum];
        for(auto i = 0, j=0; i < m_nThreadNum; i++,j++)
        {
            pthread_create(m_threadBuf + i, NULL, ThreadFunc, this);
            if(cpuNum > m_nThreadNum){

                pthread_join(m_threadBuf[i], NULL);
                CPU_ZERO(&cpu_info);
                CPU_SET(j,&cpu_info);
                pthread_setaffinity_np(m_threadBuf[i],sizeof(cpu_set_t),&cpu_info);
            }

        }

    }
}

void ThreadPool::Stop()
{
    bool bRun;

    pthread_mutex_lock (&m_mutex);
    bRun = m_bRunning;
    if(m_bRunning)
        m_bRunning = false;
    pthread_mutex_unlock (&m_mutex);

    if(bRun)
    {	
        std::cout<<"ThreadPool:: Stop"<<std::endl;
        
        /*for(int i = 0; i<m_nThreadNum; i++)
        {
            pthread_join(m_threadBuf[i], NULL);
        }*/
       
        while(!m_deque.empty())
        {
            TaskBase * pTask = m_deque.pop_front();
            delete  pTask;
        }
        m_deque.set_stop();
    }
}

void ThreadPool::Join()
{
    if (NULL == m_threadBuf)
        return;
    
	for(int i = 0; i<m_nThreadNum; i++)
    {
        pthread_join(m_threadBuf[i], NULL);
    }
}

void ThreadPool::SetThreadNum(const int num)
{
    this->m_nThreadNum = num;
}

int ThreadPool::GetThreadNum()
{
    return m_nThreadNum;
}

int ThreadPool::GetTaskNum()
{
    return m_deque.size();
}

TaskBase* ThreadPool::Take()
{
    return m_deque.pop_front();
}

void ThreadPool::AddTask(TaskBase *task)
{
    m_deque.push_back(task);
    pthread_cond_broadcast(&m_con);
}

void* ThreadPool::ThreadFunc(void * arg)
{
    std::cout<< "ThreadPool: ThreadFunc start\n" <<std::endl;

    int nResult = 0;
    bool bRun = false;

    auto pool = static_cast<ThreadPool *>(arg);
    
    while(true)
    {
        pthread_mutex_lock (&pool->m_mutex);
        bRun = pool->m_bRunning;
        pthread_mutex_unlock (&pool->m_mutex);
        if (!bRun)
            break;

        TaskBase*  pTask = pool->Take();
        if(NULL == pTask)
            break;

        try
        {
            nResult = pTask->Run();
        }
        catch(const std::exception& e)
        {
            std::cout<<"ThreadPool: func threw "<< typeid(e).name() << " exception: " << e.what()<<std::endl;
            nResult = -99;
        }
        catch(...)
        {
            std::cout<<"ThreadPool: func threw unhandled non-exception object"<<std::endl;
            nResult = -98;
        }

        delete pTask;                   //执行完的任务对象,需要释放

        if (nResult < 0)                //  某个分段任务失败则终止所有下载
            pool->Stop();

    }

    std::cout<< "ThreadPool: ThreadFunc end\n" <<std::endl;

    return (void *)NULL;
}
