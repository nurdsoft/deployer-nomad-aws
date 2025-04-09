# nomad-deploy-trigger

This contains all logic necessary to deploy a lambda function in AWS attached to an 
API Gateway to trigger nomad jobs from the internet while keeping the nomad 
infrastructure private.  This is particularly useful in the case of CI/CD with
Github, Gitlab etc.

# Usage 

```shell
$ make deploy
```

This will run through the complete pipeline and

- compile the go binary
- package it up
- deploy it to AWS lambda with the appropriate configuration
