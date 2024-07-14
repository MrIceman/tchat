tchat is a terminal-based chat application.
No fancy GUI. You can use the terminal (e.g. from your IDE) to connect to a tchat server and start talking immediately.

## architecture
I'm using some sort of a hexagonic architecture for the client and the server.
* server -> all the server related  code, you'll find here the server (duh), the handler, domain and data/network layers. The domain will include the business logic and data will take care of connection management, writing/reading to storages
* client -> analog to server, the client contains mostly of command parsing

## outlook
Though it's still WIP, here are some ideas that I'll implement.
* centralized vs p2p mode (support both functionalities)
* highly configurable server (authentication, tls, roles, anonymous, encryption, custom storages, notifications, allow pipelining, e.g. "cat somefile.txt | tchat --channel somechannel")
  
