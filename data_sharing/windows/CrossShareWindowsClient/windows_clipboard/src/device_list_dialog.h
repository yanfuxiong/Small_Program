#ifndef DEVICE_LIST_DIALOG_H
#define DEVICE_LIST_DIALOG_H

#include <QDialog>
#include "flowlayout.h"
#include <QVBoxLayout>

namespace Ui {
class DeviceListDialog;
}

class DeviceListDialog : public QDialog
{
    Q_OBJECT

public:
    explicit DeviceListDialog(QWidget *parent = nullptr);
    ~DeviceListDialog();

private Q_SLOTS:
    void onSendDataToServer(const QByteArray &data);
    void onUpdateClientList();

private:
    Ui::DeviceListDialog *ui;
    //FlowLayout *m_flowLayout;
};

#endif // DEVICE_LIST_DIALOG_H
