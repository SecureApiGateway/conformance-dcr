# Config Templating
The dcr binary requires config supplied as json.

To create the json the jq tool is used. Template files in this directory are passed to jq along with a set of variables 
to produce a json document.

## sbat-create-dcr-config.jq

This template is used by the pipeline which creates DCR registrations for use in testing. 
See [codefresh-create-fapi-dcr.yml](../codefresh-create-fapi-dcr.yml) to see how the template is used.

### Notes on template variables

`$transport_cert` `$signing_private_key` `$transport_private_key`  variables are expected to be supplied as quoted json strings.
The reason for this is that the certs and private keys are encoded on a single line with line breaks replaced by \n strings.
This needs to be supplied as `--argjson` to jq to prevent jq escaping the \n chars as \\n.

`$redirect_uris_str` variable is expected to be supplied as a comma separated list of uri strings

Example: `"https://www.google.co.uk,https://www.certification.openid.net/test/a/forgerock-sbat-dev/callback&scope=openid%20accounts&response_type=code%20id_token"`

## sbat-conformance-config.jq

This template is used by the pipeline that runs the full DCR conformance test suite.
