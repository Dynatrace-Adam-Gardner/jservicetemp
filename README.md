# jservicetemp

Requires a secret called `jira-details` in `keptn` namespace:

```
kubectl -n keptn create secret generic jira-details \
--from-literal="jira-base-url=***" \
--from-literal="jira-username=***" \
--from-literal="jira-api-token=***" \
--from-literal="jira-project-key=***" \
--from-literal="jira-issue-type=Task" \
--from-literal="jira-create-ticket-for-problems=true" \
--from-literal="jira-create-ticket-for-evaluations=true"
```

## Set Debug Flag
Set `DEBUG` to `true` or `false` in `deploy/service.yaml` to get extra log lines:

## Deploy
```
kubectl apply -f deploy/service.yaml
```

## View Logs
```
kubectl get pods -n keptn
kubectl logs -n keptn jira-service-*-* jira-service
```
