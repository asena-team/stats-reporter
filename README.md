# Stats Reporter [![Made With Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg?color=007EC6)](http://golang.org)

Keep daily progress statistics of your free bot **without** using a **database**.

![image](https://user-images.githubusercontent.com/20264712/114569483-b2d3c480-9c7d-11eb-912e-89d81f699e73.png)

## Features
- Collecting daily server and vote count statistics.
- Sends your day's stats and growth rate live to your Discord Webhook.
- No database requirements.

## Build on Google Infrastructure

> You can probably do these operations more easily with **Terraform** or **Yaml**. But because I don't know, I prefer to follow the path below.

Remember, you need [Google Cloud SDK](https://cloud.google.com/sdk/docs/install) for setup.

**1.** First, install repository
```shell
git clone https://github.com/asena-team/stats-reporter && cd stats-reporter
```

**2.** Set static variables
```shell
# App Name
APP="stats-reporter"

# Get the current project
PROJECT=$(gcloud config get-value core/project 2> /dev/null)
```

**3.** Create service account and download credentials with IAM
```shell
# Create a service account
gcloud iam service-accounts create $APP \
    --description="Optional Stats Reporter Description" \
    --display-name=$APP

# Create and download credentials for the service account
gcloud iam service-accounts keys create credentials.json \
    --iam-account "$APP@$PROJECT.iam.gservice.account.com"
```

**4.** Enable the Google Sheets API
```shell
gcloud services enable sheets.googleapis.com
```

**5.** Create a new spreadsheet and give access to the IAM account
- Create a new Google Sheet if you don’t have one already: [sheets.new](https://sheets.new)
- Share the Google Sheet with the service account email (above) as an editor:

![1_Lo4Mpgxx_tHiIu_rz9Xlqg](https://user-images.githubusercontent.com/20264712/114568399-c16dac00-9c7c-11eb-86c9-74314db96c7f.png)

*(Click that Share button in the top-right of your Sheet.)*

![1_rXUc52V1hxP-16oRrthfbw](https://user-images.githubusercontent.com/20264712/114568590-eb26d300-9c7c-11eb-967a-af0d6a747b02.png)

*(The service account email should autocomplete for a valid service account.)*

**6.** Build container image and upload Cloud Container Registry
```shell
gcloud builds submit --tag "gcr.io/$PROJECT/$APP"
```

**7.** Deploy Cloud Run app with using container image
```shell
gcloud run deploy $APP \
    --image="gcr.io/$PROJECT/$APP" \
    --platform=managed \
    --memory=256Mi \
    --region=us-central1 \
    --args="--bot_id=$BOT_ID,--webhook_id=$WEBHOOK_ID,--webhook_token=$WEBHOOK_TOKEN,--dbl_token=$DBL_TOKEN,--sheet_id=$SHEET_ID"
```

**8.** Bind IAM Policy
```shell
# Authorize your service account with `roles/run.invoker` to access your applation
gcloud run services add-iam-policy-binding stats-reporter \
    --member="$APP@$PROJECT.iam.gserviceaccount.com" \
    --role="roles/run.invoker"
```

**9.** Create a new Cron Scheduler

* It allows to collect statistics by sending requests to our application at the end of each day.

```shell
# --timezone="Your Time Zone"
gcloud scheduler jobs create http asena-stats \ 
    --schedule="0 0 * * *" \
    --uri="https://your.cloud.run.app/run" \
    --http-method=GET \ 
    --timezone="Europe/Istanbul" \
    --oidc-service-account-email="$APP@$PROJECT.iam.gserviceaccount.com"
```

## Local Run & Test

⚠️If you are running locally as a **tool**, you need to set the `'--cli'` parameter

```shell
go run . \
  --bot_id $BOT_ID \ 
  --webhook_id $WEBHOOK_ID \
  --webhook_token $WEBHOOK_TOKEN \
  --dbl_token $DBL_TOKEN
  --sheet_id $SHEET_ID
  --cli
```

## Images

![image](https://user-images.githubusercontent.com/20264712/114570264-5cb35100-9c7e-11eb-9236-d076d3a37916.png)

*(Incoming daily report. It is in Turkish language but you can translate it into English language very easily.)*

---

![image](https://user-images.githubusercontent.com/20264712/114569726-e44c9000-9c7d-11eb-9b56-30968f9fb086.png)

*(A few spreadsheet features.)*
