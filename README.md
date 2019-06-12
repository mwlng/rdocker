## Remote docker cli 


### Overview

rdocker is a remote docker command line. You can use it to send the Docker commands or instructions to the Docker daemon which is running on the remote hosts.

### Examples

1. Display running conaianers and their ids etc. 
 ```  
rdocker -i <ssh keyfile> -H <hostname or ip> -- ps -a
```

2. Log into the conatainer with it's name/id.
```
rdocker -i <ssh keyfile> -H <hostname or ip> -- exec -it <container name or id> /bin/bash