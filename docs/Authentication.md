# eaglechat Authentication and Trust Model

This document outlines the security model for identity, authentication, and message exchange within the eaglechat network.

## 1. Core Identity

- **Identity Provider:** Each client application is its own identity provider.
- **Key Pair:** Upon first launch, each client MUST generate a public/private RSA key pair.
  - The **Public Key** serves as the user's unique, verifiable, network-wide identity.
  - The **Private Key** is the user's proof of identity and MUST remain confidential on the client's device.

## 2. Service Roles

- **Clients:** Responsible for all cryptographic operations. They generate keys, sign messages, and perform end-to-end encryption.
- **ID Managers:** Act as a simple **Public Key Directory**. Their role is to map usernames to their last known public key and IP address. ID Managers have no special cryptographic authority and do not share a common key pair.

## 3. Message Security

Communication in the network relies on two distinct cryptographic processes:

### Digital Signatures (Authenticity & Integrity)

To prove a message was sent by a specific user and was not altered in transit, all messages are digitally signed.

1.  The sender (Client A) calculates a hash of the message content.
2.  Client A encrypts this hash with their **own private key**. This encrypted hash is the signature.
3.  The signature is attached to the message.
4.  The receiver (Client B) uses Client A's **public key** to decrypt the signature, revealing the original hash.
5.  Client B calculates its own hash of the received message and compares it to the decrypted hash. If they match, the message is authentic.

### Hybrid Encryption (Confidentiality)

To ensure only the intended recipient can read a message while maintaining high performance, all messages are encrypted using a **hybrid cryptosystem** (RSA + AES). This process is abstracted into a `SecureEnvelope` that contains the final encrypted payload.

1.  **AES Key Generation:** The sender (Client A) generates a new, single-use 256-bit symmetric AES key.
2.  **Message Encryption:** The plaintext message is encrypted using this fast AES key.
3.  **Key Wrapping:** The single-use AES key is then encrypted (or "wrapped") using the recipient's (Client B's) RSA public key.
4.  **Final Assembly:** The encrypted AES key and the message ciphertext are bundled together. A digital signature is created over these two components to ensure authenticity.
5.  **Transmission:** The resulting `SecureEnvelope`, containing the wrapped key, the ciphertext, and the signature, is sent over the network.

Only Client B, using their private RSA key, can decrypt the AES key. Once the AES key is revealed, it can be used to quickly decrypt the actual message.

## 4. Trust Model: Trust On First Use (TOFU)

To prevent Man-in-the-Middle (MITM) attacks where an attacker might substitute a user's public key, the network will adopt the **Trust On First Use (TOFU)** model.

### Workflow

1.  When Client A wants to communicate with User B for the first time, it queries an ID Manager for User B's public key.
2.  Upon receiving the key, Client A **MUST** cache it locally and create a permanent association between User B's username and that specific public key.

### Security Guarantees & Risks

- **Security:** For all subsequent communications, Client A will exclusively use the cached public key for User B. This prevents an attacker from later impersonating User B, as their public key will not match the cached key.
- **Vulnerability:** This model's primary vulnerability is a potential MITM attack during the **very first key exchange**. The system trusts that the first key received for a user is the correct one.
- **Mitigation:** If a client ever queries an ID Manager and receives a *different* public key for a user it already has in its cache, the client application **MUST** treat this as a critical security event. It should alert the user that the remote identity has changed and that proceeding may be unsafe.

## 5. Key Management

- **Recovery:** This design does not currently specify a key recovery mechanism. If a user loses their private key (e.g., by losing their device), their identity is considered **unrecoverable**. They will need to generate a new key pair and establish a new identity on the network.
