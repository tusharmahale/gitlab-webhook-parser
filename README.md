# Gitlab Webhook Parser  
This parser is designed to automatically update access permissions of Gitlab protected branches when merge events occur.

## Requirements
* Gitlab instance version 11.7 or higher
* Access to the administration settings

## Installation
1. Log into the Gitlab instance and navigate to the administration settings page. 
2. Select the "Webhooks" setting.
3. Create a webhook, setting the URL to the endpoint for the parser server. 
4. Enter the secret token key, and set the event type to "Merge Request Events". 
5. Save the webhook and verify that it is working correctly.
6. Deploy the parser server and configure it to run on startup. 
7. Configure the parser server with the secret token key and the appropriate access settings for protected branches.  
8. Test the parser server to ensure that merge request events are being correctly detected and access permissions are being updated.

## Usage
Once the parser is installed and configured, it will detect when merge request events occur and will update the protected branches access settings. No user intervention is necessary. 

If any changes are needed, they can be made directly in the parser server configuration.

## Contributing 
Contributions are welcome! Please make sure that any code submitted is well tested and documented. 