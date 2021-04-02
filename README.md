# jservicetemp

![image](https://user-images.githubusercontent.com/13639658/113127685-5ffd1480-925c-11eb-8e25-1a2032c3d715.png)

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

### Retrieve User IDs
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
## Set Send Event Flag
Set `SEND_EVENT` to `true` or `false` in `deploy/service.yaml` to get an alert into your tool (or not).
Currently supports Dynatrace and the service must have `keptn_project`, `keptn_service` and `keptn_stage` variables set.

## Set Debug Flag
Set `DEBUG` to `true` or `false` in `deploy/service.yaml` to get extra log lines.

## Deploy
```
kubectl apply -f deploy/service.yaml
```

## View Logs
```
kubectl get pods -n keptn
kubectl logs -n keptn jira-service-*-* jira-service
```
