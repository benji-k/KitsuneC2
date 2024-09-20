This folder contains the SSL certificates that will be used to encrypt HTTP traffic from client<->frontend. The installation
script will generate self-signed certificates. 

However, you can put your own certificates here. To do so, export your certificates to .key and .pem formats, rename them to kc2SSL.key and kc2SSL.pem
and place them in this folder. Afterwards, restart the server. 
