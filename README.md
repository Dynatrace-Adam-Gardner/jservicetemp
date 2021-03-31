# jservicetemp

![image](https://user-images.githubusercontent.com/13639658/113127685-5ffd1480-925c-11eb-8e25-1a2032c3d715.png)

## Create Secret

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
