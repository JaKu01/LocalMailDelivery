# Local Mail Delivery with REST API
This project is a simple implementation of a local mail delivery system offering a REST API. 
It can be self-hosted to deliver mails to a local mailbox which 
can be accessed by any IMAP client inside the respective LAN.

## Usage
The project is currently in its early stages and is not yet suitable for production use. 
To run the application, it is recommended to use Docker with the provided Dockerfile. 
You can specify your username and password as the first and second arguments of the `CMD` instruction in the Dockerfile. 
Please note that this setup is intended for testing purposes only, as it is not secure. 
This approach will be updated in the future.


## Planned Features
- [ ] Deliver mails to IMAP clients inside the LAN
- [ ] REST API with authentication for sending mails
- [ ] Persistence of mails

## Architecture
![Architecture](./assets/architecture.svg)

