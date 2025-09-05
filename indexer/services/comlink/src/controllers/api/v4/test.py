import os
import binascii

# Generate a 32-byte (256-bit) private key
private_key = binascii.hexlify(os.urandom(32)).decode()



print(f"Generated private key: {private_key}")
