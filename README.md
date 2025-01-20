# pasted

`pasted` is a server to receive and store data passed over a socket connection, similar to <https://termbin.com>.

**Note: This project is under active development and will likely break at some point**

## Usage

To use `pasted`, you can pipe data to it over a socket connection. For example, to paste the contents of a file:

```sh
cat file.txt | nc pasted.example.com 9999
```

`pasted` will respond with a URL where the data can be accessed.

## Installation

To install `pasted`, you can use the provided `compose.yaml` file. This will start up pasted using sqlite as the backend.:

```sh
docker compose up -d
```

To run it locally, you can use the following command:

```sh
go run main.go --config test.yaml
```

## Supported Backends

`pasted` supports the following backends for storing files:

- `file`: Stores files on disk (default).
- `memory`: Stores files in memory. Useful for testing.
- `pgx`: Stores files in a PostgreSQL database.
- `redis`: Stores files in a Redis database.
- `sqlite`: Stores files in a SQLite database.
- `s3`: Stores files in an S3 bucket.

### TODO

- [ ] `azure`: Stores files in an Azure Blob Storage container.
- [ ] `gcs`: Stores files in a Google Cloud Storage bucket.
- [ ] `(s)ftp`: Stores files on an FTP or SFTP server.
- [ ] `smb`: Stores files on an SMB share.
- [ ] `nfsv4`: Stores files on an NFSv4 share.
- [ ] `ipfs`: Stores files on IPFS.


## Transforms

`pasted` supports the following transforms for modifying data before storing it:

- `aes`: Encrypts data using AES-256-GCM.
- `gzip`: Compresses data using Gzip.
- `base64`: Encodes data using Base64.

## Contributing

To contribute to `pasted`, please fork the repository and submit a pull request. You can also submit issues or feature requests.

## License

`pasted` is licensed under the MIT license. See [LICENSE](LICENSE) for more information.
