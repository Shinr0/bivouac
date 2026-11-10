# Data Model: Forge Portal

## CatalogItem

| Field | Type | Notes |
|-------|------|-------|
| id | string | e.g. `virtual-machine` |
| title | string | |
| description | string | |
| available | boolean | false = coming soon |

## VmParameters

| Field | Type | Default |
|-------|------|---------|
| vmName | string | required |
| machineType | string | `e2-micro` |
| zone | string | `europe-west1-b` |
| diskSizeGb | number | 20 |
| imageFamily | string | `debian-12` |
| spot | boolean | false |
| staticIp | boolean | false |
| dataDiskSizeGb | number | 0 |

## ProvisioningRequest

| Field | Type | Notes |
|-------|------|-------|
| id | string | uuid |
| serviceId | string | `virtual-machine` |
| parameters | VmParameters | |
| status | enum | submitted \| provisioning \| ready \| failed |
| createdAt | string | ISO |

## PortalSettings

| Field | Type | Notes |
|-------|------|-------|
| token | string | GitHub PAT, localStorage |
| repo | string | `owner/name` |
| mode | `mock` \| `github` | default mock |
