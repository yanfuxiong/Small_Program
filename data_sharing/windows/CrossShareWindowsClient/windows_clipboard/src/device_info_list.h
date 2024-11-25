#ifndef DEVICE_INFO_LIST_H
#define DEVICE_INFO_LIST_H

#include <QWidget>
#include "common_utils.h"
#include "common_signals.h"
#include "device_info.h"
#include <QVBoxLayout>
#include <QHBoxLayout>
#include <QScrollArea>

namespace Ui {
class DeviceInfoList;
}

class DeviceInfoList : public QWidget
{
    Q_OBJECT

public:
    explicit DeviceInfoList(QWidget *parent = nullptr);
    ~DeviceInfoList();

private Q_SLOTS:
    void onUpdateClientList();

private:
    Ui::DeviceInfoList *ui;
};

#endif // DEVICE_INFO_LIST_H
