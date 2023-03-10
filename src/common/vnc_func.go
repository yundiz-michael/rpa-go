package common

/*
   extern void onMouseCallback(char* taskName,int button,int x,int y);
   extern void onKeyboardCallback(char* taskName,int down,int key);
   extern void onConnectionCallback(char* taskName,int clientCount);

   void mouseCallbackEvent(char* taskName,int button,int x,int y) {
     onMouseCallback(taskName,button,x,y);
   }

   void keyboardCallbackEvent(char* taskName,int down,int key) {
     onKeyboardCallback(taskName,down,key);
   }

   void connectionCallbackEvent(char* taskName,int clientCount) {
     onConnectionCallback(taskName,clientCount);
   }

*/
import "C"
