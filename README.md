
Basic deployment:

1) Set Up Logging

On the remove server, create file:

```
/var/log/ZombieDice/logfile.txt
```

You probably will need `sudo`

Change the file permissions to `666`

2) Run the `deploy_scrips/setup.py` file.

It takes 3 argumens: ip address, username@ip address, locatation on remote

For example:

```
bash setup.py 12.23.34.45 ubuntu@12.23.34.45 "~/" 
```

This will deploy and run the server.


