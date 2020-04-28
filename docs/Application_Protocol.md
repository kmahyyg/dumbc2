# Commands

`UPLD`: Upload

- Send Command
- Send Length
- Send Path
- Send Data

`DWLD`: Download

- Send Command
- Send remote path
- Check Path writable and remote file exists
- Recv Length
- Recv Data

`BOOM`: Self-Destroy

- Send Command

`BASH`: Interactive shell

- Send Command
- Check Status
- Connect to shell

`INJE`: Shellcode execution

- Send Command
- Get Pingback

# Responses

`00$xx$DATA`

- `00` status code, 00 for success, 01 for failed, 02 for EEXIST
- `xx` file length, in megabytes, can be zero if not required.
- `DATA` for details, if no data, just use `\x00`.