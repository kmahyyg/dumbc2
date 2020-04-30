# dumbyc2 - Command line arguments

Dumbyc2, yet another dumb c2 tool, both server and agent.

Communication protected by TLS.

Run this program with command line arguments:

For Agent:

- `-r <ADDR:PORT>` Remote Controller IP

For Controller:

- `-l <ADDR:PORT>` Listen to.

General:

- `-C Directory`, optional, certificate location, directory should be absolute path if possible, 
else no guarantee the program will provided to find the file correctly, default to `~`.

If you are running the first time, you should run `certgen` first.

- `-h, --help`, optional, show help message.  

# Certgen - Command line arguments

Certgen will generate certificates and corresponding keys to `<Folder You Defined / Home Dir>/.dumbyc2/`.

Generate RSA 4096 bits certificates (both CA and Server) and save corresponding keys.

- `-o`, optional, output directory, default to `~`.

