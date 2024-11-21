package com.rtk.myapplication;

import android.annotation.SuppressLint;
import android.app.Activity;
import android.app.AlertDialog;
import android.app.AppOpsManager;
import android.app.Service;
import android.content.ClipData;
import android.content.ClipboardManager;
import android.content.ContentResolver;
import android.content.ContentValues;
import android.content.Context;
import android.content.DialogInterface;
import android.content.Intent;
import android.content.pm.PackageManager;
import android.database.Cursor;
import android.graphics.BitmapFactory;
import android.media.Image;
import android.media.MediaPlayer;
import android.net.ConnectivityManager;
import android.net.LinkAddress;
import android.net.LinkProperties;
import android.net.Uri;
import android.net.wifi.WifiInfo;
import android.net.wifi.WifiManager;
import android.os.Binder;
import android.os.Build;
import android.os.Bundle;
import android.os.Environment;
import android.os.Handler;
import android.os.ParcelFileDescriptor;
import android.provider.ContactsContract;
import android.provider.MediaStore;
import android.provider.OpenableColumns;
import android.provider.Settings;
import android.text.TextUtils;
import android.text.format.Formatter;
import android.util.Base64;
import android.util.Log;
import android.view.Gravity;
import android.view.LayoutInflater;
import android.view.MotionEvent;
import android.view.View;
import android.view.Window;
import android.view.WindowManager;
import android.widget.Button;
import android.widget.CheckBox;
import android.widget.CompoundButton;
import android.widget.EditText;
import android.widget.ImageView;
import android.widget.TextView;
import android.widget.Toast;

import androidx.annotation.NonNull;
import androidx.annotation.Nullable;
import androidx.core.app.ActivityCompat;
import androidx.core.content.ContextCompat;
import androidx.recyclerview.widget.LinearLayoutManager;
import androidx.recyclerview.widget.RecyclerView;

import java.io.BufferedReader;
import java.io.ByteArrayOutputStream;
import java.io.File;
import java.io.FileDescriptor;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.io.UnsupportedEncodingException;
import java.lang.reflect.Field;
import java.lang.reflect.Method;
import java.net.NetworkInterface;
import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;
import java.text.DecimalFormat;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.Enumeration;
import java.util.List;
import java.util.concurrent.atomic.AtomicReference;

import libp2p_clipboard.Callback;
import libp2p_clipboard.Libp2p_clipboard;

import android.content.ClipData;
import android.graphics.Bitmap;
import android.widget.VideoView;

import com.tencent.mmkv.MMKV;

public class TestActivity extends Activity {

    private static final int REQUEST_CODE = 1024;
    private static final String TAG = "TestActivity";
    private ClipboardManager clipboardManager;
    private int testCount = 0;
    private TextView mServerStatus;
    private TextView mClientStatus;
    private TextView mPeerMessage;
    private EditText mServerId;
    private EditText mServerIpInfo;
    private boolean mIsConnected = false;
    private String mLastString = "";
    private ImageView imageView;

    private Button btnGetImage;
    private Button btnSetImage;
    private ImageView imageView2, imageView3;
    private TextView textView;
    Bitmap bitmap;
    VideoView videview;
    Button mbutton, mbuttonpaste, buttom_w, buttom_r;
    String value;
    String mimetype;
    String sharedText;
    byte[] getbyteArray;
    String sizeInMB;
    Intent intent;
    String action;
    String base64String, clearbase64String;
    Context mContext;
    Bitmap bitmapShare;
    byte[] imageData;
    String text;
    Uri uri;
    ContentResolver resolver;
    boolean check = false;
    MMKV kv;
    boolean boxischeck;
    TextView textView_name, textView_size;

    @Override
    protected void onCreate(@Nullable Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.layout_testactivity);
        mContext = this;
        mServerStatus = (TextView) findViewById(R.id.status_server);
        mClientStatus = (TextView) findViewById(R.id.status_client);
        mPeerMessage = (TextView) findViewById(R.id.peer_message);
        mServerId = (EditText) findViewById(R.id.server_id);
        mServerIpInfo = (EditText) findViewById(R.id.server_ip_info);
        imageView = (ImageView) findViewById(R.id.imageView);
        clipboardManager = (ClipboardManager) getSystemService(Context.CLIPBOARD_SERVICE);
        clipboardManager.addPrimaryClipChangedListener(() -> {
            Log.i(TAG, "onCreate: addPrimaryClipChangedListener mIsConnected" + mIsConnected);
            if (!mIsConnected) {
                return;
            }
            Log.i(TAG, "onCreate: p1");
            CharSequence mText = clipboardManager.getText();
            if (mText != null && !TextUtils.isEmpty(mText.toString()) && !mText.toString().equals(mLastString)) {
                Log.i(TAG, "onCreate: p2");
                mLastString = mText.toString();
                Log.i(TAG, "onCreate mLastString = " + mLastString);
                Libp2p_clipboard.sendMessage(mLastString);
            }
        });

        kv = MMKV.defaultMMKV();

        //alertDialog("aa",12L);

        boxischeck = kv.decodeBool("ischeck", false);
        Log.i("lszz", "onCreate: CheckBox boxischeck===" + boxischeck);

        //lsz modify
        btnGetImage = findViewById(R.id.btnGetImage);
        btnSetImage = findViewById(R.id.btnSetImage);
        imageView2 = findViewById(R.id.imageView2);
        imageView3 = findViewById(R.id.imageView3);
        textView = findViewById(R.id.textView);
        bitmap = BitmapFactory.decodeResource(getResources(), R.drawable.bb);
        videview = findViewById(R.id.videoview);
        //imageView2.setImageBitmap(bitmap);
        ClipboardUtils.setContext(getApplication());
        ClipboardUtils clipboardUtils = ClipboardUtils.getInstance();

        btnGetImage.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Log.d("lsz", "btnGetImage hasClip+==" + clipboardUtils.hasClip());
                if (clipboardUtils.hasClip()) {
                    getClipFromClipboard();
                } else {
                    Toast.makeText(TestActivity.this, "lsz get Clipboard is empty", Toast.LENGTH_SHORT).show();
                }

            }
        });

        // 设置点击监听器，将图片设置到剪贴板
        btnSetImage.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                setClipToClipboard();
            }
        });

        String wifiip = getWifiIpAddress(this);
        Log.d("lszz", "get wifiip==" + wifiip);

        new Thread(new Runnable() {
            @Override
            public void run() {
                Libp2p_clipboard.sendAddrsFromJava(getIpAddres());
                //mdns直连
                Libp2p_clipboard.mainInit(getGolangCallBack(), wifiip, "6666", 6633);
            }
        }).start();


        //lsz add for test CrossShare
        try {
            Log.i(TAG, "lszz = activity oncreat");
            testCrossShare();
            try {
                getShare(intent, action, mimetype);
            } catch (FileNotFoundException e) {
                throw new RuntimeException(e);
            } catch (Exception e) {
                throw new RuntimeException(e);
            }
        } catch (IOException e) {
            throw new RuntimeException(e);
        }

        mbutton = findViewById(R.id.buttom);
        mbuttonpaste = findViewById(R.id.buttom_paste);
        buttom_w = findViewById(R.id.buttom_w);
        buttom_r = findViewById(R.id.buttom_r);
        textView_name = findViewById(R.id.textView_name);
        textView_size = findViewById(R.id.textView_size);

        buttom_w.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {
                Bitmap bitmap = BitmapFactory.decodeResource(getResources(), R.drawable.aa);
                ByteArrayOutputStream stream = new ByteArrayOutputStream();
                bitmap.compress(Bitmap.CompressFormat.PNG, 100, stream);
                byte[] byteArray = stream.toByteArray();

                ContentValues values = new ContentValues();
                values.put("image", byteArray);
                Uri newUri = getContentResolver().insert(UserContentProvider.CONTENT_URI, values);
                int rowsUpdated = getContentResolver().update(UserContentProvider.CONTENT_URI, values, null, null);

                //发送默认路径的文件
                String path = "/storage/emulated/0/Android/data/com.rtk.myapplication/files/log.txt";
                String cliendid = "192.168.22.211:1518";
                File file = new File(path);
                if (file.exists()) {
                    int fileSize = (int)file.length();
                    Libp2p_clipboard.sendCopyFile(path, cliendid, 0, fileSize);
                }


            }

        });

        buttom_r.setOnClickListener(new View.OnClickListener() {
            @SuppressLint("Range")
            @Override
            public void onClick(View view) {
                Log.d("lszz", "点击复制按钮监听，发送到libp2p  cursor+==aa");
                try (Cursor cursor = getContentResolver().query(UserContentProvider.CONTENT_URI, null, null, null, null)) {
                    Log.d("lszz", "点击复制按钮监听，发送到libp2p  cursor+==" + cursor);
                    if (cursor != null && cursor.moveToFirst()) {
                        @SuppressLint("Range") byte[] blob = cursor.getBlob(cursor.getColumnIndex("image"));

                        Bitmap bitmap = BitmapFactory.decodeByteArray(blob, 0, blob.length);
                        imageView3.setImageBitmap(bitmap);
                        imageView3.setVisibility(View.VISIBLE);
                    }
                } catch (Exception e) {
                    e.printStackTrace();
                }

            }

        });

        //点击粘贴按钮监听
        mbuttonpaste.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {
                //todo
                //Libp2p_clipboard.clipboardImagePaste("111");
            }

        });
        //点击复制按钮监听，发送到libp2p
        mbutton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {
                //Log.d("lszz","点击复制按钮监听，发送到libp2p  hasClip+=="+clipboardUtils.hasClip());
                Log.d("lszz", "点击复制按钮监听，发送到libp2p  text+==" + text);
                Log.d("lszz", "点击复制按钮监听，发送到libp2p  imageData+==" + imageData);
                Toast.makeText(TestActivity.this, "is empty!!", Toast.LENGTH_SHORT).show();
                if (clipboardUtils.hasClip()) {
                    //getClipFromClipboard();
                    if (text != null) {
                        sendToPC(text);
                    }
                    if (imageData != null) {
                        sendToPCIMG(imageData);

                        ContentValues values = new ContentValues();
                        values.put("image", imageData);
                        Uri newUri = getContentResolver().insert(UserContentProvider.CONTENT_URI, values);
                        getContentResolver().update(UserContentProvider.CONTENT_URI, values, null, null);

                    }

                } else {
                    Toast.makeText(TestActivity.this, "lsz get Clipboard is empty", Toast.LENGTH_SHORT).show();
                }


            }
        });


        //

       /* imageView.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {
                Log.d("lsz","lsz==aaaaaaaa");
                ClipboardManager clipboard = (ClipboardManager) getSystemService(Context.CLIPBOARD_SERVICE);
                // 檢查剪貼簿是否有內容
                if (clipboard.hasPrimaryClip()) {
                    ClipData clip = clipboard.getPrimaryClip();

                    // 檢查是否包含 URI 類型的資料
                    if (clip != null && clip.getItemCount() > 0) {
                        ClipData.Item item = clip.getItemAt(0);
                        Uri imageUri = item.getUri();
                        Log.d("lsz", "lsz==" + imageUri);
                        if (imageUri != null) {
                            // 將圖片 URI 設置到 ImageView 顯示圖片
                            ImageView imageView = (ImageView) findViewById(R.id.imageView);
                            imageView.setImageURI(imageUri);
                        }
                    }
                } else {
                    Log.d("lsz","lsz Clipboard is empty");
                    Toast.makeText(TestActivity.this, "lsz Clipboard is empty", Toast.LENGTH_SHORT).show();
                }

            }
        });*/
        findViewById(R.id.bt_connlibp2p).setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                new Thread(new Runnable() {
                    @Override
                    public void run() {
                        //Libp2p_clipboard.sendAddrsFromJava(getIpAddres());
                        Log.d("lsz", "GoLog getIpAddres()=" + getIpAddres());
                        Log.d("lsz", "GoLog mServerId.getText().toString()=" + mServerId.getText().toString());
                        Log.d("lsz", "GoLog mServerIpInfo.getText().toString()=" + mServerIpInfo.getText().toString());
                        //Libp2p_clipboard.mainInit(getGolangCallBack(), mServerId.getText().toString(), mServerIpInfo.getText().toString());

                    }
                }).start();
            }
        });

        findViewById(R.id.bt_requestfloatwindow).setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                if (!Settings.canDrawOverlays(TestActivity.this)) {
                    Intent intent = new Intent(Settings.ACTION_MANAGE_OVERLAY_PERMISSION,
                            Uri.parse("package:" + TestActivity.this.getPackageName()));
                    TestActivity.this.startActivity(intent);
                }
            }
        });
        findViewById(R.id.bt_openfloatwindow).setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Log.e("lsz", "checkFloatPermission=" + checkFloatPermission(mContext));
                if (checkFloatPermission(mContext) == true) {
                    Intent serviceIntent = new Intent(TestActivity.this, FloatClipboardService.class);
                    LibP2pUtils.updateIdAndIp(mServerId.getText().toString(), mServerIpInfo.getText().toString());
                    if (android.os.Build.VERSION.SDK_INT >= android.os.Build.VERSION_CODES.O) {
                        startForegroundService(serviceIntent);
                    } else {
                        startService(serviceIntent);
                    }
                } else {
                    Toast.makeText(TestActivity.this, "请先开启悬浮窗权限", Toast.LENGTH_SHORT).show();
                }
            }
        });
        findViewById(R.id.bt_listenclip).setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                clipboardManager.addPrimaryClipChangedListener(() -> {
                    ClipData clip = clipboardManager.getPrimaryClip();
                    if (clip != null && clip.getItemCount() > 0) {
                        ClipData.Item item = clip.getItemAt(0);
                        String text = item.getText().toString();
                        Log.i(TAG, "clipboardManager listen = " + text);
                        Libp2p_clipboard.sendMessage(text);
                    }
                });
            }
        });
        findViewById(R.id.bt_testcallback).setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                // Libp2p_clipboard.performTask("测试回调");
            }
        });
        Log.i(TAG, "lsz path = " + getExternalFilesDir(null));
        //requestPermission();
        //   getIpAddres();

    }

    private void getClibMessageLoop() {
        new Handler().postDelayed(new Runnable() {
            @Override
            public void run() {
                getClibMessage();
                getClibMessageLoop();
            }
        }, 3000);
    }

    private void setClibMessageLoop() {
        new Handler().postDelayed(new Runnable() {
            @Override
            public void run() {
                setClibMessage();
                setClibMessageLoop();
            }
        }, 3000);
    }

    private void getClibMessage() {
        if (clipboardManager.hasPrimaryClip()) {
            // 剪贴板有数据
            ClipData clipData = clipboardManager.getPrimaryClip();
            if (clipData != null && clipData.getItemCount() > 0) {
                CharSequence itemText = clipData.getItemAt(0).getText();
                // 使用 itemText 中的数据
                Log.i(TAG, "getClibMessage: itemText==" + itemText.toString());
            }
        }
    }


    private void setClibMessage() {
        ClipData clipData = ClipData.newPlainText(null, "编辑后的文本数据+" + testCount);
        clipboardManager.setPrimaryClip(clipData);
        testCount++;
    }

    @Override
    protected void onResume() {
        super.onResume();
        Log.i(TAG, "onResume: mIsConnected: " + mIsConnected);

        /*try {
            testCrossShare();
            getShare(intent, action,mimetype);
        } catch (FileNotFoundException e) {
            throw new RuntimeException(e);
        } catch (Exception e) {
            throw new RuntimeException(e);
        }*/

        /*if (mIsConnected) {
            Log.i(TAG, "onCreate: p1");
            CharSequence mText =  clipboardManager.getText();
            Log.i(TAG, "onCreate: mText: " + mText);
            if (mText != null && !TextUtils.isEmpty(mText.toString()) && !mText.toString().equals(mLastString)) {
                Log.i(TAG, "onCreate: p2");
                mLastString = mText.toString();
                Log.i(TAG, "onCreate mLastString = " + mLastString);
                Libp2p_clipboard.sendMessage(mLastString);
            }
        }*/
    }

    private void requestPermission() {
        if (ActivityCompat.checkSelfPermission(this, android.Manifest.permission.READ_EXTERNAL_STORAGE) == PackageManager.PERMISSION_GRANTED &&
                ContextCompat.checkSelfPermission(this, android.Manifest.permission.WRITE_EXTERNAL_STORAGE) == PackageManager.PERMISSION_GRANTED) {
            Toast.makeText(this, "存储权限获取成功", Toast.LENGTH_SHORT).show();
        } else {
            ActivityCompat.requestPermissions(this, new String[]{android.Manifest.permission.READ_EXTERNAL_STORAGE, android.Manifest.permission.WRITE_EXTERNAL_STORAGE}, REQUEST_CODE);
        }

    }

    @Override
    public void onRequestPermissionsResult(int requestCode, @NonNull String[] permissions, @NonNull int[] grantResults) {
        super.onRequestPermissionsResult(requestCode, permissions, grantResults);
        if (requestCode == REQUEST_CODE) {
            if (ActivityCompat.checkSelfPermission(this, android.Manifest.permission.READ_EXTERNAL_STORAGE) == PackageManager.PERMISSION_GRANTED &&
                    ContextCompat.checkSelfPermission(this, android.Manifest.permission.WRITE_EXTERNAL_STORAGE) == PackageManager.PERMISSION_GRANTED) {
                Toast.makeText(this, "存储权限获取成功", Toast.LENGTH_SHORT).show();
            } else {
                Toast.makeText(this, "存储权限获取失败", Toast.LENGTH_SHORT).show();
            }
        }
    }

    private String getIpAddres() {

        ConnectivityManager connectivityManager = (ConnectivityManager) getApplicationContext().getSystemService(Service.CONNECTIVITY_SERVICE);
        LinkProperties linkProperties = connectivityManager.getLinkProperties(connectivityManager.getActiveNetwork());
        List<LinkAddress> addressList = linkProperties.getLinkAddresses();
        StringBuffer sbf = new StringBuffer();
        for (LinkAddress linkAddress : addressList) {
            sbf.append(linkAddress.toString()).append("#");
            Log.d(TAG, "xyf getIpAddreslinkAddress.toString(): " + linkAddress.toString());
        }
        Log.d(TAG, "xyf getIpAddres: " + sbf.toString());
        return sbf.toString();
    }

    public static String getWifiIpAddress(Context context) {
        WifiManager wifiManager = (WifiManager) context.getSystemService(Context.WIFI_SERVICE);
        if (wifiManager != null && wifiManager.getConnectionInfo() != null) {
            int ipAddress = wifiManager.getConnectionInfo().getIpAddress();
            return Formatter.formatIpAddress(ipAddress);
        }
        return null;


    }

    private final static int P2P_EVENT_SERVER_CONNEDTED = 0;
    private final static int P2P_EVENT_SERVER_CONNECT_FAIL = 1;
    private final static int P2P_EVENT_CLIENT_CONNEDTED = 2;
    private final static int P2P_EVENT_CLIENT_CONNECT_FAIL = 3;

    private Callback getGolangCallBack() {
        return new Callback() {
            @Override
            public void callbackMethod(String s) {
                Log.i(TAG, "lsz GoLog callbackMethod: callbacl调用 ==abc=" + s);

                ClipData clipData = ClipData.newPlainText(null, s);
                clipboardManager.setPrimaryClip(clipData);
                new Thread(new Runnable() {
                    @Override
                    public void run() {
                        //do something takes long time in the work-thread
                        runOnUiThread(new Runnable() {
                            @Override
                            public void run() {
                                mPeerMessage.setText(s);
                            }
                        });
                    }
                }).start();


            }

            @Override
            public void callbackMethodFileConfirm(String s, long l) {
                Log.i(TAG, "lszz GoLog callbackMethodFileConfirm: amsg:String= " + s);
                Log.i(TAG, "lszz GoLog callbackMethodFileConfirm: amsg:long= " + l);
                Log.i(TAG, "lszz GoLog callbackMethodFileConfirm: amsg:boxischeck= " + boxischeck);
                runOnUiThread(new Runnable() {
                                  @Override
                                  public void run() {
                                      if (!boxischeck) {
                                          Log.i("lszz", "CheckBox boxischeck======false");
                                      } else {
                                          Log.i("lszz", "CheckBox boxischeck======true");
                                      }
                                      alertDialog(s, l);
                                  }
                              }
                );


            }

            @Override
            public void callbackMethodImage(String msg) {
                Log.i(TAG, "lszz GoLog callbackMethodImage: msg: " + msg);
                Log.i(TAG, "lszz GoLog callbackMethodImage: msg.length=: " + msg.length());

                runOnUiThread(new Runnable() {
                                  @Override
                                  public void run() {

                                      byte[] value = Base64.decode(msg, Base64.DEFAULT);
                                      Log.i("lszzz", "GoLog sendToPC: img value:==" + value);
                                      Bitmap bitmap = BitmapFactory.decodeByteArray(value, 0, value.length);
                                      imageView2.setImageBitmap(bitmap);

                                  }
                              }
                );

            }

            @Override
            public void logMessageCallback(String msg) {
                Log.i(TAG, "lsz GoLog logMessageCallback: msg: " + msg);
            }

            @Override
            public void eventCallback(long event) {
                Log.i(TAG, "lsz get GoLog eventCallBack: event: " + event);
                new Thread(new Runnable() {
                    @Override
                    public void run() {
                        //do something takes long time in the work-thread
                        runOnUiThread(new Runnable() {
                            @Override
                            public void run() {
                                switch ((int) event) {
                                    case P2P_EVENT_SERVER_CONNEDTED:
                                        Log.i(TAG, "eventCallBack: P2P_EVENT_SERVER_CONNEDTED");
                                        mServerStatus.setText("connected");
                                        mIsConnected = true;
                                        break;
                                    case P2P_EVENT_SERVER_CONNECT_FAIL:
                                        mServerStatus.setText("failed to connected");
                                        Log.i(TAG, "eventCallBack: P2P_EVENT_SERVER_CONNECT_FAIL");
                                        break;
                                    case P2P_EVENT_CLIENT_CONNEDTED:
                                        Log.i(TAG, "eventCallBack: P2P_EVENT_CLIENT_CONNEDTED");
                                        mClientStatus.setText("connected");
                                        mPeerMessage.setText("");
                                        break;
                                    case P2P_EVENT_CLIENT_CONNECT_FAIL:
                                        Log.i(TAG, "eventCallBack: P2P_EVENT_CLIENT_CONNECT_FAIL");
                                        mClientStatus.setText("failed to connected");
                                        break;
                                    default:
                                        break;
                                }
                            }
                        });
                    }
                }).start();
            }
        };
    }


    private void getClipFromClipboard() {
//本地图片 取到剪切版
       /* ClipboardManager clipboard = (ClipboardManager) getSystemService(Context.CLIPBOARD_SERVICE);
        Log.d("lsz","clipa"+clipboard.hasPrimaryClip());
        // 檢查剪貼簿是否有內容
        if (clipboard.hasPrimaryClip()) {
            ClipData clip = clipboard.getPrimaryClip();
            Log.d("lsz","clip"+clip);
            // 檢查是否包含 URI 類型的資料
            if (clip != null && clip.getItemCount() > 0) {
                ClipData.Item item = clip.getItemAt(0);
                Uri imageUri = item.getUri();
                Log.d("lsz","clip imageUri"+imageUri);
                if (imageUri != null) {
                    // 將圖片 URI 設置到 ImageView 顯示圖片
                    Log.d("lsz","clip imageUriaaaaaaaaaa");
                    imageView2.setImageURI(imageUri);
                    //imageView2.setImageBitmap(bitmap1);
                }
            }
        } else {        Toast.makeText(this, "Clipboard is empty", Toast.LENGTH_SHORT).show();    } */
//本地图片 取到剪切版end


        AtomicReference<ClipData> clipDataRef = new AtomicReference<>(null);
        ClipboardUtils clipboardUtils = ClipboardUtils.getInstance();
        clipboardUtils.getPrimaryClip(clipDataRef);
        Log.e("clip", "lsz len===hasClip " + clipboardUtils.hasClip());
        for (int i = 0; i < clipboardUtils.getItemCount(clipDataRef); i++) {
            Log.e("clip", "lsz len=getItemType" + clipboardUtils.getItemType(clipDataRef, i));
            if (clipboardUtils.getItemType(clipDataRef, i) == clipboardUtils.CLIPBOARD_DATA_TYPE_TEXT) {
                text = clipboardUtils.getTextItem(clipDataRef, i);
                textView.setText(text);
                //sendToPC(text);
            } else if (clipboardUtils.getItemType(clipDataRef, i) == clipboardUtils.CLIPBOARD_DATA_TYPE_IMAGE) {
                Bitmap bitmap1 = clipboardUtils.getImageItem(clipDataRef, i);
                Log.e("clip", "lsz len===hasClip bitmap1" + bitmap1);
                if (bitmap1 != null) {
                /*
                数组转bitmap
                */
                    //Bitmap drawableicon = BitmapFactory.decodeResource(getResources(), R.drawable.liu2);
                    //byte[] imageData = bitmapToByteArray(drawableicon); // 要转换的字节数组
                    //Bitmap bitmap3 = BitmapFactory.decodeByteArray(imageData, 0, imageData.length);
                    imageView2.setImageBitmap(bitmap1);

                    //byte[] imageData = bitmapToByteArray(bitmap1);
                    imageData = bitmapToByteArray(bitmap1);

                    /*//bitmap转byteArray
                    int bytes = bitmap1.getByteCount();
                    ByteBuffer buf = ByteBuffer.allocate(bytes);
                    bitmap1.copyPixelsToBuffer(buf);
                    byte[] byteArray = buf.array();
                    Log.e("lszz","GoLog  dddbyteArray="+byteArray.length);
                    sendToPCIMG(byteArray);*/


                    //Bitmap bitmap = BitmapFactory.decodeByteArray(imageData, 0, imageData.length);
                    //imageView3.setImageBitmap(bitmap);

                } else {
                    Toast.makeText(TestActivity.this, " Clipboard img is empty", Toast.LENGTH_SHORT).show();
                }


            } else {
                Log.e("clip", "not support format");
                Toast.makeText(TestActivity.this, "lsz111 Clipboard is empty", Toast.LENGTH_SHORT).show();
            }
        }


    }

    private void setClipToClipboard() {
        ClipboardUtils clipboardUtils = ClipboardUtils.getInstance();
        clipboardUtils.clearClip();

        AtomicReference<ClipData> clipDataRef = ClipboardUtils.createClipdataRef();
        clipboardUtils.addTextItem(clipDataRef, "test text12");
        clipboardUtils.addImageItem(clipDataRef, bitmap);
        // bitmapToByteArray(bitmap);
        clipboardUtils.setPrimaryClip(clipDataRef);


//本地图片 测试存到剪切版
       /* ClipboardManager clipboard = (ClipboardManager)getSystemService(Context.CLIPBOARD_SERVICE);
        ClipData clip;
        Uri uri = Uri.parse("android.resource://com.rtk.myapplication/" + R.drawable.liu2);
        clip = ClipData.newUri(getContentResolver(), "Image", uri);
        // 將 ClipData 放入剪貼簿
        clipboard.setPrimaryClip(clip);*/
//本地图片 测试存到剪切版 end

        /*Bitmap drawableicon = BitmapFactory.decodeResource(getResources(), R.drawable.liu2);
        ClipboardManager mClipboard = (ClipboardManager) getSystemService(Context.CLIPBOARD_SERVICE);
        Uri imageUri =getImageUri(this,drawableicon);
        ClipData theClip = ClipData.newUri(getContentResolver(), "Image", imageUri);
        mClipboard.setPrimaryClip(theClip);
*/

    }

    public Uri getImageUri(Context inContext, Bitmap inImage) {
        ByteArrayOutputStream bytes = new ByteArrayOutputStream();
        inImage.compress(Bitmap.CompressFormat.JPEG, 100, bytes);
        String path = MediaStore.Images.Media.insertImage(inContext.getContentResolver(), inImage, "Title", null);
        return Uri.parse(path);
    }

    private void sendToPC(String text) {
        Log.i(TAG, "lsz GoLog sendToPC: text:" + text);
        new Thread(new Runnable() {
            @Override
            public void run() {
                Libp2p_clipboard.sendMessage(text);
            }
        }).start();
    }

    /*
    bitmap转数组
     */
    public static byte[] bitmapToByteArray(Bitmap bitmap) {
        ByteArrayOutputStream outputStream = new ByteArrayOutputStream();
        bitmap.compress(Bitmap.CompressFormat.PNG, 100, outputStream);
        Log.d("lsz", "outputStream. imag toByteArray()=" + outputStream.toByteArray().toString());
        return outputStream.toByteArray();

    }

    public void testCrossShare() throws IOException {

        // Get intent, action and MIME type
        intent = getIntent();
        action = intent.getAction();
        mimetype = intent.getType();
        Log.d("lszz", "action=" + action + "/type==" + mimetype);
        if (Intent.ACTION_SEND.equals(action) && mimetype != null) {
            if (mimetype.startsWith("image/")) {
                Uri imageUri = (Uri) intent.getParcelableExtra(Intent.EXTRA_STREAM);
                InputStream inputStream = this.getContentResolver().openInputStream(imageUri);
                long sizeInBytes = getImageSize(inputStream);
                sizeInMB = bytekb(sizeInBytes);
                Log.d("lszz", "sizeInMB=" + sizeInMB);
            } else if (mimetype.startsWith("video/mp4")) {
                Uri uri = (Uri) intent.getParcelableExtra(Intent.EXTRA_STREAM);
                InputStream inputStream = this.getContentResolver().openInputStream(uri);
                long sizeInBytes = getImageSize(inputStream);
                sizeInMB = bytekb(sizeInBytes);
            }
        }

        //if (Intent.ACTION_SEND.equals(action) && mimetype != null) {
        //    alertDialog(intent,action,mimetype);
        //}

        //Log.d("lszz","action="+action + "/type="+mimetype);
        /*if (Intent.ACTION_SEND.equals(action) && mimetype != null) {
            if ("text/plain".equals(mimetype)) {
                sharedText = intent.getStringExtra(Intent.EXTRA_TEXT);
                //Uri ur = (Uri) intent.getParcelableExtra(Intent.EXTRA_STREAM);
                Log.d("lszz","sharedText="+sharedText);
                assert sharedText != null;
                textView.setText(sharedText.replace("\"", ""));
                //downLoad(ur);
            } else if (mimetype.startsWith("image/")) {
                Uri imageUri = (Uri) intent.getParcelableExtra(Intent.EXTRA_STREAM);
                Log.d("lszz","imageUri="+imageUri);
                InputStream inputStream = this.getContentResolver().openInputStream(imageUri);
                Bitmap bitmap = BitmapFactory.decodeStream(inputStream);
                //bitmap转byteArray
                int bytes = bitmap.getByteCount();
                ByteBuffer buf = ByteBuffer.allocate(bytes);
                bitmap.copyPixelsToBuffer(buf);
                //byte[] byteArray = buf.array();

                getbyteArray=buf.array();

                getbyteArray(bitmap);
                Log.d("lszz","bitmap byteArray="+getbyteArray);
                Log.d("lszz","bitmap="+bitmap);
                imageView2.setImageURI(imageUri);

                //测试获取的uri 转换成bitmap显示
                byte[] imageData = bitmapToByteArray(bitmap); // 要转换的字节数组
                Bitmap bitmap3 = BitmapFactory.decodeByteArray(imageData, 0, imageData.length);
                imageView3.setImageBitmap(bitmap3);

            }else if (mimetype.startsWith("video/mp4")) {
                Uri uri = (Uri) intent.getParcelableExtra(Intent.EXTRA_STREAM);
                //如果是媒体类型需要从数据库获取路径
                videview.setVideoURI(uri);
                videview.start();
                videview.setVisibility(View.VISIBLE);
            }else if (mimetype.startsWith("audio/mpeg")) {
                Uri uri = (Uri) intent.getParcelableExtra(Intent.EXTRA_STREAM);
                Log.i("lszz", "uri.getPath();:= " + uri.getPath());
                playAudio(uri);
            }else if (mimetype.startsWith("application/")) {
                Uri uri = (Uri) intent.getParcelableExtra(Intent.EXTRA_STREAM);
                /*verifyStoragePermissions();
                File internalDirectory = getExternalFilesDir("MyFolder");
                Log.i("lszz", "uri.getPath();:=internalDirectory " + internalDirectory.getPath());
                Log.i("lszz", "uri.getPath();:=internalDirectory " + internalDirectory.exists());
                */
                /*try {
                    InputStream inputStream = this.getContentResolver().openInputStream(uri);
                    File saveDir = Environment.getExternalStoragePublicDirectory(Environment.DIRECTORY_DOWNLOADS);//保存在内部存储的Download下
                    File saveFile = new File(saveDir, getFileNameFromUri(uri));
                    OutputStream outputStream = new FileOutputStream(saveFile);
                    byte[] buffer = new byte[1024];
                    int length = inputStream.read(buffer);
                    while (length > 0) {
                        outputStream.write(buffer, 0, length);
                        length = inputStream.read(buffer);
                    }
                    inputStream.close();
                    outputStream.close();
                    Toast.makeText(TestActivity.this, "文件已经保存在内部存储的Download下", Toast.LENGTH_SHORT).show();
                } catch (IOException e) {
                    e.printStackTrace();
                    // 处理异常
                }
            }



        }*/


    }

    private static final int REQUEST_EXTERNAL_STORAGE = 1;
    private static String[] PERMISSIONS_STORAGE = {
            android.Manifest.permission.WRITE_EXTERNAL_STORAGE
    };

    private void verifyStoragePermissions() {
        // Check if we have write permission
        int permission = ActivityCompat.checkSelfPermission(this, android.Manifest.permission.WRITE_EXTERNAL_STORAGE);

        if (permission != PackageManager.PERMISSION_GRANTED) {
            // We don't have permission so prompt the user
            ActivityCompat.requestPermissions(
                    this,
                    PERMISSIONS_STORAGE,
                    REQUEST_EXTERNAL_STORAGE
            );
        }
    }

    public String getRealPathFromURI(Context context, Uri contentUri) {
        String[] proj = {MediaStore.Images.Media.DATA};
        Cursor cursor = context.getContentResolver().query(contentUri, proj, null, null, null);
        int column_index = cursor.getColumnIndexOrThrow(MediaStore.Images.Media.DATA);
        cursor.moveToFirst();
        String path = cursor.getString(column_index);
        cursor.close();
        Log.d("lszz", "get file path=" + path);
        return path;
    }

    private String getFileNameFromUri(Uri zipUri) {
        Cursor returnCursor = getContentResolver().query(zipUri, null, null, null, null);
        /*
         * Get the column indexes of the data in the Cursor,
         * then get the data from the Cursor
         */
        int nameIndex = returnCursor.getColumnIndex(OpenableColumns.DISPLAY_NAME);
        returnCursor.moveToFirst();
        String fileName = returnCursor.getString(nameIndex);
        Log.d("lszz", "get fileName=" + fileName);
        returnCursor.close();
        return fileName;
    }

    private void playAudio(Uri audioUri) {
        MediaPlayer mediaPlayer = new MediaPlayer();
        try {
            mediaPlayer.setDataSource(this, audioUri);
            mediaPlayer.prepare();
            mediaPlayer.start();
        } catch (IOException e) {
            e.printStackTrace();
        }
    }


    private void downLoad(Uri uri) {
        try {
            InputStream inputStream = this.getContentResolver().openInputStream(uri);
            File saveDir = Environment.getExternalStoragePublicDirectory(Environment.DIRECTORY_DOWNLOADS);//保存在内部存储的Download下
            File saveFile = new File(saveDir, getFileNameFromUri(uri));
            OutputStream outputStream = new FileOutputStream(saveFile);
            byte[] buffer = new byte[1024];
            int length = inputStream.read(buffer);
            while (length > 0) {
                outputStream.write(buffer, 0, length);
                length = inputStream.read(buffer);
            }
            inputStream.close();
            outputStream.close();
            Toast.makeText(TestActivity.this, "文件已经保存在内部存储的Download下", Toast.LENGTH_SHORT).show();
        } catch (IOException e) {
            e.printStackTrace();
            // 处理异常
        }
    }


    private void sendToPCIMG(byte[] value) {

        runOnUiThread(new Runnable() {
            @Override
            public void run() {
                Log.i("lszzz", "GoLog sendToPC: img byte[] value:==" + value);
                Log.e("lszzz", "GoLog sendToPC: img byte[] value length:==" + value.length);
                //byte[] newBytes = new String(value, StandardCharsets.ISO_8859_1).getBytes(StandardCharsets.ISO_8859_1);
                //Log.i("lszz", "GoLog sendToPC: img newBytes:==" + newBytes);
                //String str = new String(value,StandardCharsets.UTF_8);
                //Log.i("lszz", "GoLog emoveInvalidCharacters=" + removeInvalidCharacters(str));
                //String str2 = new String(newBytes,StandardCharsets.UTF_8);
                //Log.i("lszz", "GoLog str: str str2:==" + removeInvalidCharacters(str2));

                base64String = Base64.encodeToString(value, Base64.DEFAULT);

                //Libp2p_clipboard.sendImage(removeInvalidCharacters(str));
                //Libp2p_clipboard.sendImage(removeInvalidCharacters(base64String));
                clearbase64String = removeInvalidCharacters(base64String);
                //Log.i("lszz", "GoLog android======clearbase64String.getBytes()" + clearbase64String.getBytes());
                Log.i("lszz", "GoLog android======" + clearbase64String);
                Log.i("lszz", "GoLog clearbase64String.length()=" + clearbase64String.length());
                //Libp2p_clipboard.sendImage(clearbase64String);
                //for (byte getlistvalueaa : value) {
                //    Log.e("lszz","GoLog getlistvalueaa="+getlistvalueaa);
                // }


                //Libp2p_clipboard.sendImage(value);

                /*Log.i("lszz", "GoLog base64String.length()======" + base64String.length());
                Log.i("lszz", "GoLog value.length======" + value.length);

                String aa="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";
                Log.i("lszz", "GoLog aa =aa aa=aaa" + aa);
                Log.i("lszz", "GoLog test aa=aa=" + aa.length());
                Libp2p_clipboard.sendImage(aa);

                String bb = "/9j/4AAQSkZJRgABAQAAAQABAAD/4gIoSUNDX1BST0ZJTEUAAQEAAAIYAAAAAAQwAABtbnRyUkdCIFhZWiAAAAAAAAAAAAAAAABhY3NwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAA9tYAAQAAAADTLQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAlkZXNjAAAA8AAAAHRyWFlaAAABZAAAABRnWFlaAAABeAAAABRiWFlaAAABjAAAABRyVFJDAAABoAAAAChnVFJDAAABoAAAAChiVFJDAAABoAAAACh3dHB0AAAByAAAABRjcHJ0AAAB3AAAADxtbHVjAAAAAAAAAAEAAAAMZW5VUwAAAFgAAAAcAHMAUgBHAEIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFhZWiAAAAAAAABvogAAOPUAAAOQWFlaIAAAAAAAAGKZAAC3hQAAGNpYWVogAAAAAAAAJKAAAA+EAAC2z3BhcmEAAAAAAAQAAAACZmYAAPKnAAANWQAAE9AAAApbAAAAAAAAAABYWVogAAAAAAAA9tYAAQAAAADTLW1sdWMAAAAAAAAAAQAAAAxlblVTAAAAIAAAABwARwBvAG8AZwBsAGUAIABJAG4AYwAuACAAMgAwADEANv/bAEMAAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAf/bAEMBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAf/AABEIB/wGQAMBIgACEQEDEQH/xAAfAAABBAMBAQEBAAAAAAAAAAAAAgQFBgMHCAEJCgv/xABqEAAAAgUIBwYEAwQFBwgGARUBAgADBBEhBQYSMUFRYfATcYGRobHBBxQiMtHhCCMz8RVCQwkkUlMWNGJjgiVEVGRyc6IXNXSDkpPC4gomhJSjstPyGDZFpNJVtCdGxNTkGWV1N4Wls7X/xAAeAQACAwEBAQEBAQAAAAAAAAAABAIDBQYHAQgJCv/EAE8RAAACBQkHAgUBBgQFAwQABwABAxEhMUEEE1FhcYGhwfACBQaRsdHhEiMHFDNy8VIVJEJDYsIWgqLSCCIyNJIXROIlNVNUJmNFshhlc//aAAwDAQACEQMRAD8A/dQggiDD4Xhbn2QAAw+F4W59kUaodQ8k9Tw1Q6h5IwAMESfyjs5gikSfyjs5gliON2YBhQQQQSQvyABBBBF5uvDyACCCCWo0mokcDI2axAIIIJObrw8gAnhqh1DyT1PDVDqHkgjjdmAYEEEEfRpDmuT324NsvABEn8o7OYIpEn8o7OYIoAYUEEEYABEn8o7OYIpEn8o7OYJWkhfkAYUEEEXm68PIAIIIJfsPOzMgAQQQSAAJ4aodQ8k9Tw1Q6h5IAGBBBBAAEEEEAAQQQQABBBBAAEEEEAAQQQQABBBBAAESfyjs5gikSfyjs5ggAMUEEEAAQQQQABBBBAAE8NUOoeSep4aodQ8kADAggggFwIIIIAAggggLEcbswIk/lHZzBFIk/lHZzBAWBufyjs5ggY1HWKKQQBNf04+Qk/lHZzBFIkw+F99W32RSAZk7789kJP5R2cwRr+pn+FHX5P8AD0Rr+pn+FACRJqJnEzNusHX5P8PRGh6w1dRR3+T/AA9EaHrDV1FAZwbLbf8AZHqiUyLfOO3mKY0AwBEn8o7OYIpEn8o7OYIADFBHGg/t/wDD/wCZE6H+2VAVpIX5DCgmUyqjUIYZhyRNAbw4+iT9B0lj2CM37ueNP+bAZC1BqDkjlBEmHwvvq2+yfPqUOrcvq0XI0buZdzw6VAplv4D6IUy38B9EUgkQz6CpPDsBBElHwvur2eyKQFQEEEEAAQQQQABEn8o7OYIpEn8o7OYIADFBBBARna9rn5AggiT+UdnMEBFJC/IYUEEEAnNzd+CripAggggGUaTUSOBkbNYiCJplv4D6IUy38B9EBObrw8hSCJplv4D6IUy38B9EBWFIIIIBgCCCCAAIIJk0ePD3QAMaCZNHjw90xoAAggggACJP5R2cwRSJP5R2cwQAMKCCCArSJH3kfYsetYEEEEAmjRvWeGBaoYBBBBAbEnSOzvYvmXICCCCAb9aOgv8AT3Hhag1ByT1BBAKATw1Q6h5J6nhqh1DyQAMCCCCAAIIIIAAiT+UdnMEUiT+UdnMEADCggggACCCCAAIknlDbzFFIko+F91ez2QFaSF+QUggggEUiOw1lz8+GgQQQQHz0FSeHYZieUNvMUUiSeUNvMUUgL0aN3Mu54dKgIk/lHZzBFIk/lHZzBAWBufyjs5ggfyjs5ggfyjs5ggfyjs5ggLEcbswE8obeYoH8o7OYIeP+zxQ8f9nigGfWVB4dwE8obeYoE8obeYopEk8obeYoBZJC/IKRJ/KOzmCKRJ/KOzmCArGFBBEGHwvC3PsgABh8Lwtz7ItBBAWo0aqlLUS3dadQQYfC8Lc+yKLUGoOSI/7tFeP+zxQDY9NUOoeSBag1ByRNLwvtq25jwQKPheNmfZAVpIX5AKPheNmfZPT+UdnMEUiD1Br6CgFvQdJY9go1Q6h5Ij9PP8SenqDX0FPPmZooBlHG7MLLUGoOSeoIICsCCCCAAs/0w1f+EUhTVjrHmk0f6Yav/CKQpqx1jzQC4yk8obeYopEk8obeYopAAEEEEAAQQQQABBBBAAEEEEAAQQQQAH5ag1ByT1PC1BqDknqAAIIIIAAggggACCCCAAIIIIAAggnhqh1DyQAPUEEEBb6CpPDsOpTD4Xhbn2RaCCLiYE8NUOoeSep4aodQ8kABgiT+UdnMEUiT+UdnMEXAMKCCCOI43ZgAggghN14eQAQQQRdHG7MAEEEEcABPDVDqHknqeGqHUPJKUsftPMAwIIIJajSagZQMjbrAAgggn0AEEEEAATwwPLy2Xck9RJ/KOzmCT9Z0Fj3AG9At3EfVCgW7iPqikEgAJNX5H2auY8ooUC3cR9UD+UdnMEUi4AmgW7iPqhQLdxH1RSCWI43ZgCaBbuI+qBg8LrqtnsikSfyjs5glgBhQQQQABEGDwuCzPui08NUOoeSABgQQQQABBBBAAEEEEAAQQQQABBBBAAEEEEAAQQQQABBBBAAEEEEAATwwPLy2Xck9RJ/KOzmCAAxQQQQC4EQYPC4LM+6LTw1Q6h5IAGBBBBAAESfyjs5gikSfyjs5ggLEcbsw3plv4D6IGHwvvq2+yB/KOzmCKQF6NGuta1Et/SjURElDwuvr2+yKRJPKG3mKAltuK3Iwfk/w9EwpmpeGlhxq3PTCgEhms/6n/wAaND1hq6ijul4aWHGrc9Gh6w1dRQDaNGuta1Et/SjUcC3zjt5imNMlGl4nufZXVDC5DR48PdAQGNBMmjx4e6Gjx4e6AAaPHh7oaPHh7pkQQAMNAP5n/wAqFAP5n/ypmTwwuLy23c0BD0FSeHYNj+UdnMEUiT+UdnMEUgPiON2YEEEEBYBEk8obeYopElHwvur2eyArSQvyCkEboJb6CpPDsF5yrHwHCCJplv4D6IUy38B9EqFgUjdHCN0sRxuzFaSF+QboIIJWF5yrHwBBBEn8o7OYIBhHG7MYUEEEATdeHkCJP5R2cwROkw4+yKpeGlhxq3PQFKMpvOwYUExUxuDj6pj03jJfqt5PvsuQDvoOksewcoIhWBzU6EdGNV8dnoOpI5onAzSfT0x1ZNEGRhc8YPrS2bSpCjqNdHYKpEmomcTM26wnTBRe/fgiDmIUo09nLr6Ot1XLHalJTOJ6Hj0dWyvbu1pq6WO1hsav6n9H+7s+0OO2PoOksewtna9rn5HUWnVfxhx9E970p/jV7hTkFX2hS4uLQUnWetd9mG";
                Log.i("lszz", "GoLog aa =bb bb=bbc" + bb);
                Log.i("lszz", "GoLog test =bb bb=bb=" + bb.length());
                Libp2p_clipboard.sendImage(bb);*/


            }
        });


        //byte[] valuea =Base64.decode(clearbase64String, Base64.DEFAULT);
        //Log.i("lszzz", "GoLog sendToPC: img value:==" + value);
        //Bitmap bitmap = BitmapFactory.decodeByteArray(valuea, 0, valuea.length);
        //imageView3.setImageBitmap(bitmap);


    }

    public static String encodeToStringWithoutSlash(byte[] data) {
        String base64String = Base64.encodeToString(data, Base64.DEFAULT);
        return base64String.replace("/", "");
    }

    public String base64StringToNormalString(String base64String) {
        byte[] decodedBytes = Base64.decode(base64String, Base64.DEFAULT);
        return new String(decodedBytes);
    }

    public String byteToString(byte[] data) {
        int index = data.length;
        for (int i = 0; i < data.length; i++) {
            if (data[i] == 0) {
                index = i;
                break;
            }
        }
        byte[] temp = new byte[index];
        Arrays.fill(temp, (byte) 0);
        System.arraycopy(data, 0, temp, 0, index);
        String str;
        try {
            str = new String(temp, "GBK");
        } catch (UnsupportedEncodingException e) {
            // TODO Auto-generated catch block
            e.printStackTrace();
            return "";
        }
        return str;
    }


    public byte[] getbyteArray(Bitmap bitmap) {
        int bytes = bitmap.getByteCount();
        ByteBuffer buf = ByteBuffer.allocate(bytes);
        bitmap.copyPixelsToBuffer(buf);
        byte[] byteArray = buf.array();
        return byteArray;
    }


    public void alertDialog(String s, long l) {
        LayoutInflater inflater = LayoutInflater.from(this);
        View customView = inflater.inflate(R.layout.custom_dialog, null);

        TextView textView = customView.findViewById(R.id.textclientname);
        TextView textsize = customView.findViewById(R.id.textsize);
        Button buttonPositiveButton = customView.findViewById(R.id.buttonPositiveButton);
        Button buttonNegativeButton = customView.findViewById(R.id.buttonNegativeButton);
        CheckBox myCheckBox = customView.findViewById(R.id.my_checkbox);

        textView.setText(s);
        textsize.setText(String.valueOf(l));
        Log.d("lszz", "lszz String s===s=" + s);
        Log.d("lszz", "lszz long   l===l=" + l);
        AlertDialog.Builder builder = new AlertDialog.Builder(this);

        builder.setView(customView);

        AlertDialog dialog = builder.create();

        //boolean isChecked = myCheckBox.isChecked();
        myCheckBox.setOnCheckedChangeListener(new CompoundButton.OnCheckedChangeListener() {
            @Override
            public void onCheckedChanged(CompoundButton compoundButton, boolean b) {
                if (b) {
                    Log.d("lszz", "lszz set checkbox is =" + b);
                    kv.encode("ischeck", true);
                } else {
                    Log.d("lszz", "lszz set checkbox is =" + b);
                    kv.encode("ischeck", false);
                }
            }
        });


        buttonPositiveButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {
                Libp2p_clipboard.ifClipboardPasteFile(true);
                Log.d("lszz", "lszz buttonPositiveButton");
                dialog.dismiss();
            }
        });
        buttonNegativeButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {
                Log.d("lszz", "lszz buttonNegativeButton");
                Libp2p_clipboard.ifClipboardPasteFile(false);
                dialog.dismiss();
            }
        });
        buttonPositiveButton.requestFocus();


        Window window = dialog.getWindow();
        if (window != null) {
            WindowManager.LayoutParams layoutParams = new WindowManager.LayoutParams();
            layoutParams.copyFrom(window.getAttributes());

            // 设置对话框的位置为左上角
            layoutParams.gravity = Gravity.START | Gravity.TOP;

            // 可以设置偏移量，如果不设置偏移量，对话框将紧贴屏幕左上角
            // layoutParams.x = 0; // 默认为0，不需要设置
            // layoutParams.y = 0; // 默认为0，不需要设置

            window.setAttributes(layoutParams);
        }

        dialog.show();

    }

    public void alertDialog(Intent intent, String action, String mimetype) {
        LayoutInflater inflater = LayoutInflater.from(this);
        View customView = inflater.inflate(R.layout.custom_dialog, null);

        TextView textView = customView.findViewById(R.id.textclientname);
        TextView textsize = customView.findViewById(R.id.textsize);
        textView.setText(value);
        textsize.setText(sizeInMB);

        AlertDialog.Builder builder = new AlertDialog.Builder(this);
        builder.setView(customView);


// 添加按钮
        builder.setPositiveButton("确定", new DialogInterface.OnClickListener() {
            @Override
            public void onClick(DialogInterface dialog, int which) {
                // 确定按钮的点击事件
               /* try {
                    getShare(intent, action,mimetype);
                } catch (FileNotFoundException e) {
                    throw new RuntimeException(e);
                } catch (Exception e) {
                    throw new RuntimeException(e);
                }*/
                Libp2p_clipboard.ifClipboardPasteFile(true);
            }
        });
        builder.setNegativeButton("取消", new DialogInterface.OnClickListener() {
            @Override
            public void onClick(DialogInterface dialog, int which) {
                // 取消按钮的点击事件
            }
        });

// 创建并显示对话框
        AlertDialog dialog = builder.create();
        dialog.show();


    }

    private long getImageSize(InputStream inputStream) throws IOException {
        return inputStream.available(); // 返回输入流中可用的字节数
    }


    public static String bytekb(long bytes) {
//格式化小数
        int GB = 1024 * 1024 * 1024;
        int MB = 1024 * 1024;
        int KB = 1024;

        DecimalFormat format = new DecimalFormat("#.##");
        if (bytes / GB >= 1) {
            return format.format(bytes / GB) + "GB";
        } else if (bytes / MB >= 1) {
            return format.format(bytes / MB) + "MB";
        } else if (bytes / KB >= 1) {
            return format.format(bytes / KB) + "KB";
        } else {
            return bytes + "B";
        }
    }

    public void getShare(Intent intent, String action, String mimetype) throws Exception {

        Log.d("lszz", "action=" + action + "/type=" + mimetype);
        if (Intent.ACTION_SEND.equals(action) && mimetype != null) {
            if ("text/plain".equals(mimetype)) {
                sharedText = intent.getStringExtra(Intent.EXTRA_TEXT);
                //Uri ur = (Uri) intent.getParcelableExtra(Intent.EXTRA_STREAM);
                Log.d("lszz", "sharedText=" + sharedText);
                assert sharedText != null;
                textView.setText(sharedText.replace("\"", ""));
                //downLoad(ur);
            } else if (mimetype.startsWith("image/")) {
                Uri imageUri = (Uri) intent.getParcelableExtra(Intent.EXTRA_STREAM);
                Log.d("lszz", "imageUri=" + imageUri);
                InputStream inputStream = this.getContentResolver().openInputStream(imageUri);

                long sizeInBytes = getImageSize(inputStream);
                sizeInMB = bytekb(sizeInBytes);
                Log.d("lszz", "sizeInMB=" + sizeInMB);

                bitmapShare = BitmapFactory.decodeStream(inputStream);
                imageData = bitmapToByteArray(bitmapShare);
                //bitmap转byteArray
               /* int bytes = bitmap.getByteCount();
                ByteBuffer buf = ByteBuffer.allocate(bytes);
                bitmap.copyPixelsToBuffer(buf);
                //byte[] byteArray = buf.array();

                getbyteArray=buf.array();

                getbyteArray(bitmap);
                Log.d("lszz","bitmap byteArray="+getbyteArray);*/
                Log.d("lszz", "bitmap=bitmapShare=" + bitmapShare);
                imageView2.setImageURI(imageUri);

                //获取uri 存入Provider
                InputStream inputStream2 = this.getContentResolver().openInputStream(imageUri);
                Bitmap bitmap = BitmapFactory.decodeStream(inputStream2);
                ByteArrayOutputStream stream = new ByteArrayOutputStream();
                bitmap.compress(Bitmap.CompressFormat.PNG, 100, stream);
                byte[] byteArray = stream.toByteArray();
                ContentValues values = new ContentValues();
                values.put("image", byteArray);
                Uri newUri = getContentResolver().insert(UserContentProvider.CONTENT_URI, values);
                getContentResolver().update(UserContentProvider.CONTENT_URI, values, null, null);
                //获取uri 存入Provider end

                /*//测试获取的uri 转换成bitmap显示
                byte[] imageData = bitmapToByteArray(bitmap); // 要转换的字节数组
                Bitmap bitmap3 = BitmapFactory.decodeByteArray(imageData, 0, imageData.length);

                imageView3.setImageBitmap(bitmap3);*/

            } else if (mimetype.startsWith("video/mp4")) {
                Uri uri = (Uri) intent.getParcelableExtra(Intent.EXTRA_STREAM);
                //如果是媒体类型需要从数据库获取路径
                videview.setVideoURI(uri);
                videview.start();
                videview.setVisibility(View.VISIBLE);
            } else if (mimetype.startsWith("audio/mpeg")) {
                Uri uri = (Uri) intent.getParcelableExtra(Intent.EXTRA_STREAM);
                Log.i("lszz", "uri.getPath();:= " + uri.getPath());
                playAudio(uri);
            } else if (mimetype.startsWith("application/")) {
                Uri uri = (Uri) intent.getParcelableExtra(Intent.EXTRA_STREAM);

                getRealPathFromURI(mContext, uri);
                //verifyStoragePermissions();
                //File internalDirectory = getExternalFilesDir("MyFolder");
                //Log.i("lszz", "uri.getPath();:=internalDirectory " + internalDirectory.getPath());
                //Log.i("lszz", "uri.getPath();:=internalDirectory " + internalDirectory.exists());

                try {
                    InputStream inputStream = this.getContentResolver().openInputStream(uri);
                    File saveDir = Environment.getExternalStoragePublicDirectory(Environment.DIRECTORY_DOWNLOADS);//保存在内部存储的Download下
                    File saveFile = new File(saveDir, getFileNameFromUri(uri));
                    Log.i("lszz", "uri.getPath();:=saveDir " + saveDir.getPath());
                    Log.i("lszz", "uri.getPath();:=saveFile " + saveFile.getPath());
                    OutputStream outputStream = new FileOutputStream(saveFile);
                    byte[] buffer = new byte[1024];
                    int length = inputStream.read(buffer);
                    while (length > 0) {
                        outputStream.write(buffer, 0, length);
                        length = inputStream.read(buffer);
                    }
                    inputStream.close();
                    outputStream.close();
                    Toast.makeText(TestActivity.this, "文件已经保存在内部存储的Download下", Toast.LENGTH_SHORT).show();
                } catch (IOException e) {
                    e.printStackTrace();
                    // 处理异常
                }
            }


        }

    }

    public static String removeInvalidCharacters(String base64String) {
        // 正则表达式，匹配Base64的有效字符
        String regex = "[^A-Za-z0-9+/=]";
        // 使用正则表达式替换掉非法字符
        String cleanString = base64String.replaceAll(regex, "");
        return cleanString;
    }

    @Override
    protected void onStop() {
        Log.d(TAG, "lsz onStop");
        kv.encode("ischeck", false);
        //finish();
        super.onStop();
    }

    @Override
    protected void onDestroy() {
        Log.d(TAG, "lsz onDestroy");
        super.onDestroy();
    }

    public String getNetwork() {

        String name = "";
        try {
            List<NetworkInterface> interfaces = Collections.list(NetworkInterface.getNetworkInterfaces());
            for (NetworkInterface intf : interfaces) {
                if (intf.getName().startsWith("wlan")) {

                    return "wlan0";
                } else if (intf.getName().startsWith("eth")) {
                    return "ethernet";
                }
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
        return name;
    }


    public static boolean checkFloatPermission(Context context) {
        if (Build.VERSION.SDK_INT < Build.VERSION_CODES.KITKAT) {
            return true;
        }
        if (Build.VERSION.SDK_INT < Build.VERSION_CODES.M) {
            try {
                Class<?> cls = Class.forName("android.content.Context");
                Field declaredField = cls.getDeclaredField("APP_OPS_SERVICE");
                declaredField.setAccessible(true);
                Object obj = declaredField.get(cls);
                if (!(obj instanceof String)) {
                    return false;
                }
                String str2 = (String) obj;
                obj = cls.getMethod("getSystemService", String.class).invoke(context, str2);
                cls = Class.forName("android.app.AppOpsManager");
                Field declaredField2 = cls.getDeclaredField("MODE_ALLOWED");
                declaredField2.setAccessible(true);
                Method checkOp = cls.getMethod("checkOp", Integer.TYPE, Integer.TYPE, String.class);
                int result = (Integer) checkOp.invoke(obj, 24, Binder.getCallingUid(), context.getPackageName());
                return result == declaredField2.getInt(cls);
            } catch (Exception e) {
                return false;
            }
        } else {
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
                AppOpsManager appOpsMgr = (AppOpsManager) context.getSystemService(Context.APP_OPS_SERVICE);
                if (appOpsMgr == null) return false;
                int mode = appOpsMgr.checkOpNoThrow("android:system_alert_window", android.os.Process.myUid(), context.getPackageName());
                return mode == AppOpsManager.MODE_ALLOWED || mode == AppOpsManager.MODE_IGNORED;
            } else {
                return Settings.canDrawOverlays(context);
            }
        }
    }


    public void getClientList() {
        RecyclerView recyclerView = findViewById(R.id.my_recycler_view);
        recyclerView.setLayoutManager(new LinearLayoutManager(this, LinearLayoutManager.VERTICAL, false));
        //从libp2p获取列表数据
        String getlist = Libp2p_clipboard.getClientList();
        Log.d("lszz", "getlist==" + getlist);
        String[] strArray = getlist.split("#");
        //for (String getlistvalue : strArray) {
        //    Log.d("lszz","getlistvalue="+getlistvalue);
        //}

            /*for (int i = 0; i < strArray.length; i++) {
                if (!strArray[i].isEmpty()) {
                    Log.d("lszz","subString=aaaa"+strArray[i].substring(0, 10));
                    strArray[i] = strArray[i].substring(0, 10);
                }
            }*/

        //取到的数据放入myadapter
        MyAdapter myadapter = new MyAdapter(strArray);
        myadapter.setOnItemClickListener(new MyAdapter.OnItemClickListener() {
            @Override
            public void onItemClick(View view, int position) {
                value = ((TextView) view).getText().toString();
                Toast.makeText(TestActivity.this, "你选择了：" + value, Toast.LENGTH_SHORT).show();

                //if (intent != null && Intent.ACTION_SEND.equals(action) && mimetype != null) {
                //    alertDialog(intent, action, mimetype);
                //}
                    /*Uri imageUri = (Uri) intent.getParcelableExtra(Intent.EXTRA_STREAM);
                    Log.d("lszz","imageUri="+imageUri);
                    //InputStream inputStream = this.getContentResolver().openInputStream(imageUri);
                    //Bitmap bitmap = BitmapFactory.decodeStream(inputStream);

                    Log.d("lszz","imageUri=bitmapShare=aaa="+bitmapShare);
                    if(imageUri != null && bitmapShare != null) {
                        byte[] imageData = bitmapToByteArray(bitmapShare);
                        sendToPCIMG(imageData);
                    }*/

            }
        });

        recyclerView.setAdapter(myadapter);
    }

}
