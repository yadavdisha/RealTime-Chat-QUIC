Steps to Execute the Project:-

Step 01:- Firstly fork the repo into any local vscode environment into a directory which has GO alrealdy setup and open it from there in the VSCode.
Step 02:- Open a New Terminal and start the server using:-
      go run cmd/echo/echo.go -server

This will start the server successfully.

Step 02:- Then open another Terminal and start the client using:-
go run cmd/echo/echo.go -client -username "username" (Mention any username you want)

This is start the client successfully.

Step 03:- Once the intial handshake happens and a user join in client and the same message is recieved by the server the communication is established.

Step 04:- Now you can open multiple clients simultaneously and observe the communication between all the clients.

Step 05:- Once you kill a particular terminal , in all the chat windows a message would be displayed saying User has left the chat and the same happens when a particular user joins the chat by starting a new client session.

Step 06:- We have an emoji support that can be used while typing messages and windows shortcut ( Windows +.) so when you type a message and want to insert an emoji along with it you can use this shortcut and the communication can happen like that way as well.

Step 07:- We also have a feature to list all the users currently connected or active so in any of the chat window if you type **/list** you will be able to view all the connected users currently to the Student Colaboration Chat Application Protocol.

![image](https://github.com/rohitaragde/Collaborative-Chat-Application-QUIC/assets/32512875/b93d2e19-617d-4008-aa27-eb98ffcb9280)

Step 08:- All the clients would be able to see the messages of all the other clients in the chat window as well as will be notified when someone joined or left the chat.

Step 09:- On server we get all the clients and the pdu information along with the communication as well.

Step 10:- There will be timeouts in the server so you might have to kill all the sessions and re-start the server as well as the client sessions.




