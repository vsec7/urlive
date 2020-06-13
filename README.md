# urlive

Check url is live (*HTTP status code "200 ok" only*).

Skiping for dead link :P

## Install
```
 ▶ go get -u github.com/vsec7/urlive
```

## Basic Usage

```
 ▶ urlive --help 

Urlive (Check url is live *HTTP status code "200 ok" only)

By : viloid [Sec7or - Surabaya Hacker Link]

Basic Usage :
 ▶ echo http://domain.com/path/file.ext?param=value | urlive
 ▶ cat listurls.txt | urlive -c 50"

Options :
  -H, --header        Add Header to the request
  -c, --concurrency   Increase concurrency level (*default 20)
  -x, --proxy         Add HTTP proxy
  -m, --match         Add match specific string (*Sensitive Case)
  -o, --output        Output to file

```

## Advanced Usage
```
With Increase concurrency level + additional header + match specific string + output to file
e.g : 
 ▶ cat listurls.txt | urlive -c 50 -H "Authorization: Bearer eXAmPLe" -H "Cookie: ExAmpLe" -m "b374k" -o urlive.txt

** Output : urlive.txt all url with http status code 200ok that contains 'b374k' in response body**
^warn! use flag -m/-match <string> could decrease speed bcs should read response body to match string

```

## Credit And Thanks
```
@tomnomnom (github.com/tomnomnom)
```