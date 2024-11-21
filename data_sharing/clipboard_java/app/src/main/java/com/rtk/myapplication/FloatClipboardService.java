package com.rtk.myapplication;

import android.app.AlertDialog;
import android.app.Instrumentation;
import android.app.Notification;
import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.app.PendingIntent;
import android.app.Service;
import android.content.ClipData;
import android.content.ClipboardManager;
import android.content.Context;
import android.content.DialogInterface;
import android.content.Intent;
import android.content.pm.ServiceInfo;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.graphics.Color;
import android.graphics.PixelFormat;
import android.hardware.input.InputManager;
import android.net.ConnectivityManager;
import android.net.LinkAddress;
import android.net.LinkProperties;
import android.os.Build;
import android.os.Handler;
import android.os.IBinder;
import android.os.Looper;
import android.util.Base64;
import android.util.Log;
import android.view.Gravity;
import android.view.KeyEvent;
import android.view.LayoutInflater;
import android.view.MotionEvent;
import android.view.View;
import android.view.WindowManager;
import android.view.inputmethod.InputMethodManager;
import android.widget.Button;
import android.widget.ImageView;
import android.widget.TextView;
import android.widget.Toast;

import java.io.ByteArrayOutputStream;
import java.nio.ByteBuffer;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.List;
import java.util.concurrent.atomic.AtomicReference;

import libp2p_clipboard.Callback;
import libp2p_clipboard.Libp2p_clipboard;

public class FloatClipboardService extends Service implements View.OnKeyListener{

    private static final String TAG = FloatClipboardService.class.getSimpleName();
    private WindowManager windowManager;
    private WindowManager windowManager2;
    private View floatView;
    private View floatView2;
    private ClipboardManager clipboardManager;
    private Notification notification;
    private WindowManager.LayoutParams params;
    private WindowManager.LayoutParams params2;
    private int testCount = 0;
    private float initialX, initialY;

    TextView textView;

    private View floatLayout;

    String previousText = "";
    String previousImgMD5 = "";
    ImageView  imageView;
    private View view;
    @Override
    public int onStartCommand(Intent intent, int flags, int startId) {
        return START_STICKY;
    }

    @Override
    public void onCreate() {
        super.onCreate();
        initService();
        setClibMessageLoop();

//lsz add
        //start();
        //init();


        testClipboardUtils();
    }


    public void testClipboardUtils(){

        ClipboardUtils.setContext(getApplication());
        ClipboardUtils clipboardUtils = ClipboardUtils.getInstance();
        //Log.d("lsz","GoLog clipboardUtils.hasClip()=s=="+clipboardUtils.hasClip());
        if(clipboardUtils.hasClip()){
            getClipFromClipboard();
        }else{
            Toast.makeText(this, "lsz get Clipboard is empty", Toast.LENGTH_SHORT).show();
        }

    }


    public void TestThread(){

        new Thread(new Runnable() {
            @Override
            public void run() {
                sendKeycode(KeyEvent.KEYCODE_BACK);
            }
        }).start();
    }
    public void sendKeycode(int i) {

        try {
            Instrumentation inst = new Instrumentation();
            inst.sendKeyDownUpSync(i);
        } catch (Exception e) {
            e.printStackTrace();
        }


    }

    private void init() {

        floatLayout = new Button(this);
        floatLayout.setBackgroundColor(Color.RED);

        // 设置悬浮窗布局参数
         params = new WindowManager.LayoutParams(
                WindowManager.LayoutParams.WRAP_CONTENT,
                WindowManager.LayoutParams.WRAP_CONTENT,
                WindowManager.LayoutParams.TYPE_APPLICATION_OVERLAY,

                 WindowManager.LayoutParams.FLAG_NOT_TOUCH_MODAL,


                //WindowManager.LayoutParams.FLAG_NOT_FOCUSABLE,
                PixelFormat.TRANSLUCENT);

        // 将布局添加到WindowManager
        windowManager = (WindowManager) getSystemService(WINDOW_SERVICE);
        windowManager.addView(floatLayout, params);

        floatLayout.setOnKeyListener(new View.OnKeyListener() {


            @Override
            public boolean onKey(View v, int keyCode, KeyEvent event) {
                Log.d(TAG, "Back Back key pressed"+keyCode);
               if (keyCode == KeyEvent.KEYCODE_BACK && event.getAction() == KeyEvent.ACTION_DOWN) {
                    // 处理返回键事件
                    Log.d(TAG, "lsz Back Back key pressed");
                //  TestThread();
                    return true;
                }
                return false;
            }
        });

        Log.d("lsz","float============");



        // 设置触摸事件监听
        floatLayout.setOnTouchListener(new View.OnTouchListener() {
            @Override
            public boolean onTouch(View v, MotionEvent event) {
                switch (event.getAction()) {
                    case MotionEvent.ACTION_DOWN:
                        // 记录初始坐标
                        initialX = event.getRawX() - params.x;
                        initialY = event.getRawY() - params.y;
                        break;
                    case MotionEvent.ACTION_MOVE:
                        // 更新悬浮窗位置
                        params.x = (int) (event.getRawX() - initialX);
                        params.y = (int) (event.getRawY() - initialY);
                        windowManager.updateViewLayout(v, params);
                        break;
                    case MotionEvent.ACTION_UP:
                        // 可以考虑在这里处理抬起手指后的操作，例如隐藏悬浮窗
                        break;
                }
                return false;
            }
        });
        createNotification();

        // 监听剪贴板变化
       // listenClipboard();
        new Thread(new Runnable() {
            @Override
            public void run() {
                //Libp2p_clipboard.sendAddrsFromJava(getIpAddres());
                //Libp2p_clipboard.mainInit(getGolangCallBack(), LibP2pUtils.mServerId, LibP2pUtils.mServerIpInfo);
            }
        }).start();
    }


    void start() {
     /*   Intent serviceIntent = new Intent(this, TestService.class);
        if (android.os.Build.VERSION.SDK_INT >= android.os.Build.VERSION_CODES.O) {
            startForegroundService(serviceIntent);
        } else {
            startService(serviceIntent);
        }*/

    }

    private void initService() {

        // 创建悬浮窗视图
        windowManager = (WindowManager) getApplicationContext().getSystemService(WINDOW_SERVICE);
        LayoutInflater inflater = (LayoutInflater) getSystemService(LAYOUT_INFLATER_SERVICE);
        floatView = inflater.inflate(R.layout.float_window, null);

        // 设置悬浮窗参数
        params = new WindowManager.LayoutParams(
                WindowManager.LayoutParams.WRAP_CONTENT,
                WindowManager.LayoutParams.WRAP_CONTENT,
                WindowManager.LayoutParams.TYPE_APPLICATION_OVERLAY,
                WindowManager.LayoutParams.FLAG_NOT_TOUCH_MODAL

                    | WindowManager.LayoutParams.FLAG_NOT_FOCUSABLE,

               // WindowManager.LayoutParams.FLAG_NOT_TOUCH_MODAL,
               // WindowManager.LayoutParams.FLAG_WATCH_OUTSIDE_TOUCH,


                PixelFormat.TRANSLUCENT);

        params.gravity = Gravity.TOP | Gravity.START;
        params.x = 0;
        params.y = 100;


        /*floatView.setOnKeyListener(new View.OnKeyListener() {
            @Override
            public boolean onKey(View view, int i, KeyEvent keyEvent) {
                Log.i(TAG, "lsz back11122view="+view);
                if (i == keyEvent.KEYCODE_BACK) {
                    Log.i(TAG, "back11122");
                }
                return false;
            }
        });*/


        floatView.setOnKeyListener(new View.OnKeyListener() {
            @Override
            public boolean onKey(View v, int keyCode, KeyEvent event) {
                if (keyCode == KeyEvent.KEYCODE_BACK && event.getAction() == KeyEvent.ACTION_UP) {
                    // 处理返回键事件
                    Log.i(TAG, "lsz back11122view=aaa");
                    return true; // 返回true表示消费了这个事件，不会继续传递
                }
                return false;
            }
        });

        // 添加悬浮窗到WindowManager
        windowManager.addView(floatView, params);

        floatView.findViewById(R.id.float_text_view).setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                // 处理点击事件
                Log.i(TAG, "lsz back11122view=");
                params.flags = WindowManager.LayoutParams.FLAG_NOT_TOUCH_MODAL;
                windowManager.updateViewLayout(floatView, params);
            }
        });

        //Log.d("lsz","floatView.isFocused()=="+floatView.isFocused());
        //Log.d("lsz","floatView.hasFocus()=="+floatView.hasFocus());




        windowManager2 = (WindowManager) getApplicationContext().getSystemService(WINDOW_SERVICE);
        LayoutInflater inflater2 = (LayoutInflater) getSystemService(LAYOUT_INFLATER_SERVICE);
        floatView2 = inflater2.inflate(R.layout.float_window, null);


        // 设置悬浮窗参数
        params2 = new WindowManager.LayoutParams(
                WindowManager.LayoutParams.WRAP_CONTENT,
                WindowManager.LayoutParams.WRAP_CONTENT,
                WindowManager.LayoutParams.TYPE_APPLICATION_OVERLAY,
                WindowManager.LayoutParams.FLAG_NOT_TOUCH_MODAL,


                //WindowManager.LayoutParams.FLAG_NOT_FOCUSABLE,

                // WindowManager.LayoutParams.FLAG_NOT_TOUCH_MODAL,
                // WindowManager.LayoutParams.FLAG_WATCH_OUTSIDE_TOUCH,


                PixelFormat.TRANSLUCENT);

        params2.gravity = Gravity.TOP | Gravity.START;
        params2.x = 0;
        params2.y = 300;
        //windowManager2.addView(floatView2, params2);

        /*Log.d("lsz","floatView.isFocused()=="+floatView.isFocused());
        Log.d("lsz","floatView.hasFocus()=="+floatView.hasFocus());
        Log.d("lsz","floatView2.isFocused()=="+floatView2.isFocused());
        Log.d("lsz","floatView2.hasFocus()=="+floatView2.hasFocus());
*/


        // 设置触摸事件监听
        floatView.setOnTouchListener(new View.OnTouchListener() {
            @Override
            public boolean onTouch(View v, MotionEvent event) {
                switch (event.getAction()) {
                    case MotionEvent.ACTION_DOWN:
                        // 记录初始坐标
                        initialX = event.getRawX() - params.x;
                        initialY = event.getRawY() - params.y;
                        break;
                    case MotionEvent.ACTION_MOVE:
                        // 更新悬浮窗位置
                        params.x = (int) (event.getRawX() - initialX);
                        params.y = (int) (event.getRawY() - initialY);
                        windowManager.updateViewLayout(v, params);
                        break;
                    case MotionEvent.ACTION_UP:
                        // 可以考虑在这里处理抬起手指后的操作，例如隐藏悬浮窗
                        break;
                }
                return false;
            }
        });
        // 创建通知
        createNotification();

        // 监听剪贴板变化

    //    listenClipboard();
        new Thread(new Runnable() {
            @Override
            public void run() {
                Libp2p_clipboard.sendAddrsFromJava(getIpAddres());
                //Libp2p_clipboard.mainInit(getGolangCallBack(), LibP2pUtils.mServerId, LibP2pUtils.mServerIpInfo);
                //Libp2p_clipboard.mainInit(getGolangCallBack(), "QmYsgB9x85LDSn8aon1kaaFnJP5LbZNv3apJVi67gnr8gv", "/ip4/192.168.2.107/tcp/7999/p2p/");

            }
        }).start();
    }

    private String getIpAddres() {
        ConnectivityManager connectivityManager = (ConnectivityManager) getApplicationContext().getSystemService(Service.CONNECTIVITY_SERVICE);
        LinkProperties linkProperties = connectivityManager.getLinkProperties(connectivityManager.getActiveNetwork());
        List<LinkAddress> addressList = linkProperties.getLinkAddresses();
        StringBuffer sbf = new StringBuffer();
        for (LinkAddress linkAddress : addressList) {
            sbf.append(linkAddress.toString()).append("#");
        }
        Log.d(TAG, "getIpAddres: " + sbf.toString());
        return sbf.toString();
    }

    private void updateFocus(boolean focusable) {
        if (focusable) {
            params.flags = WindowManager.LayoutParams.FLAG_NOT_TOUCH_MODAL;
            windowManager.updateViewLayout(floatView, params);
        } else {
            params.flags = WindowManager.LayoutParams.FLAG_NOT_FOCUSABLE
                    | WindowManager.LayoutParams.FLAG_NOT_TOUCH_MODAL;
            windowManager.updateViewLayout(floatView, params);
        }
    }

    private void createNotification() {
        // 创建通知渠道（API 26+）
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            CharSequence name = "Float Clipboard Service";
            String channelId = "float_clipboard_channel_id";
            NotificationChannel channel = new NotificationChannel(channelId, name, NotificationManager.IMPORTANCE_LOW);
            NotificationManager notificationManager = (NotificationManager) getSystemService(Context.NOTIFICATION_SERVICE);
            notificationManager.createNotificationChannel(channel);
        }

        // 创建通知
        Intent intent = new Intent(this, FloatClipboardService.class);
        PendingIntent pendingIntent = PendingIntent.getActivity(this, 0, intent, PendingIntent.FLAG_MUTABLE);
        notification = new Notification.Builder(this, "float_clipboard_channel_id")
                .setContentTitle("剪贴板悬浮窗服务")
                .setContentText("点击打开应用")
                .setSmallIcon(R.drawable.ic_launcher_foreground)
                .setContentIntent(pendingIntent)
                .build();

        // 启动前台服务
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.Q) {
            startForeground(1, notification, ServiceInfo.FOREGROUND_SERVICE_TYPE_SPECIAL_USE);
        } else {
            startForeground(1, notification);
        }
    }

    private void listenClipboard() {
        clipboardManager = (ClipboardManager) getSystemService(CLIPBOARD_SERVICE);
        clipboardManager.addPrimaryClipChangedListener(() -> {
            ClipData clip = clipboardManager.getPrimaryClip();
            if (clip != null && clip.getItemCount() > 0) {
                ClipData.Item item = clip.getItemAt(0);
                String text = item.getText().toString();
                Log.i(TAG, "listenClipboard: " + text);
//                updateTextView(text);
            }else{
                Log.i(TAG, "listenClipboard: is null");
            }
        });
    }

    private void updateTextView(String text) {
        if (floatView != null) {
            TextView textView = floatView.findViewById(R.id.float_text_view);
            textView.setText(text);
        }
    }

    private void getClibMessage() {
        ClipData clip = clipboardManager.getPrimaryClip();
        Log.i(TAG, "lsz getClibMessage: clip=" + clip);
        if (clip != null && clip.getItemCount() > 0) {
            ClipData.Item item = clip.getItemAt(0);
            String text = item.getText().toString();
            Log.i(TAG, " 11 lsz getClibMessage: " + text);
            sendToPC(text);
        //    updateTextView(text);
        }
        //   updateFocus(false);
    }

    private void sendToPC(String text) {

        new Handler().postDelayed(new Runnable() {
            @Override
            public void run() {
                String currentText = text;
                if (!currentText.equals(previousText)) {
                    Log.i(TAG, "lszz GoLog sendToPC: text:"+text );
                    Libp2p_clipboard.sendMessage(text);
                }
                previousText = currentText;
            }
        }, 100);


        /*Log.i(TAG, "sendToPC: text:" + text);
        new Thread(new Runnable() {
            @Override
            public void run() {
                Libp2p_clipboard.sendMessage(text);
            }
        }).start();*/
    }

    private void sendToPCIMG(byte[] value) {
        Log.i(TAG, "sendToPC: img byte:" + value);
        /*new Thread(new Runnable() {
            @Override
            public void run() {
                String base64String = Base64.encodeToString(value, Base64.DEFAULT);
                String  clearbase64String=removeInvalidCharacters(base64String);
                Libp2p_clipboard.sendImage(clearbase64String);
            }
        }).start();*/

        new Handler().postDelayed(new Runnable() {
            @Override
            public void run() {
                String currentimgmd5 ;
                MessageDigest md = null;
                try {
                    md = MessageDigest.getInstance("MD5");
                } catch (NoSuchAlgorithmException e) {
                    throw new RuntimeException(e);
                }
                byte[] md5Bytes = md.digest(value);

                StringBuilder sb = new StringBuilder();
                for (byte b : md5Bytes) {
                    sb.append(String.format("%02X", b));
                }
                currentimgmd5=sb.toString();
                if (!currentimgmd5.equals(previousImgMD5)) {
                    //Log.i(TAG, "lszz GoLog sendToPC: text:"+currentimgmd5 );
                    String base64String = Base64.encodeToString(value, Base64.DEFAULT);
                    String  clearbase64String=removeInvalidCharacters(base64String);
                    //Libp2p_clipboard.sendImage(clearbase64String);
                }

                previousImgMD5=currentimgmd5;
            }
        }, 100);

    }

    private void setClibMessageLoop() {
        new Handler().postDelayed(new Runnable() {
            @Override
            public void run() {
                //setClibMessage();

              //  getClibMessage();

                setClibMessageLoop();
                testClipboardUtils();

            }
        }, 3000);
    }

    private void setClibMessage() {
        ClipData clipData = ClipData.newPlainText(null, "编辑后的文本数据+" + testCount);
        clipboardManager.setPrimaryClip(clipData);
        testCount++;
    }

    @Override
    public IBinder onBind(Intent intent) {
        return null;
    }

    @Override
    public void onDestroy() {
        if (floatView != null) windowManager.removeView(floatView);
        if (notification != null) stopForeground(true);
        Log.d(TAG,"onDestroy");

        super.onDestroy();
    }

    private final static int P2P_EVENT_SERVER_CONNEDTED = 0;
    private final static int P2P_EVENT_SERVER_CONNECT_FAIL = 1;
    private final static int P2P_EVENT_CLIENT_CONNEDTED = 2;
    private final static int P2P_EVENT_CLIENT_CONNECT_FAIL = 3;

    private Callback getGolangCallBack() {
        return new Callback() {
            @Override
            public void callbackMethod(String s) {
                Log.i(TAG, "lsz GoLog callmsg callbackMethod: callback调用 ==" + s);
                //ClipData clipData = ClipData.newPlainText(null, s);
                //clipboardManager.setPrimaryClip(clipData);
                /*ClipboardUtils clipboardUtils = ClipboardUtils.getInstance();

                AtomicReference<ClipData> clipDataRef = ClipboardUtils.createClipdataRef();
                clipboardUtils.addTextItem(clipDataRef, s);

                clipboardUtils.setPrimaryClip(clipDataRef);*/

               // setClipToClipboard(s,null);
            }

            @Override
            public void callbackMethodFileConfirm(String s, long l) {

            }

            @Override
            public void logMessageCallback(String msg) {
                Log.i(TAG, "logMessageCallback: msg: " + msg);
            }

            @Override
            public void callbackMethodImage(String msg) {
                Log.i(TAG, "lsz callmsg GoLog callbackMethodImage: msg: " + msg.toString());
                //setClipToClipboard(null,msg);

                getImg(msg);
            }

            @Override
            public void eventCallback(long event) {
                Log.i(TAG, "eventCallBack: event2: " + event);
            }
        };
    }

    public void getImg(String msg){
        byte[] value =Base64.decode(msg, Base64.DEFAULT);
        Log.d("lszzz", "GoLog sendToPC: img value:==" + value);
        Bitmap bitmap = BitmapFactory.decodeByteArray(value, 0, value.length);
        imageView.setImageBitmap(bitmap);
    }
    @Override
    public boolean onKey(View view, int i, KeyEvent keyEvent) {
        Log.i(TAG, "keyEvent: " + i);
        if (keyEvent.getKeyCode() == keyEvent.KEYCODE_BACK) {
            Log.i(TAG, "back");
        }
        return false;
    }

    public static String removeInvalidCharacters(String base64String) {
        // 正则表达式，匹配Base64的有效字符
        String regex = "[^A-Za-z0-9+/=]";
        // 使用正则表达式替换掉非法字符
        String cleanString = base64String.replaceAll(regex, "");
        return cleanString;
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
                    imageView2.setImageURI(imageUri);
                    //imageView2.setImageBitmap(bitmap1);
                }
            }
        } else {        Toast.makeText(this, "Clipboard is empty", Toast.LENGTH_SHORT).show();    } */
//本地图片 取到剪切版end

        AtomicReference<ClipData> clipDataRef = new AtomicReference<>(null);
        ClipboardUtils clipboardUtils = ClipboardUtils.getInstance();
        clipboardUtils.getPrimaryClip(clipDataRef);
        //Log.e("clip", "lsz len===hasClip GoLog=clipboardUtils.hasClip()=" + clipboardUtils.hasClip());
        for (int i = 0; i < clipboardUtils.getItemCount(clipDataRef); i++) {
            Log.e("clip", "GoLog lsz len=getItemType 0：TETX 1：IMG=" + clipboardUtils.getItemType(clipDataRef, i));
            if (clipboardUtils.getItemType(clipDataRef, i) == clipboardUtils.CLIPBOARD_DATA_TYPE_TEXT) {
                String text = clipboardUtils.getTextItem(clipDataRef, i);
                TextView textView = floatView.findViewById(R.id.text_id);
                textView.setText("get="+text);
                sendToPC(text);
                //Log.e("clip", "lsz len===hasClip text GoLog=" + text);
            } else if (clipboardUtils.getItemType(clipDataRef, i) == clipboardUtils.CLIPBOARD_DATA_TYPE_IMAGE) {
                Bitmap bitmap1 = clipboardUtils.getImageItem(clipDataRef, i);
                Log.e("clip", "lsz GoLog len===hasClip bitmap1==" + bitmap1);
                if(bitmap1 != null) {
                /*
                数组转bitmap
                */
                    //Bitmap drawableicon = BitmapFactory.decodeResource(getResources(), R.drawable.liu2);
                    //byte[] imageData = bitmapToByteArray(drawableicon); // 要转换的字节数组
                    //Bitmap bitmap3 = BitmapFactory.decodeByteArray(imageData, 0, imageData.length);
                    //imageView2.setImageBitmap(bitmap1);

                    //bitmap转byteArray
                    //int bytes = bitmap1.getByteCount();
                    //ByteBuffer buf = ByteBuffer.allocate(bytes);
                    //bitmap1.copyPixelsToBuffer(buf);
                    //byte[] byteArray = buf.array();

                    imageView = floatView.findViewById(R.id.imageView);
                    imageView.setImageBitmap(bitmap1);

                    byte[] imageData = bitmapToByteArray(bitmap1);

                    sendToPCIMG(imageData);


                }else{
                    Toast.makeText(this, " Clipboard img is empty", Toast.LENGTH_SHORT).show();
                }


            } else {
                Log.e("clip", "not support format");
                Toast.makeText(this, "lsz111 Clipboard is empty", Toast.LENGTH_SHORT).show();
            }
        }


    }


    private void setClipToClipboard(String string, String msg) {
        ClipboardUtils clipboardUtils = ClipboardUtils.getInstance();
        //clipboardUtils.clearClip();
        AtomicReference<ClipData> clipDataRef = ClipboardUtils.createClipdataRef();
        if(!string.equals("") || string != null){
            clipboardUtils.addTextItem(clipDataRef, "test text12");
        }

        if(!msg.equals("") || msg != null){

            byte[] decodedBytes = Base64.decode(msg, Base64.DEFAULT);
            Bitmap bitmap = BitmapFactory.decodeByteArray(decodedBytes, 0, decodedBytes.length);

            clipboardUtils.addImageItem(clipDataRef, bitmap);
        }

        clipboardUtils.setPrimaryClip(clipDataRef);

    }


    public static byte[] bitmapToByteArray(Bitmap bitmap) {
        ByteArrayOutputStream outputStream = new ByteArrayOutputStream();
        bitmap.compress(Bitmap.CompressFormat.JPEG, 100, outputStream);
        Log.d("lsz", "outputStream. imag toByteArray()=" + outputStream.toByteArray().toString());
        return outputStream.toByteArray();

    }

    }




