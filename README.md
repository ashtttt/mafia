Once the MFA is enabled, if you are a heavy user of AWS CLI or Terraform or a developer using any of AWS SDKs, you need to generate temporary AWS session credentials using AWS STS api by passing MFA code.

This litle command line tool automatically generates the MFA token, calls the AWS STS api for temporaty kesy and updates shared credential (.aws/credentials) file with retrieved temporary session keys before they expire. It's called "MaFiA" :)

**Usage**
1. Download the executable from, and place it in a location where PATH is set. 
2. When you are enabling MFA on AWS console, note down the secret key beneath the QR code image. You need to click on "Show secret key for manual configuration" link to open it. We just need this key, rest all you can continue configuring Authy as usual. If you have already configured MFA, you may need to deactivate and re-activate again to get to the key.
3. Now set or export the above as "TOPT_SECRET". "mafia" tool uses this environment variable to find the key and generate MFA token. Please make sure to set this ENV variable as permanent, if you lose this environment variable "mafia" will not be able to generate MFA key and you need to redo the setup all over again. (or at least save the key in a safe place, so can set it again if lose the environment variable).
4. Hope you already have your static keys in .aws/credential file. If they are under "[default]" profile, move them under a different profile name.
5. Ex:
    ```
	[nonprod]
	aws_access_key_id     = XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
	aws_secret_access_key = XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
	```
6. Now you can run "mafia keys --profile=<nonprod>", which will create a "[default]" profile with temporary session credentials. If you leave the command window open, it will automatically rotate the keys for every 12 hours.
7. Run "mafia -help" for other options.
