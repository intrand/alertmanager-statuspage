# About (brief)
alertmanager-statuspage receives a webhook ***from alertmanager***, looks for certain labels and updates a component's status on statuspage.io via the API.

# Requirements

1. An account on statuspage.io (by Atlassian),

    > [Sign up](https://www.atlassian.com/try/cloud/signup?bundle=statuspage) for statuspage.io (by Atlassian). You will need to verify your email address and enroll in the plan of your choosing. A free plan is availble for personal use.

    <br/>

2. An API token from whoever has privileges to update the components/pages,

    > Follow [these instructions](https://developer.statuspage.io/) to get API access set up correctly.

    <br/>

3. Prometheus metrics with certain labels,

    | label | purpose |
    | --- | --- |
    | **statuspageio_page** | the ID of the corresponding statuspage.io page |
    | **statuspageio_component** | the ID of the corresponding statuspage.io page component |
    | statuspageio_severity | one of the strings listed in the [documentation](https://developer.statuspage.io/#operation/patchPagesPageIdComponentsComponentId) |

    Note: **bold** = required

    <br/>

4. Alertmanager sending a webhook to this program.

    Example alertmanager config:

    ```
    route:
      repeat_interval: 12h
      receiver: statuspageio

    receivers:
      - name: 'statuspageio'
        webhook_configs:
          - url: 'http://alertmanager-statuspage:8080'
    ```

# Kubernetes

Example manifests are in [/deploy/kubernetes](/deploy/kubernetes) of this repository.

You may find the image registry here: https://hub.docker.com/r/intrand/alertmanager-statuspage
