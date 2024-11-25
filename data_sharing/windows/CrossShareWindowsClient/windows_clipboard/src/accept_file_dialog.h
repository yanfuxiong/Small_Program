#ifndef ACCEPT_FILE_DIALOG_H
#define ACCEPT_FILE_DIALOG_H

#include <QDialog>

namespace Ui {
class AcceptFileDialog;
}

class AcceptFileDialog : public QDialog
{
    Q_OBJECT

public:
    explicit AcceptFileDialog(QWidget *parent = nullptr);
    ~AcceptFileDialog();

    void setFileInfo(const QString &fileInfo);

private slots:
    void on_accept_btn_clicked();
    void on_reject_btn_clicked();

private:
    Ui::AcceptFileDialog *ui;
};

#endif // ACCEPT_FILE_DIALOG_H
