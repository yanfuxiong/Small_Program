package com.rtk.myapplication;

import android.app.Application;
import android.util.Log;

import com.tencent.mmkv.MMKV;

public class MyApplication extends Application {

    @Override
    public void onCreate() {
        super.onCreate();

        String rootDir = MMKV.initialize(this);
        Log.i("MMKV", "mmkv root: " + rootDir);
    }
}
