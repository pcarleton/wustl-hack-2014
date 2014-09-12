package com.wustlhack.messengerdemo.app;

import android.app.Activity;
import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.content.IntentFilter;
import android.support.v4.content.LocalBroadcastManager;


import android.os.AsyncTask;
import android.os.Bundle;
import android.util.Log;
import android.widget.TextView;

import android.widget.Toast;

import com.google.android.gms.gcm.GoogleCloudMessaging;


import java.io.IOException;


public class MainActivity extends Activity {
    public static final String RECEIVE_MESSAGE = "com.wustlhack.messengerdemo.RECEIVE_MESSAGE";

    String PROJECT_NUMBER = "56902053111";

    private GoogleCloudMessaging gcm;
    private String regid;

    private TextView messageView;


    // Set up a broadcast receiver
    private BroadcastReceiver bReceiver = new BroadcastReceiver() {
        @Override
        public void onReceive(Context context, Intent intent) {
            if(intent.getAction().equals(RECEIVE_MESSAGE)) {
                String message = intent.getStringExtra("message");
                messageView.setText(message);
            }
        }
    };

    // registerWithGCM contacts the GCM server and logs the ID it receives.
    public void registerWithGCM(){
        new AsyncTask<Void, Void, String>() {
            @Override
            protected String doInBackground(Void... params) {
                String msg = "";
                try {
                    if (gcm == null) {
                        gcm = GoogleCloudMessaging.getInstance(getApplicationContext());
                    }
                    String regid = gcm.register(PROJECT_NUMBER);
                    // In reality, we would want to send this to the server so it can reach us.
                    // For the sake of simplicity for this demo, we'll just copy and paste it from
                    // the logs.
                    msg = "Device registered, registration ID=" + regid;
                    Log.i("GCM", msg);
                } catch (IOException ex) {
                    msg = "Error :" + ex.getMessage();
                }
                return msg;
            }
        }.execute(null, null, null);
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        messageView = (TextView) findViewById(R.id.messageView);

	// Tell GCM we want to receive messages.
        registerWithGCM();

	// Tell our local broadcast manager we want to receive messages and handle
	// them with our 'bReceiver'.
        LocalBroadcastManager bManager = LocalBroadcastManager.getInstance(this);
        IntentFilter intentFilter = new IntentFilter();
        intentFilter.addAction(RECEIVE_MESSAGE);
        bManager.registerReceiver(bReceiver, intentFilter);
    }
}
