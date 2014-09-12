package com.wustlhack.messengerdemo.app;

import android.app.Activity;
import android.content.Context;
import android.content.Intent;
import android.os.Bundle;
import android.support.v4.content.LocalBroadcastManager;
import android.support.v4.content.WakefulBroadcastReceiver;
import android.util.Log;

import com.google.android.gms.gcm.GoogleCloudMessaging;

// GcmBroadcastReceiver receives the broadcast from Google Play
// Services and forwards it to a local broadcast manager.
public class GcmBroadcastReceiver extends WakefulBroadcastReceiver {
    @Override
    public void onReceive(Context context, Intent intent) {
        Bundle extras = intent.getExtras();
        Intent RTReturn = new Intent(MainActivity.RECEIVE_MESSAGE);
        RTReturn.putExtra("message", extras.getString("message"));
        LocalBroadcastManager.getInstance(context).sendBroadcast(RTReturn);
        setResultCode(Activity.RESULT_OK);
    }
}