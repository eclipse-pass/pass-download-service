# PASS manuscript download service

[![Build Status](https://travis-ci.com/oa-pass/pass-download-service.svg?branch=master)](https://travis-ci.com/oa-pass/pass-download-service)

Contains a partial impmementation of the PASS download service

## Usage

    pass-download-service serve

## API

The implementation has a simple provisional API

### Lookup DOI
Looks up a DOI and returns JSON containing available author-accepted manuscripts

```
GET http://<HOSTNAME>:<PORT>/lookup?doi=<DOI>
```

Example: `curl [http://localhost:6502/download?doi=10.1038%2Fnature12373](http://localhost:6502/download?doi=10.1038%2Fnature12373)`

Returns:
```
{
  "manuscripts": [
    {
      "description": "oa repository (via OAI-PMH doi match)",
      "location": "http://europepmc.org/articles/pmc4221854?pdf=render"
    }
  ]
}
```

### Download DOI
Given a DOI and a manuscript URL (from a previous lookup), will download the manuscript at the given URL into Fedora, and
return the URL of the Fedora object containing the downloaded binary.  Its up to the client to later on create a PASS `File` entity that
points to the resulting Fedora URL as content.

If the URL does not match any URLs from the corresponding lookup query, the request will fail with a "bad request" error code.

The response body and `Location` header will contain the Fedora binary URL

POST with an empty body:
```
POST  http://<HOSTNAME>:<PORT>/lookup?doi=<DOI>&url=<URL>
```

Example:
`curl -X POST [http://localhost:6502/download?doi=10.1038%2Fnature12373&url=http%3A%2F%2Feuropepmc.org%2Farticles%2Fpmc4221854%3Fpdf%3Drender](http://localhost:6502/download?doi=10.1038%2Fnature12373&url=http%3A%2F%2Feuropepmc.org%2Farticles%2Fpmc4221854%3Fpdf%3Drender)`

Result:
``
http://localhost:8080/fcrepo/rest/files/b3/b6/e7/e6/b3b6e7e6-57e0-47e0-b6b1-5f7271f3c76a
``

## Configuration

For cli flags, see `pass-download-service help`

Environment variables are as follows:

* `DOWNLOAD_SERVICE_PORT` - Port to serve the download service on (default `6502`)
* `DOWNLOAD_SERVICE_DEST` - Fedora container URI where binaries will be downloaded into
* `UNPAYWALL_REQUEST_EMAIL` - E-mail address that will be sent with unpaywall requests
* `UNPAYWALL_BASEURI` - BaseURL of the unpaywall service.
* `PASS_EXTERNAL_FEDORA_BASEURL` - Public facing PASS Fedora Baseurl
* `PASS_FEDORA_BASEURL` - Internal Fedora Baseurl
* `$PASS_FEDORA_USER` - Fedora username
* `$PASS_FEDORA_PASSWORD` - Fedora password