# Go LTI

Public repository for Canvas LTI 1.3 implementation in Go.

## Generating a new key pair

```bash
openssl genrsa -out keys/private.pem 2048
openssl rsa -in keys/private.pem -pubout -out keys/public.pem
```

## Useful links

- [Canvas LTI 1.3 Documentation](https://documentation.instructure.com/doc/api/file.tools_intro.html)
