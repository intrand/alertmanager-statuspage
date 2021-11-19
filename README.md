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
    | **statuspageio_component** | the ID of the corresponding statuspage.io component |
    | statuspageio_severity | one of the strings listed in the [documentation](https://developer.statuspage.io/#operation/patchPagesPageIdComponentsComponentId) |

    > Note: **bold** = required

    <br/>

4. Alertmanager sending a webhook to this program.

    Example alertmanager config:

    ```yaml
    # your global section may look different
    global:
      resolve_timeout: 5m

    route:
      # your configuration may vary
      receiver: my_default_receiver
      group_by: ["alertname"]
      group_wait: 30s
      group_interval: 5m
      repeat_interval: 24h

      routes:
        # send watchdog alerts into a blackhole
        # (comes with kube-prometheus)
        - match:
            alertname: Watchdog
          receiver: null

        # send any alert with our labels to the statuspage receiver
        - match_re:
            statuspageio_page: ^.+$
            statuspageio_component: ^.+$
          receiver: statuspageio
          continue: true

    receivers:
      # not strictly necessary to have
      - name: my_default_receiver
        <fill me in>

      # this is required to make alertmanager-statuspage function
      - name: statuspageio
        webhook_configs:
          - url: http://alertmanager-statuspage:8080
    ```

    > In case you don't have or want a default receiver, simply omit it from the list of receivers and change `route.receiver` to `null`.

# Deployment

## Kubernetes / Docker (recommended)

Example manifests are in [/deploy/kubernetes](/deploy/kubernetes) of this repository.

You may find the containerd image at the [docker hub registry](https://hub.docker.com/r/intrand/alertmanager-statuspage).

## SystemD (untested)

<!-- 1. Download the binary using `curl`, move it to `/usr/local/bin` and make it executable:

    ```sh
    curl -sSL -o ~/amsp https://github.com/intrand/alertmanager-statuspage/releases/latest/todo && \
    sudo mv ~/amsp /usr/local/bin/alertmanager-statuspage && \
    sudo chmod 755 /usr/local/bin/alertmanager-statuspage
    ```
-->

1. Getting the binary:

    > I haven't added binaries to github yet. In the meantime, you may build your own. Don't forget to make it executable!

    ```
    sudo chmod 755 /usr/local/bin/alertmanager-statuspage
    ```

3. Configure systemd:


    ```ini
    sudo tee /etc/systemd/system/alertmanager-statuspage 1>/dev/null <<EOF
    [Unit]
    Description=alertmanager receiver for statuspage.io
    Documentation=https://github.com/intrand/alertmanager-statuspage
    After=network.target

    [Service]
    Environment="token=123-abc-456-efg"
    Environment="listen_address=0.0.0.0:8080"
    ExecStart=/usr/local/bin/alertmanager-statuspage
    KillMode=process
    Restart=on-failure
    RestartSec=15s

    [Install]
    WantedBy=multi-user.target
    EOF

    ```

4. Reload systemd

    ```
    sudo systemctl daemon-reload
    ```

5. Start the service

    ```sh
    sudo systemctl start alertmanager-statuspage
    ```

## Manually

```sh
token="123-abc-456-efg" ./alertmanager-statuspage -listen.address 0.0.0.0:8080
```

# Thank you!

[Benjojo](https://github.com/benjojo) helped me a ton by providing the [alertmanager-discord code](https://github.com/benjojo/alertmanager-discord) as a starting point for this project, and during testing. Thank you very much!
