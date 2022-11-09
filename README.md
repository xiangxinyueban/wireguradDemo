# p2p-wiregurad
Wireguard Server  
Contains three components  
1. Central Server (backend for user mgmt)
2. BootStrap Server (libp2p BootStrap Server)
3. Wireguard Server (Wireguard Relay)

Central Server supports User register, user login, user activation with activation code, generate activation code
. Central Server also supports server wireguard server management. You can add new wireguard server just provide username and password to get server online. And Central Server will monitor the status of wireguard server with heartbeat. If server offline. You can remove the server.  

And the central server can monitor the status of users.
❗ For learning purposes only, please do not use for illegal purposes.