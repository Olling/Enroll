# Enroll
Enroll allows easy enrollment of Linux systems into [Enrolld](https://github.com/Olling/Enrolld)


## Configuration Examples

### Main configuraton file (default location: /etc/enroll/enroll.conf)
``` 
{
  "URL": "https://EnrolldServer:8304",
  "ConfigFragments": "/etc/enroll/conf.d",
  "Payload": {
    "ServerID": "fqdn.example.com",
    "Inventories": ["Copenhagen", "WebServer"],
    "AnsibleProperties": {
      "Environment": "production"
    }
  }
}
```

### Configuraton fragment (default location: /etc/enroll/conf.d/*)
``` 
{
  "Inventories": ["Customer_1"],
  "AnsibleProperties": {
    "CustomerUsage": "Training"
  }
}
```

## Examples

``` 
$ enroll --status
Status: Not enrolled
```
``` 
$ enroll --enroll
Status: Enrolled
```
