## Remote docker cli 

### Overview
rdocker is a remote docker command line. You can use it to send the Docker commands or instructions to the Docker daemon which is running on the remote hosts.

## Prerequisite
It requires Docker client to be installed on your computer.

### Examples
1. Display running conaianers and their ids etc. 
 ```  
rdocker -u <ssh user> -i <ssh keyfile> -H <hostname or ip> -- ps -a
```

2. Log into the conatainer with it's name/id.
```
rdocker -u <ssh user> -i <ssh keyfile> -H <hostname or ip> -- exec -it <container name or id> /bin/bash