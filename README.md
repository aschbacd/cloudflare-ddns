# Cloudflare DDNS

A simple tool that allows you to use Cloudflare as a ddns provider so that you can use custom domains for free without having a static public ip address. Unfortunately .cf, .ga, .gq, .ml, and .tk TLDs are not supported anymore by Cloudflare.

Before using the tool a configuration file must be created. To do this simply run the following command:

```bash
cloudflare-ddns configure
```

After starting the executable with the `configure` command the configurator will ask you for your authentication email address and key. For authentication email use your Cloudflare account email. For authentication key use your global Cloudflare api key that can be viewed under My Profile / API Tokens.

![Screenshot API Key](images/screenshot_api_key.png)

After finishing the configuration you can add the binary to Crontab on Linux or to Task Scheduler in Windows to automate your domain updates. Therefore you can use the `update` command.

Thanks to [@janunterkofler](https://github.com/janunterkofler) for the initial bash script!
