# pwnpaste
Pwnpaste queries haveibeenpwned.com for compromised emails. Pastes respective to the compromised emails are outputted to stdout or exported to stdout in HTML format. Pwnpaste is meant to act as a quick way to identify compromises and gather information useful to assessments.
##Authors
* [Jonathan Broche](https://github.com/gojhonny)
* [Tom Steele](https://github.com/tomsteele)


## Download and Install
-----
Download the latest compiled release for your OS [here](https://github.com/gojhonny/pwnpaste/releases/latest).

```
tar -zxvf pwnpaste*.tar.gz
cd pwnpaste*amd64
./pwnpaste -h
```



## Help
-----
```
Usage of ./pwnpaste:
  -c int
    	Maximum number of concurrent requests (default 10)
  -html
    	Output HTML (has webcache links)
  -i string
    	File containing new line delimited email addresses
  -v	Print version and exit
```

## Examples
-----

###Text Table Output
```
./pwnpaste foo@bar.com
  `( ◔ ౪◔)´  I'm done ...0

pwnpaste 1.0.0

+-------------+-----------------------------------------+-------------------------------------------------------------------------------------+
|    EMAIL    |              PASTEBIN URL               |                                   CACHEDVIEW URL                                    |
+-------------+-----------------------------------------+-------------------------------------------------------------------------------------+
| foo@bar.com | https://pastebin.com/raw.php?i=vwXXWCEN | https://webcache.googleusercontent.com/search?q=cache:https://pastebin.com/vwXXWCEN |
|             | https://pastebin.com/raw.php?i=qAegkpzu | https://webcache.googleusercontent.com/search?q=cache:https://pastebin.com/qAegkpzu |
|             | https://pastebin.com/raw.php?i=Lg80iL8k | https://webcache.googleusercontent.com/search?q=cache:https://pastebin.com/Lg80iL8k |
|             | https://pastebin.com/raw.php?i=uMq1W2mx | https://webcache.googleusercontent.com/search?q=cache:https://pastebin.com/uMq1W2mx |
|             | https://pastebin.com/raw.php?i=9ZKSRx5i | https://webcache.googleusercontent.com/search?q=cache:https://pastebin.com/9ZKSRx5i |
|             | https://pastebin.com/raw.php?i=L6fZS5VC | https://webcache.googleusercontent.com/search?q=cache:https://pastebin.com/L6fZS5VC |
|             | https://pastebin.com/raw.php?i=b6taeWri | https://webcache.googleusercontent.com/search?q=cache:https://pastebin.com/b6taeWri |
|             | https://pastebin.com/raw.php?i=ba6LmF9Z | https://webcache.googleusercontent.com/search?q=cache:https://pastebin.com/ba6LmF9Z |
|             | https://pastebin.com/raw.php?i=wXb5W8GV | https://webcache.googleusercontent.com/search?q=cache:https://pastebin.com/wXb5W8GV |
|             | https://pastebin.com/raw.php?i=EE8GM0ed | https://webcache.googleusercontent.com/search?q=cache:https://pastebin.com/EE8GM0ed |
|             | https://pastebin.com/raw.php?i=8Q0BvKD8 | https://webcache.googleusercontent.com/search?q=cache:https://pastebin.com/8Q0BvKD8 |
|             | https://pastebin.com/raw.php?i=C4GdBDnP | https://webcache.googleusercontent.com/search?q=cache:https://pastebin.com/C4GdBDnP |
+-------------+-----------------------------------------+-------------------------------------------------------------------------------------+
```
###HTML Output
```
./pwnpaste -i emails.txt -html > foo.html
```
