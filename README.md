# StegoSecretS

StegoSecretS empowers the SSS with steganography

```
stego encrypt -f/--file file.txt -p/--parts -t/--threshold -o/--output
# if no file a message can be input and embedded inside the images

stego decrypt -f/--file file.aes --master-key file.key -k/--key -i/--img
# if no file is provided the it will get the message from the images/keys
```


```
# out/file_20220102_230313
file
file.checksum // clear data checksum
file.aes // encrypted data
file.aes.checksum // encrypted data checksum
file.key // Base64 encoded master key
// Base64 Partial keys
file.1.key
file.2.key
file.3.key
image.1.jpg
image.1.jpg.checksum
image.2.jpg.checksum
image.3.jpg.checksum
```


cat 0.jpg.checksum | sha256sum --check