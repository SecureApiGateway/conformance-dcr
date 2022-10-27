# Fork of OB Conformance DCR tool
Upstream repo: https://github.com/OpenBankingUK/conformance-dcr

## Repo structure

The existing master branch will replicate the master in the upstream.

A new branch: `sbat-master` has been created and made the default branch of the repo.
All SBAT related changes must be made to this branch, the branch is protected so feature branch and PR flow must be followed.

The upstream master will periodically be synced and can then be merged into sbat-master if we want to make use of the new commits.

# SBAT Customisations
This fork has been created to enable SBAT customisations to be applied to the conformance tool.

In addition to running the conformance suite, the tool has been customised to enable the registration of test TPP applications.
This customisation reuses the test code but only runs the DCR32CreateSoftwareClient scenario. This is used to create test applications in development environments.

## New configuration options
To customise behaviour, it has been necessary to create additional configuration options which can be used in the config.json

| Key                                  | Type    | Usage                                                                                                                                                                                   |
|--------------------------------------|---------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| create_software_client_only          | boolean | Used to control whether to run the whole suite or to only create a software client <br />Disabled by default                                                                            |
| preferred_token_endpoint_auth_method | string  | (Optional) If specified then the DCR request will try to use this value as the token_endpoint_auth_method in the registration<br /><br />For SBAT we currently prefer `private_key_jwt` |


## SBAT Custom Docker Image
[Dockerfile-sbat](Dockerfile-sbat) creates an image based on alpine, whereas the original Dockerfile produces a distroless image

Using alpine as the base allow us to run bash and do post-processing of the dcr output.
In the pipeline we use this to capture the client_id that was created.

Currently, the image has to be built manually by doing: `make build_sbat_image`
The image produced then needs to be pushed to our container registry
