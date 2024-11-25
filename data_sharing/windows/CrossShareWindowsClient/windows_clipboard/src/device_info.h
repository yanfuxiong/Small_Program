#ifndef DEVICE_INFO_H
#define DEVICE_INFO_H

#include <QWidget>
#include <QLabel>
#include <QPushButton>
#include <QPointer>
#include <QHBoxLayout>
#include "common_utils.h"

namespace Ui {
class DeviceInfo;
}

class DeviceInfo : public QWidget
{
    Q_OBJECT

public:
    explicit DeviceInfo(QWidget *parent = nullptr);
    ~DeviceInfo();

    void setDeviceInfo(const nlohmann::json &deviceInfoJson);
    QByteArray getClientID() const;

private slots:
    void onOpenFile();
    void onUpdateProgressInfoWithID(int currentVal, const QByteArray &clientID);
    void onUpdateControlStatus(bool status);

private:
    Ui::DeviceInfo *ui;
    nlohmann::json m_deviceInfoJson;

    QLabel *m_imageLabel;
    QLabel *m_deviceInfoLabel;
    QPushButton *m_fileIcon;
};

#endif // DEVICE_INFO_H
