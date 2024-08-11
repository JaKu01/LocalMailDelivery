# Local Mail Delivery with REST API
This project is a simple implementation of a local mail delivery system offering a REST API. 
It can be self-hosted to deliver mails to a local mailbox which 
can be accessed by any IMAP client inside the respective LAN.

## Usage
The project is currently in its early stages and is not yet ready for production use. 
It is recommended to use Docker with the provided Dockerfile to run the application. 
You can configure your username, password, certificates, and keys through the following environment variables:

```
USERNAME=yourusername
PASSWORD=yourpassword
CERTIFICATE_PATH=/path/to/certificate
KEY_PATH=/path/to/key
```

## Planned Features
- [ ] Deliver mails to IMAP clients inside the LAN
- [ ] REST API with authentication for sending mails
- [ ] Persistence of mails

## Architecture
![Architecture](./assets/architecture.svg)

