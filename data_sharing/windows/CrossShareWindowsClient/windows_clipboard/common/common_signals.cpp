#include "common_signals.h"

CommonSignals* CommonSignals::m_instance = nullptr;


CommonSignals::CommonSignals()
{

}

CommonSignals *CommonSignals::getInstance()
{
    if (m_instance == nullptr) {
        m_instance = new CommonSignals;
    }
    return m_instance;
}
