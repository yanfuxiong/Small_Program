#include "accept_file_dialog.h"
#include "ui_accept_file_dialog.h"
#include "common_signals.h"

AcceptFileDialog::AcceptFileDialog(QWidget *parent) :
    QDialog(parent),
    ui(new Ui::AcceptFileDialog)
{
    ui->setupUi(this);
    ui->icon_label->clear();
}

AcceptFileDialog::~AcceptFileDialog()
{
    delete ui;
}

void AcceptFileDialog::setFileInfo(const QString &fileInfo)
{
    ui->info_label->setText(fileInfo);
}

void AcceptFileDialog::on_accept_btn_clicked()
{
    Q_EMIT CommonSignals::getInstance()->userAcceptFile(true);
    accept();
}

void AcceptFileDialog::on_reject_btn_clicked()
{
    Q_EMIT CommonSignals::getInstance()->userAcceptFile(false);
    reject();
}
