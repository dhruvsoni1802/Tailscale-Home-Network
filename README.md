# Tailscale Home Network
This is a simple project that allows you to share files between your devices using Tailscale.

## Usage for Storage Node server
Whenever you are executing the server binary on the Mac or Windows machine, you need to set the TS_AUTHKEY environment variable. You can get the auth key from the Tailscale dashboard.

To set the TS_AUTHKEY environment variable, you can run the following command on Mac:
```
export TS_AUTHKEY=<your-auth-key>
```

On Windows, you can set the TS_AUTHKEY environment variable by running the following command:
```
set TS_AUTHKEY=<your-auth-key>
```

Once you have set the TS_AUTHKEY environment variable, you can execute the server binary.


## Sample commands from client

```
curl -X POST -F "file=@/path/to/file" http://storage-node:8080/upload
```

```
curl http://storage-node:8080/health
```

