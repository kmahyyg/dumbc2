# dumbyc2 - Command line arguments

Dumbyc2, yet another dumb c2 tool, both server and agent.

Communication protected by TLS.

Run this program with command line arguments:

-  `-s`, run as controller.
-  `-c`, run as agent.

`-s` and `-c` are conflict, you must choose one.

- `-b`, optional, run as bind agent, will try to bind port, must co-operate with `-c`.

if not `-b` is not defined, it works as reverse connection.

- `-H IP`, required, define the remote IP or bind IP.
- `-P Port`, required, define the remote Port or bind Port.

- `-C Directory`, optional, certificate location, directory should be absolute path if possible, 
else no guarantee the program will provided to find the file correctly, default to `~`.

If you are running the first time, you should run `certgen` first.

Certgen will generate certificates and corresponding keys to `<Folder You Defined / Home Dir>/.dumbyc2/`.

- `-h, --help`, optional, show help message.  

# Certgen - Command line arguments

Generate RSA 4096 bits certificates (both CA and Server) and save corresponding keys.

- `-o`, optional, output directory, default to `~`.

