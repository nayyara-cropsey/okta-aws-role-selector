# Okta AWS Role Selector

Middleware app for selecting a AWS role from SAML assertion.

Okta AWS Role Selector is an app for extracting AWS roles from a SAML assertion similar to the AWS SAML console: https://signin.aws.amazon.com/saml.

## Usecase
Okta AWS Role Selector is used to co-ordinate role selection for an app that authenticate users using an AWS SAML assertion (setup with Okta).
Such an app might use the SAML assertion with STS and **AssumeWithSAML** to obtain a temporary AWS STS token for internal use.

Because the SAML assertions contains multiple AWS roles that the user can assume, you can use this app as a middleware for role selection using a browser-based workflow.

## Okta Setup
For setting up an app to provide the SAML assertion, see [Okta AWS authentication](https://saml-doc.okta.com/SAML_Docs/How-to-Configure-SAML-2.0-for-Amazon-Web-Service#scenarioB).

## Quick Start

* Edit `config.yaml` and provide an Okta app's metadata. Without this the app won't start
* Update `config.yaml` and enter rest of AWS accounts info and register apps per account
* Use `make run` to run example server. 

### Docker

Latest Docker image is pushed to `nayyarasamuel7/okta-aws-role-selector`. To run a container for your config:
* Mount directory with your config onto `/root/config`
* When running the container provide these parameters to the docker run: `-c config/<NAME_OF_YOUR_CONFIG_FILE>`. Skip this step if file is named `config.yaml`

##### Example run:

```bash 
docker run -p 80:80 -v $(HOME):/root/config nayyarasamuel7/okta-aws-role-selector:latest  -c config/my_config.yaml
```

## Screenshots

<img height="400px" src="https://raw.githubusercontent.com/nayyara-samuel/okta-aws-role-selector/master/images/role-selector.png">
<img height="300px" src="https://raw.githubusercontent.com/nayyara-samuel/okta-aws-role-selector/master/images/example.png">
