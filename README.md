This service reacts to `evaluation.finished` and `remediation.finished` events.

It will:
- Create a JIRA ticket with the pertinent details
- Any labels in the keptn YAML files will be created as labels in JIRA
- Labels for `keptn_project`, `keptn_service` and `keptn_stage` will be created
- For quality evaluations, an additional `keptn_result` JIRA label will be created
- Send an event notification into your tool (currently supports Dynatrace)
- The JIRA ticket will link directly to the sequence in Keptn's bridge
- The event will link directly to the sequence in Keptn's bridge

## Quality Evaluations
![image](https://user-images.githubusercontent.com/13639658/113381981-acaf2f80-93c3-11eb-9ba6-34017e88f2ac.png)

## Remediation Output
![image](https://user-images.githubusercontent.com/13639658/113382057-e8e29000-93c3-11eb-92c2-7ec33d76ad9c.png)

## How To Deploy
1. Gather the required information below
2. Create your `kubectl` secret
3. Customise the Keptn URL and Keptn Bridge URL in the `deploy/service.yaml` file
4. If required, disable event sending by setting `SEND_EVENT` to `"false"` in `deploy/service.yaml`
5. If required, enable debug mode by setting `DEBUG` to `"true"` in `deploy/service.yaml`
6. Deploy with `kubectl apply -f deploy/service.yaml`

## Debugging
This now uses default kubernetes logging

1. Note the pod name of the `jira-service` using `kubectl get pods -n keptn`
2. `kubectl logs -n keptn jira-service-*-* jira-service`

## Required Information

- JIRA Base URL (without trailing slash) eg. `https://abc123.atlassian.net`
- JIRA Username eg. `joe.smith@example.com`
- JIRA ID for Ticket Reporter (see below for how to retrieve)
- JIRA ID for Ticket Assignee (if different from reporter ID)
- JIRA API Token ([generate one here](https://id.atlassian.com/manage/api-tokens))
- JIRA Project Key. Take this from the URL. Eg. PROJ is the project code for `https://abc123.atlassian.net/projects/PROJ/issues`
- JIRA Issue Type eg. Task, Bug, Epic etc. Defaults to `Task`.
- Keptn base URL (eg. http://localhost:8080 or however you've exposed Keptn)
- Keptn bridge URL (eg. http://localhost:8080/bridge)

## Retrieve User IDs (IMPORTANT)
JIRA now required the User ID for both the ticket reporter and the assignee.

Retrieve these by clicking your profile icon (top right) then go to profile and grab your ID from the end of the URL:

![image](https://user-images.githubusercontent.com/13639658/113224119-0a615000-92ce-11eb-9abd-693efa2ac612.png)

## Create Secret

Requires a secret called `jira-details` in `keptn` namespace:

```
kubectl -n keptn create secret generic jira-details \
--from-literal="jira-base-url=***" \
--from-literal="jira-username=***" \
--from-literal="jira-reporter-user-id=***" \
--from-literal="jira-assignee-user-id=***" \
--from-literal="jira-api-token=***" \
--from-literal="jira-project-key=***" \
--from-literal="jira-issue-type=Task" \
--from-literal="jira-create-ticket-for-problems=true" \
--from-literal="jira-create-ticket-for-evaluations=true"
```
## Sending Events
Set `SEND_EVENT` to `true` or `false` in `deploy/service.yaml` to get an alert into your tool (or not).
Currently supports Dynatrace and the service must have `keptn_project`, `keptn_service` and `keptn_stage` variables set.
