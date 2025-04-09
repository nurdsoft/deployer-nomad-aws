# nomad-deploy-trigger

This contains all logic necessary to deploy a lambda function in AWS attached to an 
API Gateway to trigger nomad jobs from the internet while keeping the nomad 
infrastructure private.  This is particularly useful in the case of CI/CD with
Github, Gitlab etc.
