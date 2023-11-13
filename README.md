# findingway

Inspired and indebted to the [Rust project of the same name](https://github.com/epitaque/findingway/) by [epitaque](https://github.com/epitaque). You can see it in action in [Aether PUG DSR](https://discord.gg/aetherpugdsr) in the #pf-checks channel.

findingway scrapes https://xivpf.com/listings every 3 minutes, collects the resulting listings, and posts them onto a Discord channel of your choice. Note that xivpf.com is not particularly accurate and includes private listings; there does not seem to be a way to segment them out at the current time.

## Running

findingway ingests its configuration file at `./config.yaml` to determine what to parse.

findingway requires one environment variable to start:

* **DISCORD_TOKEN**: You have to create a [Discord bot for findingway](https://discord.com/developers/applications). Once you've done so, you can add the bot token here.

findingway also accepts one optional environment variable:

* **ONCE**: If present, findingway will run only once and then exit successfully. Otherwise it will run perpetually and update the target channel every three minutes.

I'm not totally sure if findingway can "just run" in other Discords, even if added. The emojis it uses are present only in APD, and bots can't always use emojis across Discords. If it can't be run in other Discords, I can create a configuration file for mapping roles and jobs to emojis -- someone just open an issue and let me know.

## Deployment

The repository automatically builds Docker images; you can access them if you want to run findingway in your own Discord.

I run this in Fargate for Aether PUG DSR. Here's a task definition you might find useful:

```
{
  "ipcMode": null,
  "executionRoleArn": "arn:aws:iam::AWS_ACCOUNT_ID:role/ecsTaskExecutionRole",
  "containerDefinitions": [
    {
      "dnsSearchDomains": null,
      "environmentFiles": null,
      "logConfiguration": {
        "logDriver": "awslogs",
        "secretOptions": null,
        "options": {
          "awslogs-group": "/ecs/findingway",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      },
      "entryPoint": null,
      "portMappings": [],
      "command": null,
      "linuxParameters": null,
      "cpu": 0,
      "environment": [
        {
          "name": "DATA_CENTRE",
          "value": "Aether"
        },
        {
          "name": "DISCORD_CHANNEL_ID",
          "value": "your discord channel ID"
        },
        {
          "name": "DISCORD_TOKEN",
          "value": "your discord token"
        },
        {
          "name": "DUTY",
          "value": "Dragonsong's Reprise (Ultimate)"
        }
      ],
      "resourceRequirements": null,
      "ulimits": null,
      "dnsServers": null,
      "mountPoints": [],
      "workingDirectory": null,
      "secrets": null,
      "dockerSecurityOptions": null,
      "memory": null,
      "memoryReservation": null,
      "volumesFrom": [],
      "stopTimeout": null,
      "image": "ghcr.io/veraticus/findingway:main",
      "startTimeout": null,
      "firelensConfiguration": null,
      "dependsOn": null,
      "disableNetworking": null,
      "interactive": null,
      "healthCheck": null,
      "essential": true,
      "links": null,
      "hostname": null,
      "extraHosts": null,
      "pseudoTerminal": null,
      "user": null,
      "readonlyRootFilesystem": null,
      "dockerLabels": null,
      "systemControls": null,
      "privileged": null,
      "name": "findingway"
    }
  ],
  "placementConstraints": [],
  "memory": "512",
  "taskRoleArn": "arn:aws:iam::AWS_ACCOUNT_ID:role/ecsTaskExecutionRole",
  "compatibilities": [
    "EC2",
    "FARGATE"
  ],
  "taskDefinitionArn": "arn:aws:ecs:us-east-1:AWS_ACCOUNT_ID:task-definition/findingway:1",
  "family": "findingway",
  "requiresAttributes": [
    {
      "targetId": null,
      "targetType": null,
      "value": null,
      "name": "com.amazonaws.ecs.capability.logging-driver.awslogs"
    },
    {
      "targetId": null,
      "targetType": null,
      "value": null,
      "name": "ecs.capability.execution-role-awslogs"
    },
    {
      "targetId": null,
      "targetType": null,
      "value": null,
      "name": "com.amazonaws.ecs.capability.docker-remote-api.1.19"
    },
    {
      "targetId": null,
      "targetType": null,
      "value": null,
      "name": "com.amazonaws.ecs.capability.task-iam-role"
    },
    {
      "targetId": null,
      "targetType": null,
      "value": null,
      "name": "com.amazonaws.ecs.capability.docker-remote-api.1.18"
    },
    {
      "targetId": null,
      "targetType": null,
      "value": null,
      "name": "ecs.capability.task-eni"
    }
  ],
  "pidMode": null,
  "requiresCompatibilities": [
    "FARGATE"
  ],
  "networkMode": "awsvpc",
  "runtimePlatform": {
    "operatingSystemFamily": "LINUX",
    "cpuArchitecture": null
  },
  "cpu": "256",
  "revision": 1,
  "status": "ACTIVE",
  "inferenceAccelerators": null,
  "proxyConfiguration": null,
  "volumes": []
}
```
