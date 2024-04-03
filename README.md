# grand - CLI tool for generating cryptographically-secure random byte strings

## Installation

```
go install github.com/soroushj/grand@latest
```

## Usage

```
grand [flags]
```

The flags are:

```
-e encoding
    output encoding; one of:
      "hex"   - base16
      "b64s"  - base64, standard alphabet
      "b64sr" - base64, standard alphabet, no padding
      "b64u"  - base64, URL-safe alphabet
      "b64ur" - base64, URL-safe alphabet, no padding
      "b32s"  - base32, standard alphabet
      "b32sr" - base32, standard alphabet, no padding
      "b32h"  - base32, extended hex alphabet
      "b32hr" - base32, extended hex alphabet, no padding
      (default "hex")
-n int
    number of random byte strings to generate (default 1)
-s size
    size of random byte strings; an integer or an inclusive range,
    e.g. "16-32" (if a range is specified, the size of each byte
    string will be a cryptographically-secure random number in the
    range) (default "16")
```

An example for generating two 32-byte base64url-encoded random byte strings:

```
$ grand -e b64u -s 32 -n 2
8cYgWnTqRq1RkDJBB-PKnf1pp7svVox_c0bR8lqctIM=
5XzQ2In5cO2Rz5GD8VcHwJJZFS7iCXydbdEbnyVLfJE=
```
